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
		email, phone_primary, phone_other1, phone_other2 FROM customers WHERE username = $1`, userName)

	if err := row.Scan(
		&customer.UserName,
		&customer.FirstName,
		&customer.LastName,
		&customer.Email,
		&customer.PhonePrimary,
		&customer.PhoneOther1,
		&customer.PhoneOther2); err != nil {
		return models.GetCustomerResponse{}, err
	}
	return customer, nil
}

const createCustomerNameQuery = `INSERT INTO customers 
(username, password_hash, first_name, last_name, email, phone_primary, phone_other1,phone_other2) VALUES 
(@username,@password_hash,@first_name,@last_name,@email,@phone_primary,@phone_other1,@phone_other2) 
RETURNING username`

func (pg *postgres) CreateCustomer(ctx context.Context, newCustomer models.CreateCustomerParams) (string, error) {
	rows := pg.db.QueryRow(ctx, createCustomerNameQuery, pgx.NamedArgs(utils.StructToMap(newCustomer, "db")))
	var u string
	err := rows.Scan(&u)
	return u, err
}

const updateCustomerQuery = `
UPDATE customers
SET first_name = @first_name, last_name=@last_name , 
email=@email, phone_primary=@primary_phone, phone_other1 = @phone_other1, phone_other2 = @phone_other2
WHERE username = @username
`

func (pg *postgres) UpdateCustomer(ctx context.Context, updatedCustomer models.UpdateCustomerParams) error {
	_, err := pg.db.Exec(ctx, updateCustomerQuery, pgx.NamedArgs(utils.StructToMap(updatedCustomer, "db")))
	return err
}

const updateCustomerPasswordQuery = `
UPDATE customers
SET password_hash = @password_hash
WHERE username = @username
`

func (pg *postgres) UpdateCustomerPassword(ctx context.Context, username string, password_hash string) error {
	_, err := pg.db.Exec(ctx, updateCustomerPasswordQuery, pgx.NamedArgs{"username": username, "password_hash": password_hash})
	return err
}
