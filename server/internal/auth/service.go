package auth

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Register(username, password string) (int, error) {
	return 1, nil
}
