package service

import "strings"

const (
	DemoUsername = "student"
	DemoPassword = "student"
	DemoToken    = "demo-token"
	DemoSubject  = "student"
)

type Service struct{}

func New() *Service {
	return &Service{}
}

func (s *Service) Login(username, password string) (string, bool) {
	if username == DemoUsername && password == DemoPassword {
		return DemoToken, true
	}

	return "", false
}

func (s *Service) VerifyAuthorizationHeader(header string) (string, bool) {
	token, ok := parseBearerToken(header)
	if !ok {
		return "", false
	}

	if token != DemoToken {
		return "", false
	}

	return DemoSubject, true
}

func parseBearerToken(header string) (string, bool) {
	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 {
		return "", false
	}

	if parts[0] != "Bearer" || strings.TrimSpace(parts[1]) == "" {
		return "", false
	}

	return strings.TrimSpace(parts[1]), true
}
