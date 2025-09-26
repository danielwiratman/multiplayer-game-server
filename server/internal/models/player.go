package models

type UpdatePlayerRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}
