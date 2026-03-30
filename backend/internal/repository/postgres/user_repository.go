package postgres

import (
	"context" // Used for database operations
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Kineth-t/CS464-g1t10-project/internal/model"
)

// UserRepository handles user-related database operations
type UserRepository struct {
	db *pgxpool.Pool // PostgreSQL connection pool
}

// Constructor
func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

// Create inserts a new user into the database
func (r *UserRepository) Create(u model.User) (model.User, error) {

	// Insert user into DB and return generated ID
	row := r.db.QueryRow(context.Background(),
		`INSERT INTO users (
			username, password, phone_number,
			street, city, state, country, zip_code, role
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id`,
		u.Username,
		u.Password, // hashed before storing
		u.PhoneNumber,
		u.Address.Street,
		u.Address.City,
		u.Address.State,
		u.Address.Country,
		u.Address.ZipCode,
		u.Role,
	)

	// Scan returned ID into struct
	if err := row.Scan(&u.ID); err != nil {
		// Likely caused by UNIQUE constraint (duplicate username)
		return model.User{}, errors.New("username already taken")
	}

	return u, nil
}

// FindByUsername retrieves a user by username
func (r *UserRepository) FindByUsername(username string) (model.User, error) {
	var u model.User

	// Query user by username
	row := r.db.QueryRow(context.Background(),
		`SELECT 
			id, username, password, phone_number,
			street, city, state, country, zip_code, role
		 FROM users 
		 WHERE username = $1`,
		username,
	)

	// Map DB fields -> struct fields
	err := row.Scan(
		&u.ID,
		&u.Username,
		&u.Password,
		&u.PhoneNumber,
		&u.Address.Street,
		&u.Address.City,
		&u.Address.State,
		&u.Address.Country,
		&u.Address.ZipCode,
		&u.Role,
	)

	if err != nil {
		// If no user found (or other error)
		return model.User{}, errors.New("user not found")
	}

	return u, nil
}