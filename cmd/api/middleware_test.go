package main

import (
	"errors"
	"gearboxd/internal/assert"
	"gearboxd/internal/data"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthenticate(t *testing.T) {
	tests := []struct {
		name            string
		authHeader      string
		getForTokenFunc func(scope, token string) (*data.User, error)
		wantCode        int
		wantUser        *data.User
	}{
		{
			name:       "no auth header sets anonymous user",
			authHeader: "",
			wantCode:   http.StatusOK,
			wantUser:   data.AnonymousUser,
		},
		{
			name:       "valid token sets user in context",
			authHeader: "Bearer ABCDEFGHIJKLMNOPQRSTUVWXYZ",
			getForTokenFunc: func(scope, token string) (*data.User, error) {
				return &data.User{ID: 1, Activated: true}, nil
			},
			wantCode: http.StatusOK,
			wantUser: &data.User{ID: 1, Activated: true},
		},
		{
			name:       "malformed auth header",
			authHeader: "NotBearer something",
			wantCode:   http.StatusUnauthorized,
		},
		{
			name:       "auth header with no space",
			authHeader: "Bearertoken",
			wantCode:   http.StatusUnauthorized,
		},
		{
			name:       "invalid token format",
			authHeader: "Bearer short",
			wantCode:   http.StatusUnauthorized,
		},
		{
			name:       "token not found in DB",
			authHeader: "Bearer ABCDEFGHIJKLMNOPQRSTUVWXYZ",
			getForTokenFunc: func(scope, token string) (*data.User, error) {
				return nil, data.ErrRecordNotFound
			},
			wantCode: http.StatusUnauthorized,
		},
		{
			name:       "database error",
			authHeader: "Bearer ABCDEFGHIJKLMNOPQRSTUVWXYZ",
			getForTokenFunc: func(scope, token string) (*data.User, error) {
				return nil, errors.New("connection refused")
			},
			wantCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := newTestApplication(t, nil, &data.Models{
				Users: &MockUserStore{GetForTokenFunc: tt.getForTokenFunc},
			})

			var gotUser *data.User
			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				gotUser = app.contextGetUser(r)
				w.WriteHeader(http.StatusOK)
			})

			rr := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			app.authenticate(next).ServeHTTP(rr, req)

			assert.Equal(t, rr.Code, tt.wantCode)

			if tt.wantUser != nil {
				assert.Equal(t, gotUser.ID, tt.wantUser.ID)
			}
		})
	}
}

func TestRequireAuthenticatedUser(t *testing.T) {
	tests := []struct {
		name        string
		contextUser *data.User
		code        int
	}{
		{
			name:        "Anonymous user",
			contextUser: data.AnonymousUser,
			code:        http.StatusUnauthorized,
		},
		{
			name:        "Authenticated user",
			contextUser: &data.User{ID: 1},
			code:        http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := newTestApplication(t, nil, nil)
			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			rr := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/", nil)
			req = app.contextSetUser(req, tt.contextUser)
			app.requireAuthenticatedUser(next).ServeHTTP(rr, req)

			assert.Equal(t, tt.code, rr.Code)
		})
	}
}

func TestRequireActivatedUser(t *testing.T) {
	tests := []struct {
		name        string
		contextUser *data.User
		code        int
	}{
		{
			name:        "User activated",
			contextUser: &data.User{Activated: true},
			code:        http.StatusOK,
		},
		{
			name:        "User not activated",
			contextUser: &data.User{Activated: false},
			code:        http.StatusForbidden,
		},
		{
			name:        "Anonymous user not authenticated",
			contextUser: data.AnonymousUser,
			code:        http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := newTestApplication(t, nil, nil)
			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			rr := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/", nil)
			req = app.contextSetUser(req, tt.contextUser)

			app.requireActivatedUser(next).ServeHTTP(rr, req)
			assert.Equal(t, rr.Code, tt.code)
		})
	}
}

func TestRequirePermission(t *testing.T) {
	tests := []struct {
		name               string
		contextUser        *data.User
		getAllForUserFunc  func(userID int64) (data.Permissions, error)
		requiredPermission string

		code int
	}{
		{
			name:        "User has permission",
			contextUser: &data.User{ID: 1, Activated: true},
			getAllForUserFunc: func(userID int64) (data.Permissions, error) {
				return data.Permissions{"cars:write"}, nil
			},
			requiredPermission: "cars:write",
			code:               http.StatusOK,
		},
		{
			name:        "User has no permission",
			contextUser: &data.User{ID: 1, Activated: true},
			getAllForUserFunc: func(userID int64) (data.Permissions, error) {
				return data.Permissions{"cars:read"}, nil
			},
			requiredPermission: "cars:write",
			code:               http.StatusForbidden,
		},
		{
			name:        "User not activated",
			contextUser: &data.User{ID: 1, Activated: false},
			getAllForUserFunc: func(userID int64) (data.Permissions, error) {
				return nil, nil
			},
			requiredPermission: "cars:write",
			code:               http.StatusForbidden,
		},
		{
			name:        "Anonymous user",
			contextUser: data.AnonymousUser,
			getAllForUserFunc: func(userID int64) (data.Permissions, error) {
				return nil, nil
			},
			requiredPermission: "cars:write",
			code:               http.StatusUnauthorized,
		},
		{
			name:        "Database error",
			contextUser: &data.User{ID: 1, Activated: true},
			getAllForUserFunc: func(userID int64) (data.Permissions, error) {
				return nil, errors.New("connection refused")
			},
			requiredPermission: "cars:write",
			code:               http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := newTestApplication(t,
				nil,
				&data.Models{
					Permissions: MockPermissionStore{
						GetAllForUserFunc: tt.getAllForUserFunc},
				},
			)

			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			rr := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/", nil)
			req = app.contextSetUser(req, tt.contextUser)

			app.requirePermission(tt.requiredPermission, next).ServeHTTP(rr, req)
			assert.Equal(t, rr.Code, tt.code)
		})
	}
}
