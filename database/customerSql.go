package DB

import (
	"context"

	models "github.com/DeltaCapstone/ChoiceMoversBackend/models"
	"github.com/DeltaCapstone/ChoiceMoversBackend/utils"
	"github.com/jackc/pgx/v5"
)

// /////////////////////////////////////////////////////////////////////////////////////////////////
// Customer Route Queries
func (pg *postgres) GetCustomerById(ctx context.Context, id int) (models.Customer, error) {
	var customer models.Customer
	row := pg.db.QueryRow(ctx,
		`SELECT customer_id, username,first_name, last_name, 
		email, phone_primary FROM customers WHERE customer_id = $1`, id)

	if err := row.Scan(&customer.ID, &customer.UserName, &customer.FirstName, &customer.LastName, &customer.Email, &customer.PhonePrimary); err != nil {
		return customer, err
	}
	return customer, nil
}

func (pg *postgres) GetCustomerByUserName(ctx context.Context, userName string) (models.Customer, error) {
	var customer models.Customer
	row := pg.db.QueryRow(ctx,
		`SELECT customer_id, username,first_name, last_name, 
		email, phone_primary FROM customers WHERE username = $1`, userName)

	if err := row.Scan(&customer.ID, &customer.UserName, &customer.FirstName, &customer.LastName, &customer.Email, &customer.PhonePrimary); err != nil {
		return customer, err
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
WHERE customer_id = $8
`

func (pg *postgres) UpdateCustomer(ctx context.Context, updatedCustomer models.Customer) error {
	_, err := pg.db.Exec(ctx, updateCustomerQuery,
		updatedCustomer.UserName, updatedCustomer.PasswordHash,
		updatedCustomer.FirstName, updatedCustomer.LastName, updatedCustomer.Email,
		updatedCustomer.PhonePrimary, updatedCustomer.PhoneOther)

	return err
}
