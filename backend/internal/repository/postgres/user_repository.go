package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Kineth-t/CS464-g1t10-project/internal/model"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(u model.User) (model.User, error) {
	row := r.db.QueryRow(context.Background(),
		`INSERT INTO users (username, password, phone_number, street, city, state, country, zip_code, role)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		 RETURNING id`,
		u.Username, u.Password, u.PhoneNumber,
		u.Address.Street, u.Address.City, u.Address.State, u.Address.Country, u.Address.ZipCode,
		u.Role,
	)
	if err := row.Scan(&u.ID); err != nil {
		return model.User{}, errors.New("username already taken")
	}
	return u, nil
}

func (r *UserRepository) FindByUsername(username string) (model.User, error) {
	var u model.User
	row := r.db.QueryRow(context.Background(),
		`SELECT id, username, password, phone_number, street, city, state, country, zip_code, role
		 FROM users WHERE username = $1`,
		username,
	)
	err := row.Scan(
		&u.ID, &u.Username, &u.Password, &u.PhoneNumber,
		&u.Address.Street, &u.Address.City, &u.Address.State, &u.Address.Country, &u.Address.ZipCode,
		&u.Role,
	)
	if err != nil {
		return model.User{}, errors.New("user not found")
	}
	return u, nil
}