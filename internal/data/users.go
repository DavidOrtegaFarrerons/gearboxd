package data

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"errors"
	"gearboxd/internal/validator"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrDuplicateUsername = errors.New("duplicate username")
	ErrDuplicateEmail    = errors.New("duplicate email")
)

type password struct {
	plaintext *string
	hash      []byte
}

func (p *password) Set(plaintextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return err
	}

	p.plaintext = &plaintextPassword
	p.hash = hash

	return nil
}

func (p *password) Matches(plaintextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plaintextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}

var AnonymousUser = &User{}

type User struct {
	ID        int64     `json:"id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	Password  password  `json:"-"`
	Activated bool      `json:"activated"`
	Version   int       `json:"-"`
	CreatedAt time.Time `json:"createdAt"`
}

func (u *User) IsAnonymous() bool {
	return u == AnonymousUser
}

func ValidateUsername(v *validator.Validator, username string) {
	v.Check(username != "", "username", "cannot be empty")
	v.Check(len(username) >= 3, "username", "cannot be less than 3 characters")
	v.Check(len(username) <= 64, "username", "cannot be more than 64 characters")
}

func ValidateEmail(v *validator.Validator, email string) {
	v.Check(email != "", "email", "cannot be empty")
	v.Check(validator.Matches(email, validator.EmailRX), "email", "format is not valid")
}

func ValidatePassword(v *validator.Validator, password string) {
	v.Check(password != "", "password", "cannot be empty")
	v.Check(len(password) >= 8, "password", "cannot be less 8 characters")
	v.Check(len(password) <= 64, "password", "cannot be more than 64 characters")
}

func ValidateUser(user *User, v *validator.Validator) {
	ValidateUsername(v, user.Username)
	ValidateEmail(v, user.Email)
	ValidatePassword(v, *user.Password.plaintext)
}

type UserModelInterface interface {
	Insert(user *User) error
	Update(user *User) error
	GetForToken(scope, tokenPlaintext string) (*User, error)
	GetByEmail(email string) (*User, error)
}

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(user *User) error {
	query := `INSERT INTO users (email, username, password_hash) 
	VALUES ($1, $2, $3)
	RETURNING id, version, created_at`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, user.Email, user.Username, user.Password.hash).Scan(&user.ID, &user.Version, &user.CreatedAt)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "users_username_key"):
			return ErrDuplicateUsername
		case strings.Contains(err.Error(), "users_email_key"):
			return ErrDuplicateEmail
		default:
			return err
		}
	}

	return nil
}

func (m *UserModel) Update(user *User) error {
	query := `
		UPDATE users
		SET username = $1, email = $2, password_hash = $3, activated = $4, version = version + 1
		WHERE id = $5 AND version = $6
		RETURNING version`

	args := []any{
		user.Username,
		user.Email,
		user.Password.hash,
		user.Activated,
		user.ID,
		user.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&user.Version)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "users_username_key"):
			return ErrDuplicateUsername
		case strings.Contains(err.Error(), "users_email_key"):
			return ErrDuplicateEmail
		default:
			return err
		}
	}

	return nil
}

func (m *UserModel) GetForToken(scope, tokenPlaintext string) (*User, error) {
	tokenHash := sha256.Sum256([]byte(tokenPlaintext))

	query := `SELECT users.id, users.created_at, users.username, users.email, users.password_hash, users.activated, users.version
	FROM users
	INNER JOIN tokens
	ON users.id = tokens.user_id
	WHERE tokens.hash = $1
	AND tokens.scope = $2
	AND tokens.expiry > $3`

	args := []any{tokenHash[:], scope, time.Now()}

	var user User

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.Username,
		&user.Email,
		&user.Password.hash,
		&user.Activated,
		&user.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}

func (m UserModel) GetByEmail(email string) (*User, error) {
	query := `
		SELECT id, created_at, username, email, password_hash, activated, version
		FROM users
		WHERE email = $1`

	var user User

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.Username,
		&user.Email,
		&user.Password.hash,
		&user.Activated,
		&user.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err

		}
	}

	return &user, nil
}
