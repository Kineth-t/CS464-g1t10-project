package postgres

import (
	"context" // Used for DB operations
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Kineth-t/CS464-g1t10-project/internal/model"
)

// PhoneRepository handles all phone-related DB operations
type PhoneRepository struct {
	db *pgxpool.Pool // Connection pool to PostgreSQL
}

// Constructor
func NewPhoneRepository(db *pgxpool.Pool) *PhoneRepository {
	return &PhoneRepository{db: db}
}

// GetAll retrieves all phones from database
func (r *PhoneRepository) GetAll() []model.Phone {

	// Execute SELECT query
	rows, err := r.db.Query(context.Background(),
		`SELECT id, brand, model, price, stock, description FROM phones`)
	if err != nil {
		return nil // returns nil if query fails
	}
	defer rows.Close()

	var phones []model.Phone

	// Loop through result rows
	for rows.Next() {
		var p model.Phone

		// Map DB columns -> struct fields
		rows.Scan(&p.ID, &p.Brand, &p.Model, &p.Price, &p.Stock, &p.Description)

		phones = append(phones, p)
	}

	return phones
}

// GetByID retrieves a phone by its ID
func (r *PhoneRepository) GetByID(id int) (model.Phone, error) {
	var p model.Phone

	// Query a single row
	row := r.db.QueryRow(context.Background(),
		`SELECT id, brand, model, price, stock, description 
		 FROM phones WHERE id = $1`, id)

	// Scan result into struct
	if err := row.Scan(&p.ID, &p.Brand, &p.Model, &p.Price, &p.Stock, &p.Description); err != nil {
		return model.Phone{}, errors.New("phone not found")
	}

	return p, nil
}

// Create inserts a new phone into database
func (r *PhoneRepository) Create(p model.Phone) model.Phone {

	// Insert and return generated ID
	r.db.QueryRow(context.Background(),
		`INSERT INTO phones (brand, model, price, stock, description)
		 VALUES ($1, $2, $3, $4, $5) 
		 RETURNING id`,
		p.Brand, p.Model, p.Price, p.Stock, p.Description,
	).Scan(&p.ID)

	return p
}

// Update modifies an existing phone
func (r *PhoneRepository) Update(p model.Phone) error {

	// Execute UPDATE query
	result, err := r.db.Exec(context.Background(),
		`UPDATE phones 
		 SET brand=$1, model=$2, price=$3, stock=$4, description=$5 
		 WHERE id=$6`,
		p.Brand, p.Model, p.Price, p.Stock, p.Description, p.ID,
	)
	if err != nil {
		return err
	}

	// If no rows affected, phone does not exist
	if result.RowsAffected() == 0 {
		return errors.New("phone not found")
	}

	return nil
}

// Delete removes a phone by ID
func (r *PhoneRepository) Delete(id int) error {

	// Execute DELETE query
	result, err := r.db.Exec(context.Background(),
		`DELETE FROM phones WHERE id=$1`, id)
	if err != nil {
		return err
	}

	// If nothing deleted, phone not found
	if result.RowsAffected() == 0 {
		return errors.New("phone not found")
	}

	return nil
}