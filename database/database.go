package DB

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type postgres struct {
	db *pgxpool.Pool
}

var (
	PgInstance *postgres
	pgOnce     sync.Once
)

func NewPG(ctx context.Context) (*postgres, error) {
	pgOnce.Do(func() {
		config, err := pgxpool.ParseConfig(fmt.Sprintf("user=%s password=%s dbname=%s host=db port=5432 sslmode=disable", os.Getenv("PGUSER"), os.Getenv("PGPASSWORD"), os.Getenv("PGDATABASE")))
		if err != nil {
			log.Fatalf("unable to parse PostgreSQL configuration: %v", err)
		}

		db, err := pgxpool.NewWithConfig(ctx, config)
		if err != nil {
			// Check if the error is a PgError
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) {
				log.Fatalf("PostgreSQL error - Code: %s, Message: %s", pgErr.Code, pgErr.Message)
			} else {
				// If it's not a PgError, log the general error
				log.Fatalf("unable to create connection pool: %v", err)
			}
		}
		PgInstance = &postgres{db}
	})

	return PgInstance, nil
}

func (pg *postgres) Ping(ctx context.Context) error {
	return pg.db.Ping(ctx)
}

func (pg *postgres) Close() {
	pg.db.Close()
}

// /////////////////////////////////////////////////////////////////////////////////////////////////
// Customer Route Queries
func (pg *postgres) GetCustomerById(ctx context.Context, id string) ([]Customer, error) {
	var customers []Customer
	var rows pgx.Rows
	var err error
	ID, er := strconv.Atoi(id)
	if er != nil {
		return nil, fmt.Errorf("id is not an integer: %v", err)
	}
	rows, err = pg.db.Query(ctx, "SELECT customer_id, username,first_name, last_name, email, phone_primary FROM customers WHERE customer_id = $1", ID)

	if err != nil {
		return nil, fmt.Errorf("error querying database: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var customer Customer
		if err := rows.Scan(&customer.ID, &customer.UserName, &customer.FirstName, &customer.LastName, &customer.Email, &customer.PhonePrimary); err != nil {
			return nil, fmt.Errorf("error reading row: %v", err)
		}
		customers = append(customers, customer)
	}
	return customers, nil
}

const createCustomerNameQuery = `INSERT INTO customers 
(username, password_hash, first_name, last_name, email, phone_primary, phone_other) VALUES 
(@username,@password_hash,@first_name,@last_name,@email,@phone_primary,@phone_other) 
RETURNING customer_id`

func (pg *postgres) CreateCustomer(ctx context.Context, newCustomer Customer) (int, error) {
	var newid int
	err := pg.db.QueryRow(ctx, createCustomerNameQuery, pgx.NamedArgs(structToMap(newCustomer, "db"))).Scan(&newid)
	return newid, err
}

const updateCustomerQuery = `
UPDATE customers
SET username = $1, password_hash = $2, first_name = $3, last_name = $4, email =$5, phone_primary = $6, phone_other = $7
WHERE customer_id = $8
`

func (pg *postgres) UpdateCustomer(ctx context.Context, updatedCustomer Customer) error {
	_, err := pg.db.Exec(ctx, updateCustomerQuery,
		updatedCustomer.UserName, updatedCustomer.PasswordHash,
		updatedCustomer.FirstName, updatedCustomer.LastName, updatedCustomer.Email,
		updatedCustomer.PhonePrimary, updatedCustomer.PhoneOther, updatedCustomer.ID)

	return err
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////
//Employee Route Queries

func (pg *postgres) GetEmployeeList(ctx context.Context) ([]Employee, error) {
	var employees []Employee
	var rows pgx.Rows
	var err error

	rows, err = pg.db.Query(ctx,
		"SELECT employee_id, username,first_name, last_name, email, phone_primary, employee_type FROM employees")

	if err != nil {
		return nil, fmt.Errorf("error querying database: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var employee Employee
		if err := rows.Scan(&employee.ID, &employee.UserName, &employee.FirstName, &employee.LastName, &employee.Email, &employee.PhonePrimary, &employee.EmployeeType); err != nil {
			return nil, fmt.Errorf("error reading row: %v", err)
		}
		employees = append(employees, employee)
	}
	return employees, nil
}

const createEmployeeNameQuery = `INSERT INTO employees 
(username, password_hash, first_name, last_name, email, phone_primary, phone_other, employee_type) VALUES 
(@username,@password_hash,@first_name,@last_name,@email,@phone_primary,@phone_other,@employee_type) 
RETURNING employee_id `

func (pg *postgres) CreateEmployee(ctx context.Context, newEmployee Employee) (int, error) {
	var newid int
	err := pg.db.QueryRow(ctx, createEmployeeNameQuery, pgx.NamedArgs(structToMap(newEmployee, "db"))).Scan(&newid)
	return newid, err
}

const updateEmployeeQuery = `
UPDATE employees
SET username = $1, password_hash = $2, first_name = $3, last_name = $4, email =$5, phone_primary = $6, phone_other = $7, employee_type = $8
WHERE employee_id = $9`

func (pg *postgres) UpdateEmployee(ctx context.Context, updatedEmployee Employee) error {
	_, err := pg.db.Exec(ctx, updateEmployeeQuery,
		updatedEmployee.UserName, updatedEmployee.PasswordHash,
		updatedEmployee.FirstName, updatedEmployee.LastName, updatedEmployee.Email,
		updatedEmployee.PhonePrimary, updatedEmployee.PhoneOther, updatedEmployee.EmployeeType, updatedEmployee.ID)
	return err
}

////////////////////////////////////////////////
//Jobs

// shell
func (pg *postgres) GetJobsByStatus(ctx context.Context, status string) ([]Job, error) {
	var jobs []Job
	return jobs, nil
}
