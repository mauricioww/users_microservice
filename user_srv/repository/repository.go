package repository

import (
	"context"
	"database/sql"

	"github.com/go-kit/log"
	"github.com/mauricioww/user_microsrv/errors"
	"github.com/mauricioww/user_microsrv/user_srv/entities"
)

const (
	createUserSQL = `
		INSERT INTO USERS(email, pwd_hash, age)
			VALUES (?, ?, ?)
	`

	authenticateSQL = `
		SELECT u.pwd_hash FROM USERS u WHERE u.email = ?
	`

	updateUserSQL = `
		UPDATE USERS SET email = ?, pwd_hash = ?, age = ?
			WHERE id = ? 
	`

	getUserByIDSQL = `
		SELECT u.email, u.pwd_hash, u.age
			FROM USERS u WHERE u.id = ? AND u.active = true
	`

	softDeleteUserSQL = `
		UPDATE USERS SET active = false
			WHERE id = ?	`
)

// UserRepositorier describes the methods used to do DB operations
type UserRepositorier interface {
	CreateUser(ctx context.Context, user entities.User) (int, error)
	Authenticate(ctx context.Context, session *entities.Session) (string, error)
	UpdateUser(ctx context.Context, information entities.Update) (entities.User, error)
	GetUser(ctx context.Context, id int) (entities.User, error)
	DeleteUser(ctx context.Context, id int) (bool, error)
}

// UserRepository implements the UserRepositorier interface
type UserRepository struct {
	db     *sql.DB
	logger log.Logger
}

// NewUserRepository returns a usreRepository pointer type
func NewUserRepository(mysqlDb *sql.DB, l log.Logger) *UserRepository {
	return &UserRepository{
		db:     mysqlDb,
		logger: log.With(l, "repository", "mysql"),
	}
}

// CreateUser does the DB operation to create a new user with the given data
func (r *UserRepository) CreateUser(ctx context.Context, user entities.User) (int, error) {
	id, err := r.db.ExecContext(ctx, createUserSQL, user.Email, user.Password, user.Age)

	if err != nil {
		return -1, errors.NewInternalError()
	}

	n, _ := id.LastInsertId()
	return int(n), nil
}

// Authenticate fetchs the password for the specific user and return it
func (r *UserRepository) Authenticate(ctx context.Context, session *entities.Session) (string, error) {
	var hash string

	err := r.db.QueryRow(authenticateSQL, session.Email).Scan(&hash)

	if err == sql.ErrNoRows {
		return "", errors.NewUserNotFoundError()
	}

	if err != nil {
		return "", errors.NewInternalError()
	}

	return hash, nil
}

// UpdateUser replaces the current information within DB with the given data
func (r *UserRepository) UpdateUser(ctx context.Context, update entities.Update) (entities.User, error) {
	var u entities.User

	if err := r.db.QueryRow(getUserByIDSQL, update.UserID).Scan(); err == sql.ErrNoRows {
		return u, errors.NewUserNotFoundError()
	}

	if _, err := r.db.ExecContext(ctx, updateUserSQL, update.Email, update.Password, update.Age, update.UserID); err != nil {
		return u, errors.NewInternalError()
	}

	_ = r.db.QueryRow(getUserByIDSQL, update.UserID).Scan(&u.Email, &u.Password, &u.Age)
	return u, nil
}

// GetUser fetchs the information within the DB about a specific user
func (r *UserRepository) GetUser(ctx context.Context, id int) (entities.User, error) {
	var u entities.User

	err := r.db.QueryRow(getUserByIDSQL, id).Scan(&u.Email, &u.Password, &u.Age)

	if err == sql.ErrNoRows {
		return entities.User{}, errors.NewUserNotFoundError()
	}

	if err != nil {
		return entities.User{}, errors.NewInternalError()
	}

	return u, nil
}

// DeleteUser does a soft delete operation over a specific user
func (r *UserRepository) DeleteUser(ctx context.Context, id int) (bool, error) {

	if err := r.db.QueryRow(getUserByIDSQL, id).Scan(); err == sql.ErrNoRows {
		return false, errors.NewUserNotFoundError()
	}

	if _, err := r.db.ExecContext(ctx, softDeleteUserSQL, id); err != nil {
		return false, errors.NewInternalError()
	}

	return true, nil
}
