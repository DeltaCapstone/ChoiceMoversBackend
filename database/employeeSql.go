package DB

import (
	"context"
	"fmt"
	"strconv"

	"github.com/DeltaCapstone/ChoiceMoversBackend/utils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

///////////////////////////////////////////////////////////////////////////////////////////////////////////
//Employee Route Queries

func (pg *postgres) GetEmployeeById(ctx context.Context, id string) (Employee, error) {
	var employee Employee
	var err error
	ID, er := strconv.Atoi(id)
	if er != nil {
		return employee, fmt.Errorf("id is not an integer: %v", err)
	}
	row := pg.db.QueryRow(ctx,
		`SELECT employee_id, username,first_name, last_name, 
		email, phone_primary FROM employees WHERE employee_id = $1`, ID)

	if err := row.Scan(&employee.ID, &employee.UserName, &employee.FirstName, &employee.LastName, &employee.Email, &employee.PhonePrimary); err != nil {
		return employee, fmt.Errorf("error reading row: %v", err)
	}

	return employee, nil
}

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

type CreateEmployeeParams struct {
	UserName     string        `db:"username" json:"userName"`
	PasswordHash string        `db:"password_hash" json:"passwordHash"`
	FirstName    string        `db:"first_name" json:"firstName"`
	LastName     string        `db:"last_name" json:"lastName"`
	Email        string        `db:"email" json:"email"`
	PhonePrimary pgtype.Text   `db:"phone_primary" json:"phonePrimary"`
	PhoneOther   []pgtype.Text `db:"phone_other" json:"phoneOther"`
	EmployeeType string        `db:"employee_type" json:"employeeType"`
}

const createEmployeeNameQuery = `INSERT INTO employees 
(username, password_hash, first_name, last_name, email, phone_primary, phone_other, employee_type) VALUES 
(@username,@password_hash,@first_name,@last_name,@email,@phone_primary,@phone_other,@employee_type) `

func (pg *postgres) CreateEmployee(ctx context.Context, newEmployee CreateEmployeeParams) (string, error) {
	row := pg.db.QueryRow(ctx, createEmployeeNameQuery, pgx.NamedArgs(utils.StructToMap(newEmployee, "db")))
	var u string
	err := row.Scan()
	return u, err
}

const updateEmployeeQuery = `
UPDATE employees
SET username = $1, password_hash = $2, first_name = $3, last_name = $4, email =$5, 
phone_primary = $6, phone_other = $7, employee_type = $8
WHERE employee_id = $9`

func (pg *postgres) UpdateEmployee(ctx context.Context, updatedEmployee Employee) error {
	_, err := pg.db.Exec(ctx, updateEmployeeQuery,
		updatedEmployee.UserName, updatedEmployee.PasswordHash,
		updatedEmployee.FirstName, updatedEmployee.LastName, updatedEmployee.Email,
		updatedEmployee.PhonePrimary, updatedEmployee.PhoneOther, updatedEmployee.EmployeeType)
	return err
}
