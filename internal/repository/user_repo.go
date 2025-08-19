package repository

import (
	"GoRent/internal/domain/user"
	"GoRent/internal/repository/db"
	"context"
	"errors"
	"fmt"
)

type UserRepository interface {
	Create(ctx context.Context, u *user.User) error
	GetByEmail(ctx context.Context, email string) (*user.User, error)
	GetByID(ctx context.Context, id string) (*user.User, error)
	UpdateRole(ctx context.Context, id string, role user.Role) error
	GetAllUsers(ctx context.Context) ([]*user.User, error)
}

type userRepo struct {
	db *db.DB
}

func NewUserRepository(db *db.DB) UserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) Create(ctx context.Context, u *user.User) error {
	query := `INSERT INTO users (id, name, email, password_hash, role) VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.ExecContext(ctx, query, u.ID, u.Name, u.Email, u.PasswordHash, u.Role)
	return err
}

func (r *userRepo) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	query := `SELECT id, name, email, password_hash, role FROM users WHERE email = $1`
	row := r.db.QueryRowContext(ctx, query, email)

	var u user.User
	err := row.Scan(&u.ID, &u.Name, &u.Email, &u.PasswordHash, &u.Role)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (r *userRepo) GetByID(ctx context.Context, id string) (*user.User, error) {
	query := `SELECT id, name, email, password_hash, role FROM users WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, id)

	var u user.User
	err := row.Scan(&u.ID, &u.Name, &u.Email, &u.PasswordHash, &u.Role)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (r *userRepo) GetAllUsers(ctx context.Context) ([]*user.User, error) {
	query := `SELECT id, name, email, role FROM users`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("database query failed: %w", err)
	}
	defer rows.Close()

	var users []*user.User
	for rows.Next() {
		var u user.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Role); err != nil {
			return nil, fmt.Errorf("row scanning failed: %w", err)
		}
		users = append(users, &u)
	}

	return users, nil
}

func (r *userRepo) UpdateRole(ctx context.Context, id string, role user.Role) error {
	query := `UPDATE users SET role = $1 WHERE id = $2`
	result, err := r.db.ExecContext(ctx, query, role, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil
}
