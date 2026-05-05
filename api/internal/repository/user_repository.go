package repository

import (
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"lakoo/backend/internal/domain"
)

type UserRepository interface {
	Create(user *domain.User) error
	GetByEmail(email string) (*domain.User, error)
	GetByID(id string) (*domain.User, error)
	UpdatePassword(id string, hashedPassword string) error
	UpdateProfile(user *domain.User) error
	GetByTenantID(tenantID string) ([]*domain.User, error)
	Delete(id string) error
}

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(u *domain.User) error {
	query := `
		INSERT INTO users (id, tenant_id, name, email, password, role, created_at, updated_at) 
		VALUES (:id, :tenant_id, :name, :email, :password, :role, :created_at, :updated_at)`
	
	_, err := r.db.NamedExec(query, u)
	return err
}

func (r *userRepository) GetByEmail(email string) (*domain.User, error) {
	var u domain.User
	err := r.db.Get(&u, "SELECT * FROM users WHERE email = ? AND deleted_at IS NULL", email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func (r *userRepository) GetByID(id string) (*domain.User, error) {
	var u domain.User
	err := r.db.Get(&u, "SELECT * FROM users WHERE id = ? AND deleted_at IS NULL", id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func (r *userRepository) UpdatePassword(id string, hashedPassword string) error {
	_, err := r.db.Exec("UPDATE users SET password = ? WHERE id = ?", hashedPassword, id)
	return err
}

func (r *userRepository) UpdateProfile(u *domain.User) error {
	query := `UPDATE users SET name = :name, email = :email, updated_at = :updated_at WHERE id = :id`
	_, err := r.db.NamedExec(query, u)
	return err
}

func (r *userRepository) GetByTenantID(tenantID string) ([]*domain.User, error) {
	var users []*domain.User
	err := r.db.Select(&users, "SELECT * FROM users WHERE tenant_id = ? AND deleted_at IS NULL", tenantID)
	return users, err
}

func (r *userRepository) Delete(id string) error {
	_, err := r.db.Exec("UPDATE users SET deleted_at = NOW() WHERE id = ?", id)
	return err
}
