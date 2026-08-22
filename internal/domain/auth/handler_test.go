package auth

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type stubAuthService struct {
	loginResult    *AuthResult
	loginErr       error
	registerResult *AuthResult
	registerErr    error
}

func (s stubAuthService) Login(_, _ string) (*AuthResult, error) {
	return s.loginResult, s.loginErr
}

func (s stubAuthService) Register(_, _, _ string) (*AuthResult, error) {
	return s.registerResult, s.registerErr
}

func TestHandleLoginErrorMapping(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		wantStatus int
		wantBody   string
		notInBody  string
	}{
		{
			name:       "invalid credentials",
			err:        ErrInvalidCredentials,
			wantStatus: http.StatusUnauthorized,
			wantBody:   "invalid credentials",
		},
		{
			name:       "internal error",
			err:        errors.New("database password leaked"),
			wantStatus: http.StatusInternalServerError,
			wantBody:   "internal server error",
			notInBody:  "database password leaked",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			handler := NewHandler(stubAuthService{loginErr: test.err})
			request := httptest.NewRequest(http.MethodPost, "/api/login", strings.NewReader(`{"email":"user@example.com","password":"password"}`))
			response := httptest.NewRecorder()

			handler.HandleLogin().ServeHTTP(response, request)

			if response.Code != test.wantStatus {
				t.Fatalf("status = %d, want %d", response.Code, test.wantStatus)
			}
			if !strings.Contains(response.Body.String(), test.wantBody) {
				t.Errorf("body = %q, want it to contain %q", response.Body.String(), test.wantBody)
			}
			if test.notInBody != "" && strings.Contains(response.Body.String(), test.notInBody) {
				t.Errorf("body exposes internal error: %q", response.Body.String())
			}
		})
	}
}

func TestHandleRegisterErrorMapping(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		err        error
		wantStatus int
		wantBody   string
	}{
		{
			name:       "malformed body",
			body:       `{`,
			wantStatus: http.StatusBadRequest,
			wantBody:   "invalid request body",
		},
		{
			name:       "invalid fields",
			body:       `{"email":"bad","username":"x","password":"short"}`,
			wantStatus: http.StatusBadRequest,
			wantBody:   "invalid request",
		},
		{
			name:       "email taken",
			body:       `{"email":"user@example.com","username":"user","password":"password"}`,
			err:        ErrEmailTaken,
			wantStatus: http.StatusConflict,
			wantBody:   "email already exists",
		},
		{
			name:       "username taken",
			body:       `{"email":"user@example.com","username":"user","password":"password"}`,
			err:        ErrUsernameTaken,
			wantStatus: http.StatusConflict,
			wantBody:   "username already exists",
		},
		{
			name:       "internal error",
			body:       `{"email":"user@example.com","username":"user","password":"password"}`,
			err:        errors.New("database exploded"),
			wantStatus: http.StatusInternalServerError,
			wantBody:   "internal server error",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			handler := NewHandler(stubAuthService{registerErr: test.err})
			request := httptest.NewRequest(http.MethodPost, "/api/register", strings.NewReader(test.body))
			response := httptest.NewRecorder()

			handler.HandleRegister().ServeHTTP(response, request)

			if response.Code != test.wantStatus {
				t.Fatalf("status = %d, want %d", response.Code, test.wantStatus)
			}
			if !strings.Contains(response.Body.String(), test.wantBody) {
				t.Errorf("body = %q, want it to contain %q", response.Body.String(), test.wantBody)
			}
		})
	}
}
