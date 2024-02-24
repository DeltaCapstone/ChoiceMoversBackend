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
		rows, err = pg.db.Query(ctx, "select customer_id, username, email, phone_primary from customers where id = $1", ID)
	} else {
		rows, err = pg.db.Query(ctx, "select customer_id, username, email, phone_primary from customers")
	}
	if err != nil {
		return nil, fmt.Errorf("error querying database: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var customer Customer
		if err := rows.Scan(&customer.ID, &customer.UserName, &customer.Email, &customer.PhonePrimary); err != nil {
			return nil, fmt.Errorf("error reading row: %v", err)
		}
		customers = append(customers, customer)
	}
	return customers, nil
}

func (pg *postgres) CreateCustomer(ctx context.Context, newCustomer Customer) (int, error) {
	user_id := 0
	return user_id, nil
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
		rows, err = pg.db.Query(ctx, "select employee_id, username, email, phone_primary, employee_type from employees where id = $1", ID)
	} else {
		rows, err = pg.db.Query(ctx, "select employee_id, username, email, phone_primary, employee_type from employees")
	}
	if err != nil {
		return nil, fmt.Errorf("error querying database: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var employee Employee
		if err := rows.Scan(&employee.ID, &employee.UserName, &employee.Email, &employee.PhonePrimary, &employee.EmployeeType); err != nil {
			return nil, fmt.Errorf("error reading row: %v", err)
		}
		employees = append(employees, employee)
	}
	return employees, nil
}

func (pg *postgres) CreateEmployee(ctx context.Context, newEmployee Employee) (int, error) {
	user_id := 0
	return user_id, nil
}
