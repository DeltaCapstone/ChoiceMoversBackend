package DB

import (
	"context"

	models "github.com/DeltaCapstone/ChoiceMoversBackend/models"
	"github.com/DeltaCapstone/ChoiceMoversBackend/utils"
	"github.com/jackc/pgx/v5"
)

// /////////////////////////////////////////////////////////////////////////////////////////////////////////
// Employee Route Queries
func (pg *postgres) GetEmployeeByUsername(ctx context.Context, username string) (models.GetEmployeeResponse, error) {
	var employee models.GetEmployeeResponse
	row := pg.db.QueryRow(ctx,
		`SELECT username, first_name, last_name, 
		email, phone_primary, employee_type FROM employees WHERE username = $1`, username)

	if err := row.Scan(&employee.UserName, &employee.FirstName, &employee.LastName,
		&employee.Email, &employee.PhonePrimary, &employee.EmployeeType); err != nil {
		return employee, err
	}
	return employee, nil
}

func (pg *postgres) GetEmployeeHashByUserName(ctx context.Context, userName string) (string, error) {
	var hash string
	row := pg.db.QueryRow(ctx,
		`SELECT password_hash FROM employees WHERE username = $1`, userName)

	if err := row.Scan(&hash); err != nil {
		return "", err
	}
	return hash, nil
}

func (pg *postgres) GetEmployeeRole(ctx context.Context, userName string) (string, error) {
	var role string
	row := pg.db.QueryRow(ctx,
		`SELECT employee_type FROM employees WHERE username = $1`, userName)

	if err := row.Scan(&role); err != nil {
		return "", err
	}
	return role, nil
}

func (pg *postgres) DeleteEmployeeByUsername(ctx context.Context, username string) error {
	_, err := pg.db.Exec(ctx, `DELETE FROM employees WHERE username = $1`, username)
	return err
}

func (pg *postgres) GetEmployeeList(ctx context.Context) ([]models.GetEmployeeResponse, error) {
	var employees []models.GetEmployeeResponse
	var rows pgx.Rows
	var err error

	rows, err = pg.db.Query(ctx,
		"SELECT username,first_name, last_name, email, phone_primary, employee_type FROM employees")

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var employee models.GetEmployeeResponse
		if err := rows.Scan(&employee.UserName, &employee.FirstName, &employee.LastName, &employee.Email, &employee.PhonePrimary, &employee.EmployeeType); err != nil {
			return nil, err
		}
		employees = append(employees, employee)
	}
	return employees, nil
}

const createEmployeeNameQuery = `INSERT INTO employees 
(username, password_hash, first_name, last_name, email, phone_primary, phone_other, employee_type) VALUES 
(@username,@password_hash,@first_name,@last_name,@email,@phone_primary,@phone_other,@employee_type) `

func (pg *postgres) CreateEmployee(ctx context.Context, newEmployee models.CreateEmployeeParams) error {
	_, err := pg.db.Exec(ctx, createEmployeeNameQuery, pgx.NamedArgs(utils.StructToMap(newEmployee, "db")))
	return err
}

const updateEmployeeQuery = `
UPDATE employees
SET username = @username, password_hash = @password_hash, first_name = @first_name, last_name = @last_name, email = @email, 
phone_primary = @phone_primary, phone_other = @phone_other, employee_type = @employee_type
WHERE username = @username`

func (pg *postgres) UpdateEmployee(ctx context.Context, updatedEmployee models.Employee) error {
	_, err := pg.db.Exec(ctx, updateEmployeeQuery, pgx.NamedArgs(utils.StructToMap(updatedEmployee, "db")))
	return err
}
