package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Kineth-t/CS464-g1t10-project/internal/model"
)

type PhoneRepository struct {
	db *pgxpool.Pool
}

func NewPhoneRepository(db *pgxpool.Pool) *PhoneRepository {
	return &PhoneRepository{db: db}
}

func (r *PhoneRepository) GetAll() []model.Phone {
	rows, err := r.db.Query(context.Background(),
		`SELECT id, brand, model, price, stock, description FROM phones`)
	if err != nil {
		return nil
	}
	defer rows.Close()
	var phones []model.Phone
	for rows.Next() {
		var p model.Phone
		rows.Scan(&p.ID, &p.Brand, &p.Model, &p.Price, &p.Stock, &p.Description)
		phones = append(phones, p)
	}
	return phones
}

func (r *PhoneRepository) GetByID(id int) (model.Phone, error) {
	var p model.Phone
	row := r.db.QueryRow(context.Background(),
		`SELECT id, brand, model, price, stock, description FROM phones WHERE id = $1`, id)
	if err := row.Scan(&p.ID, &p.Brand, &p.Model, &p.Price, &p.Stock, &p.Description); err != nil {
		return model.Phone{}, errors.New("phone not found")
	}
	return p, nil
}

func (r *PhoneRepository) Create(p model.Phone) model.Phone {
	r.db.QueryRow(context.Background(),
		`INSERT INTO phones (brand, model, price, stock, description)
		 VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		p.Brand, p.Model, p.Price, p.Stock, p.Description,
	).Scan(&p.ID)
	return p
}

func (r *PhoneRepository) Update(p model.Phone) error {
	result, err := r.db.Exec(context.Background(),
		`UPDATE phones SET brand=$1, model=$2, price=$3, stock=$4, description=$5 WHERE id=$6`,
		p.Brand, p.Model, p.Price, p.Stock, p.Description, p.ID,
	)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return errors.New("phone not found")
	}
	return nil
}

func (r *PhoneRepository) Delete(id int) error {
	result, err := r.db.Exec(context.Background(),
		`DELETE FROM phones WHERE id=$1`, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return errors.New("phone not found")
	}
	return nil
}