package data

import (
	"gearboxd/internal/assert"
	"gearboxd/internal/validator"
	"testing"
)

func TestValidateUsername(t *testing.T) {
	tests := []struct {
		title          string
		username       string
		valid          bool
		expectedErrors map[string]string
	}{
		{
			"valid username",
			"david",
			true,
			nil,
		},
		{
			"username is empty",
			"",
			false,
			map[string]string{"username": "cannot be empty"},
		},
		{
			"username is shorter than 3 characters",
			"da",
			false,
			map[string]string{"username": "cannot be less than 3 characters"},
		},
		{
			"username is longer than 64 characters",
			"daviddaviddaviddaviddaviddaviddaviddaviddaviddaviddaviddaviddavid", //65 chars
			false,
			map[string]string{"username": "cannot be more than 64 characters"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			v := validator.New()

			ValidateUsername(v, tt.username)

			assert.Equal(t, v.Valid(), tt.valid)
			if tt.valid == false && v.Valid() == false {
				assert.Equal(t, v.Errors, tt.expectedErrors)
			}
		})
	}
}

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		title          string
		email          string
		valid          bool
		expectedErrors map[string]string
	}{
		{
			"valid email",
			"david@gmail.com",
			true,
			nil,
		},
		{
			"email is empty",
			"",
			false,
			map[string]string{"email": "cannot be empty"},
		},
		{
			"email is not formatted correctly",
			"david.com",
			false,
			map[string]string{"email": "format is not valid"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			v := validator.New()

			ValidateEmail(v, tt.email)

			assert.Equal(t, v.Valid(), tt.valid)
			if tt.valid == false && v.Valid() == false {
				assert.Equal(t, v.Errors, tt.expectedErrors)
			}
		})
	}
}

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		title          string
		password       string
		valid          bool
		expectedErrors map[string]string
	}{
		{
			"valid password",
			"safepassword",
			true,
			nil,
		},
		{
			"password is empty",
			"",
			false,
			map[string]string{"password": "cannot be empty"},
		},
		{
			"password is shorter than 8 characters",
			"safepas",
			false,
			map[string]string{"password": "cannot be less 8 characters"},
		},
		{
			"password is longer than 64 characters",
			"daviddaviddaviddaviddaviddaviddaviddaviddaviddaviddaviddaviddavid", //65 chars
			false,
			map[string]string{"password": "cannot be more than 64 characters"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			v := validator.New()

			ValidatePassword(v, tt.password)

			assert.Equal(t, v.Valid(), tt.valid)
			if tt.valid == false && v.Valid() == false {
				assert.Equal(t, v.Errors, tt.expectedErrors)
			}
		})
	}
}

func TestValidateUser(t *testing.T) {
	tests := []struct {
		title          string
		username       string
		email          string
		password       string
		valid          bool
		expectedErrors map[string]string
	}{
		{
			title:    "valid user",
			username: "david",
			email:    "david@gmail.com",
			password: "safepassword",
			valid:    true,
		},
		{
			title:    "invalid username",
			username: "da",
			email:    "david@gmail.com",
			password: "safepassword",
			valid:    false,
			expectedErrors: map[string]string{
				"username": "cannot be less than 3 characters",
			},
		},
		{
			title:    "invalid email",
			username: "david",
			email:    "david.com",
			password: "safepassword",
			valid:    false,
			expectedErrors: map[string]string{
				"email": "format is not valid",
			},
		},
		{
			title:    "invalid password",
			username: "david",
			email:    "david@gmail.com",
			password: "short",
			valid:    false,
			expectedErrors: map[string]string{
				"password": "cannot be less 8 characters",
			},
		},
		{
			title:    "multiple invalid fields",
			username: "",
			email:    "david.com",
			password: "",
			valid:    false,
			expectedErrors: map[string]string{
				"username": "cannot be empty",
				"email":    "format is not valid",
				"password": "cannot be empty",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			user := &User{
				Username: tt.username,
				Email:    tt.email,
			}
			user.Password.plaintext = &tt.password

			v := validator.New()

			ValidateUser(user, v)

			assert.Equal(t, v.Valid(), tt.valid)
			if !tt.valid {
				assert.Equal(t, v.Errors, tt.expectedErrors)
			}
		})
	}
}
