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
func (pg *postgres) GetCustomers(ctx context.Context, id string) ([]Customer, error) {
	var customers []Customer
	var rows pgx.Rows
	var err error
	if id != "" {
		ID, e := strconv.Atoi(id)
		if e != nil {
			return nil, fmt.Errorf("id is not an integer: %v", err)
		}
		rows, err = pg.db.Query(ctx, "SELECT customer_id, username,first_name, last_name, email, phone_primary FROM customers WHERE customer_id = $1", ID)
	} else {
		rows, err = pg.db.Query(ctx, "SELECT customer_id, username,first_name, last_name, email, phone_primary FROM customers")
	}
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

func (pg *postgres) CreateCustomer(ctx context.Context, newCustomer Customer) (int, error) {
	var newid int
	//chech if username or email exists

	//insert
	query := `INSERT INTO customers 
			(username, password_hash, first_name, last_name, email, phone_primary, phone_other) VALUES 
			(@username,@password_hash,@first_name,@last_name,@email,@phone_primary,@phone_other) 
			ON CONFLICT DO NOTHING RETURNING customer_id`

	pg.db.QueryRow(ctx, query, pgx.NamedArgs(structToMap(newCustomer, "db"))).Scan(&newid)
	if newid == 0 {
		return newid, fmt.Errorf("error inserting to database: could not create user")
	}
	return newid, nil
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////
//Employee Route Queries

func (pg *postgres) GetEmployees(ctx context.Context, id string) ([]Employee, error) {
	var employees []Employee
	var rows pgx.Rows
	var err error
	if id != "" {
		ID, e := strconv.Atoi(id)
		if e != nil {
			return nil, fmt.Errorf("id is not an integer: %v", err)
		}
		rows, err = pg.db.Query(ctx, "SELECT employee_id, username,first_name, last_name, email, phone_primary, employee_type FROM employees FROM employee_id = $1", ID)
	} else {
		rows, err = pg.db.Query(ctx, "SELECT employee_id, username,first_name, last_name, email, phone_primary, employee_type FROM employees")
	}
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

func (pg *postgres) CreateEmployee(ctx context.Context, newEmployee Employee) (int, error) {
	var newid int
	//check if username or email exits

	//insert
	query := `INSERT INTO employees 
			(username, password_hash, first_name, last_name, email, phone_primary, phone_other, employee_type) VALUES 
			(@username,@password_hash,@first_name,@last_name,@email,@phone_primary,@phone_other,@employee_type) 
			ON CONFLICT DO NOTHING RETURNING employee_id `
	pg.db.QueryRow(ctx, query, pgx.NamedArgs(structToMap(newEmployee, "db"))).Scan(&newid)
	if newid == 0 {
		return newid, fmt.Errorf("error inserting to database: could not create user")
	}
	return newid, nil
}
