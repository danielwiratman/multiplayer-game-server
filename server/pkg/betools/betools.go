// Package betools is a collection of tools for http server made by Daniel W.
package betools

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"server/internal/models"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type Middleware func(http.Handler) http.Handler

type Route struct {
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
	Middlewares []Middleware
}

func WithMiddlewares(mws []Middleware, rs []Route) []Route {
	var routes []Route
	for _, r := range rs {
		routes = append(routes, Route{
			Method:      r.Method,
			Pattern:     r.Pattern,
			HandlerFunc: r.HandlerFunc,
			Middlewares: append(mws, r.Middlewares...),
		})
	}
	return routes
}

type Controller interface {
	GetRoutes() []Route
}

type Router struct {
	controllers []Controller
}

func NewRouter(controllers ...Controller) *Router {
	return &Router{
		controllers: controllers,
	}
}

func (me *Router) Route(r chi.Router) {
	for _, controller := range me.controllers {
		for _, route := range controller.GetRoutes() {
			r.Group(func(r chi.Router) {
				r.Use(
					SliceMap(route.Middlewares, func(m Middleware) func(http.Handler) http.Handler {
						return (func(next http.Handler) http.Handler)(m)
					})...,
				)
				r.MethodFunc(
					route.Method,
					route.Pattern,
					route.HandlerFunc,
				)
			})
		}
	}
}

func SliceMap[T any, U any](s []T, transform func(el T) U) []U {
	var res []U
	for _, el := range s {
		res = append(res, transform(el))
	}
	return res
}

type FieldError struct {
	Key     string `json:"key"`
	Message string `json:"message"`
	Param   string `json:"param"`
}

type MessageData struct {
	Message string `json:"message"`
}

type Pagination struct {
	Total int `json:"total"`
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

type Response struct {
	Success    bool         `json:"success"`
	Data       any          `json:"data"`
	Error      string       `json:"error,omitempty"`
	Errors     []FieldError `json:"errors,omitempty"`
	Pagination *Pagination  `json:"pagination,omitempty"`
}

type SSEEvents struct {
	Event string `json:"event"`
	Data  any    `json:"data"`
	Error error  `json:"error"`
}

func SendEventsResponse(w http.ResponseWriter, flusher http.Flusher, event SSEEvents) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	fmt.Fprintf(w, "event: %s\n", event.Event)
	if event.Error != nil {
		payloadBytes, err := json.Marshal(map[string]string{
			"error": event.Error.Error(),
		})
		if err != nil {
			return
		}
		fmt.Fprintf(w, "data: %s\n\n", payloadBytes)

	} else {
		payloadBytes, err := json.Marshal(event.Data)
		if err != nil {
			return
		}
		fmt.Fprintf(w, "data: %s\n\n", payloadBytes)

	}
	flusher.Flush()
}

func SendErrorResponse(w http.ResponseWriter, code int, data any, errors ...FieldError) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)

	var message string

	switch t := data.(type) {
	case string:
		message = t
	case error:
		message = t.Error()
	default:
		message = http.StatusText(code)
	}

	res := &Response{
		Success: false,
		Error:   message,
		Errors:  errors,
	}
	json.NewEncoder(w).Encode(res)
}

func SendSuccessResponse(w http.ResponseWriter, code int, args ...any) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)

	var responseData any
	var pagination *Pagination

	switch l := len(args); {
	case l > 1:
		// parse pagination if exists
		paginationObj, ok := args[1].(Pagination)
		if ok {
			pagination = &paginationObj
		}
		fallthrough
	case l > 0:
		// parse response data if exists
		switch t := args[0].(type) {
		case string:
			responseData = &MessageData{Message: t}
		default:
			responseData = &t
		}
		fallthrough
	default:
		if responseData == nil {
			// response data not set, use status text as message
			responseData = MessageData{Message: http.StatusText(code)}
		}
		res := &Response{
			Success:    true,
			Data:       responseData,
			Pagination: pagination,
		}

		json.NewEncoder(w).Encode(res)
	}
}

func SendOKResponse(w http.ResponseWriter, args ...any) {
	SendSuccessResponse(w, http.StatusOK, args...)
}

func SendCreatedResponse(w http.ResponseWriter, args ...any) {
	SendSuccessResponse(w, http.StatusCreated, args...)
}

type BodyParserOptions struct {
	Validate bool
	Field    string
}

func BodyParser[T any](opts ...BodyParserOptions) Middleware {
	opt := BodyParserOptions{
		Validate: true,
		Field:    "",
	}
	if len(opts) > 0 {
		opt = opts[0]
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var body []byte
			if opt.Field != "" {
				body = []byte(r.FormValue(opt.Field))
			} else {
				var err error
				body, err = io.ReadAll(io.LimitReader(r.Body, 10<<14))
				if err != nil {
					SendErrorResponse(w, http.StatusBadRequest, "read body failed")
					return
				}
			}

			var data T
			if err := json.Unmarshal(body, &data); err != nil {
				SendErrorResponse(w, http.StatusBadRequest, "json unmarshal failed")
				return
			}

			if opt.Validate {
				fieldErrors := Validate(&data)
				if len(fieldErrors) > 0 {
					SendErrorResponse(
						w,
						http.StatusBadRequest,
						"validation failed",
						SliceMap(fieldErrors, func(el ValidationError) FieldError {
							return FieldError(el)
						})...,
					)
					return
				}
			}

			next.ServeHTTP(w, SetContext(r, CtxKeyBody, data))
		})
	}
}

var V *validator.Validate = validator.New()

type ValidationError struct {
	Key     string
	Message string
	Param   string
}

func Validate[T any](data *T) []ValidationError {
	err := V.Struct(data)
	if err == nil {
		return nil
	}
	if _, ok := err.(validator.ValidationErrors); !ok {
		return nil
	}

	var errors []ValidationError
	for _, err := range err.(validator.ValidationErrors) {
		fieldName := err.StructField()
		field, _ := reflect.TypeOf(data).Elem().FieldByName(fieldName)
		if jsonTag, ok := field.Tag.Lookup("json"); ok {
			name := strings.SplitN(jsonTag, ",", 2)[0]
			if name != "-" {
				fieldName = name
			}
		}
		errors = append(errors, ValidationError{
			Key:     fieldName,
			Message: err.Tag(),
			Param:   err.Param(),
		})
	}
	return errors
}

type CtxKey int

func getContext[T any](r *http.Request, ctxLabel CtxKey) T {
	var def T

	val := r.Context().Value(ctxLabel)
	if val == nil {
		return def
	}

	if val, ok := val.(T); ok {
		return val
	}

	return def
}

func SetContext(r *http.Request, key CtxKey, val any) *http.Request {
	ctx := context.WithValue(r.Context(), key, val)
	return r.WithContext(ctx)
}

/* BELOW THIS IS EDITABLE */

const (
	CtxKeyBody CtxKey = iota
	CtxKeyAuth
)

func GetBodyCtx[T any](r *http.Request) T {
	return getContext[T](r, CtxKeyBody)
}

func GetAuthCtx(r *http.Request) models.Account {
	return getContext[models.Account](r, CtxKeyAuth)
}
