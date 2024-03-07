package DB

import (
	"context"

	models "github.com/DeltaCapstone/ChoiceMoversBackend/models"
	"github.com/DeltaCapstone/ChoiceMoversBackend/utils"
	"github.com/jackc/pgx/v5"
)

// /////////////////////////////////////////////////////////////////////////////////////////////////
// Customer Route Queries

func (pg *postgres) GetCustomerCredentials(ctx context.Context, userName string) (string, error) {
	var hash string
	row := pg.db.QueryRow(ctx,
		`SELECT password_hash FROM customers WHERE username = $1`, userName)

	if err := row.Scan(&hash); err != nil {
		return "", err
	}
	return hash, nil
}

func (pg *postgres) GetCustomerByUserName(ctx context.Context, userName string) (models.GetCustomerResponse, error) {
	var customer models.GetCustomerResponse
	row := pg.db.QueryRow(ctx,
		`SELECT username,first_name, last_name, 
		email, phone_primary, phone_other FROM customers WHERE username = $1`, userName)

	if err := row.Scan(
		&customer.UserName,
		&customer.FirstName,
		&customer.LastName,
		&customer.Email,
		&customer.PhonePrimary,
		&customer.PhoneOther); err != nil {
		return models.GetCustomerResponse{}, err
	}
	return customer, nil
}

const createCustomerNameQuery = `INSERT INTO customers 
(username, password_hash, first_name, last_name, email, phone_primary, phone_other) VALUES 
(@username,@password_hash,@first_name,@last_name,@email,@phone_primary,@phone_other) 
RETURNING username`

func (pg *postgres) CreateCustomer(ctx context.Context, newCustomer models.CreateCustomerParams) (string, error) {
	rows := pg.db.QueryRow(ctx, createCustomerNameQuery, pgx.NamedArgs(utils.StructToMap(newCustomer, "db")))
	var u string
	err := rows.Scan(&u)
	return u, err
}

const updateCustomerQuery = `
UPDATE customers
SET username = $1, password_hash = $2, first_name = $3, last_name = $4, email =$5, phone_primary = $6, phone_other = $7
WHERE username = $1
`

func (pg *postgres) UpdateCustomer(ctx context.Context, updatedCustomer models.Customer) error {
	_, err := pg.db.Exec(ctx, updateCustomerQuery,
		updatedCustomer.UserName, updatedCustomer.PasswordHash,
		updatedCustomer.FirstName, updatedCustomer.LastName, updatedCustomer.Email,
		updatedCustomer.PhonePrimary, updatedCustomer.PhoneOther)

	return err
}
