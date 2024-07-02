package repository

import (
	"context"
	"database/sql"
	"log"

	"backend-app/models"
)

type UserRepository struct {
	DB *sql.DB
}

// NewUserRepository initializes a new UserRepository
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		DB: db,
	}
}

// CreateUser creates a new user in the database
func (ur *UserRepository) CreateUser(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (username, email, password_hash)
		VALUES ($1, $2, $3)
		RETURNING id
	`
	err := ur.DB.QueryRowContext(
		ctx,
		query,
		user.Username,
		user.Email,
		user.PasswordHash,
	).Scan(&user.ID)
	if err != nil {
		log.Println("Error creating user:", err)
		return err
	}
	return nil
}

// UpdateUser updates an existing user in the database
func (ur *UserRepository) UpdateUser(user *models.User) error {
	query := `
		UPDATE users
		SET username = $1, email = $2, password_hash = $3
		WHERE id = $4
	`
	_, err := ur.DB.Exec(
		query,
		user.Username,
		user.Email,
		user.PasswordHash,
		user.ID,
	)
	if err != nil {
		log.Println("Error updating user:", err)
		return err
	}
	return nil
}

// DeleteUser deletes a user from the database by ID
func (ur *UserRepository) DeleteUser(userID int64) error {
	query := `
		DELETE FROM users
		WHERE id = $1
	`
	_, err := ur.DB.Exec(query, userID)
	if err != nil {
		log.Println("Error deleting user:", err)
		return err
	}
	return nil
}

// GetUserByID retrieves a user from the database by ID
func (ur *UserRepository) GetUserByID(userID int64) (*models.User, error) {
	var user models.User
	query := `
		SELECT id, username, email, password_hash
		FROM users
		WHERE id = $1
	`
	err := ur.DB.QueryRow(query, userID).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
	)
	if err != nil {
		log.Println("Error retrieving user:", err)
		return nil, err
	}
	return &user, nil
}

// GetUserByEmail retrieves a user from the database by email
func (ur *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	query := `
		SELECT id, username, email, password_hash
		FROM users
		WHERE email = $1
	`
	err := ur.DB.QueryRow(query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
	)
	if err != nil {
		log.Println("Error retrieving user by email:", err)
		return nil, err
	}
	return &user, nil
}
