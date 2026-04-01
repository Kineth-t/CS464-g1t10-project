package postgres

import (
	"context" // Used for DB operations
	"errors"
	"fmt"

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
func (r *PhoneRepository) GetAll() ([]model.Phone, error) {
	// Execute SQL query to fetch all phones
	rows, err := r.db.Query(context.Background(),
		`SELECT id, brand, model, price, stock, description, image_url
		 FROM phones
		 ORDER BY id ASC`)
	if err != nil {
		return nil, fmt.Errorf("failed to query phones: %w", err)
	}
	defer rows.Close()

	var phones []model.Phone

	// Iterate through result set
	for rows.Next() {
		var p model.Phone

		// Scan row into struct
		if err := rows.Scan(
			&p.ID,
			&p.Brand,
			&p.Model,
			&p.Price,
			&p.Stock,
			&p.Description,
			&p.ImageURL,
		); err != nil {
			return nil, fmt.Errorf("failed to scan phone row: %w", err)
		}

		phones = append(phones, p)
	}

	// Check for errors that occurred during iteration
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error while iterating rows: %w", err)
	}

	return phones, nil
}

// GetByID retrieves a phone by its ID
func (r *PhoneRepository) GetByID(id int) (model.Phone, error) {
	var p model.Phone

	// Query a single row
	row := r.db.QueryRow(context.Background(),
		`SELECT id, brand, model, price, stock, description, image_url
		 FROM phones WHERE id = $1`, id)

	// Scan result into struct
	if err := row.Scan(&p.ID, &p.Brand, &p.Model, &p.Price, &p.Stock, &p.Description, &p.ImageURL); err != nil {
		return model.Phone{}, errors.New("phone not found")
	}

	return p, nil
}

// CheckStockAndReserve checks stock inside a transaction with FOR UPDATE to prevent
// two users from adding the last item to their carts at the same time.
// Returns the current price so the cart item always captures the price at time of add.
func (r *PhoneRepository) CheckStockAndReserve(phoneID, quantity int) (float64, error) {
	// Begin transaction
	tx, err := r.db.Begin(context.Background())
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(context.Background())

	// Lock the phone row — concurrent adds to cart for the same phone will queue here
	var stock int
	var price float64
	err = tx.QueryRow(context.Background(),
		`SELECT stock, price FROM phones WHERE id=$1 FOR UPDATE`, phoneID,
	).Scan(&stock, &price)
	if err != nil {
		return 0, errors.New("phone not found")
	}

	// Validate stock is sufficient
	if stock < quantity {
		return 0, errors.New("insufficient stock")
	}

	// Commit — releases the lock without changing any data.
	// Stock is only deducted at checkout, not when adding to cart.
	return price, tx.Commit(context.Background())
}

// Create inserts a new phone into database
func (r *PhoneRepository) Create(p model.Phone) model.Phone {

	// Insert and return generated ID
	r.db.QueryRow(context.Background(),
		`INSERT INTO phones (brand, model, price, stock, description, image_url)
		 VALUES ($1, $2, $3, $4, $5, $6) 
		 RETURNING id`,
		p.Brand, p.Model, p.Price, p.Stock, p.Description, p.ImageURL,
	).Scan(&p.ID)

	return p
}

// Update modifies an existing phone
func (r *PhoneRepository) Update(p model.Phone) error {

	// Execute UPDATE query
	result, err := r.db.Exec(context.Background(),
		`UPDATE phones 
		 SET brand=$1, model=$2, price=$3, stock=$4, description=$5, image_url=$6 
		 WHERE id=$7`,
		p.Brand, p.Model, p.Price, p.Stock, p.Description, p.ImageURL, p.ID,
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
