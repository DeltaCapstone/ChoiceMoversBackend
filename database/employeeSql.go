package DB

import (
	"context"

	models "github.com/DeltaCapstone/ChoiceMoversBackend/models"
	"github.com/DeltaCapstone/ChoiceMoversBackend/utils"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// /////////////////////////////////////////////////////////////////////////////////////////////////////////
// Employee Route Queries
func (pg *postgres) GetEmployeeCredentials(ctx context.Context, userName string) (string, error) {
	var hash string
	row := pg.db.QueryRow(ctx,
		`SELECT password_hash FROM employees WHERE username = $1`, userName)

	if err := row.Scan(&hash); err != nil {
		return "", err
	}
	return hash, nil
}

func (pg *postgres) GetEmployeeByUsername(ctx context.Context, username string) (models.GetEmployeeResponse, error) {
	var employee models.GetEmployeeResponse
	row := pg.db.QueryRow(ctx,
		`SELECT username, first_name, last_name, 
		email, phone_primary,phone_other1,phone_other2, employee_type, employee_priority FROM employees WHERE username = $1`, username)

	if err := row.Scan(
		&employee.UserName,
		&employee.FirstName,
		&employee.LastName,
		&employee.Email,
		&employee.PhonePrimary,
		&employee.PhoneOther1,
		&employee.PhoneOther2,
		&employee.EmployeeType,
		&employee.EmployeePriority); err != nil {
		return employee, err
	}
	return employee, nil
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
		"SELECT username,first_name, last_name, email, phone_primary,phone_other1,phone_other2, employee_type, employee_priority FROM employees")

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var employee models.GetEmployeeResponse
		if err := rows.Scan(
			&employee.UserName,
			&employee.FirstName,
			&employee.LastName,
			&employee.Email,
			&employee.PhonePrimary,
			&employee.PhoneOther1,
			&employee.PhoneOther2,
			&employee.EmployeeType,
			&employee.EmployeePriority); err != nil {
			return nil, err
		}
		employees = append(employees, employee)
	}
	return employees, nil
}

const addSignupQuery = `INSERT INTO employee_signup (id,email,employee_type,employee_priority,signup_token, expires_at, used)
VALUES
(@id, @email, @employee_type, @employee_priority, @signup_token, @expires_at, @used)`

func (pg *postgres) AddEmployeeSignup(ctx context.Context, newEmployeeSignUp models.EmployeeSignup) error {
	_, err := pg.db.Exec(ctx, addSignupQuery, pgx.NamedArgs(utils.StructToMap(newEmployeeSignUp, "db")))
	return err
}

const getSignupQuery = `SELECT id,email,employee_type,employee_priority,signup_token, expires_at, used 
FROM employee_signup where id=$1`

func (pg *postgres) GetEmployeeSignup(ctx context.Context, id uuid.UUID) (models.EmployeeSignup, error) {
	row := pg.db.QueryRow(ctx, getSignupQuery, id)
	var es models.EmployeeSignup
	if err := scanStruct(row, &es); err != nil {
		return es, err
	}
	return es, nil
}

func (pg *postgres) UseEmployeeSignup(ctx context.Context, id uuid.UUID) error {
	_, err := pg.db.Exec(ctx, "UPDATE employee_signup SET used=true WHERE id=$1", id)
	return err
}

const createEmployeeNameQuery = `INSERT INTO employees 
(username, password_hash, first_name, last_name, email,
	phone_primary, phone_other1,phone_other2, employee_type,employee_priority) 
VALUES 
(@username,@password_hash,@first_name,@last_name,@email,
	@phone_primary,@phone_other1,phone_other2,@employee_type,@employee_priority) `

func (pg *postgres) CreateEmployee(ctx context.Context, newEmployee models.CreateEmployeeParams) error {
	_, err := pg.db.Exec(ctx, createEmployeeNameQuery, pgx.NamedArgs(utils.StructToMap(newEmployee, "db")))
	return err
}

const updateEmployeeQuery = `
UPDATE employees
SET first_name = @first_name, last_name = @last_name, email = @email, 
phone_primary = @phone_primary, phone_other1 = @phone_other1, phone_other2=@phone_other2
WHERE username = @username`

func (pg *postgres) UpdateEmployee(ctx context.Context, updatedEmployee models.UpdateEmployeeParams) error {
	_, err := pg.db.Exec(ctx, updateEmployeeQuery, pgx.NamedArgs(utils.StructToMap(updatedEmployee, "db")))
	return err
}

const updateEmployeePasswordQuery = `
UPDATE employees
SET password_hash = @password_hash
WHERE username = @username
`

func (pg *postgres) UpdateEmployeePassword(ctx context.Context, username string, password_hash string) error {
	_, err := pg.db.Exec(ctx, updateEmployeePasswordQuery, pgx.NamedArgs{"username": username, "password_hash": password_hash})
	return err
}

const updateEmployeeTypePriorityQuery = `
UPDATE employees
SET employee_type = @employee_type,
employee_priority = @employee_priority
where username = @username`

func (pg *postgres) UpdateEmployeeTypePriority(ctx context.Context, update models.UpdateEmployeeTypePriorityParams) error {
	_, err := pg.db.Exec(ctx, updateEmployeeTypePriorityQuery,
		pgx.NamedArgs{
			"username":          update.UserName,
			"employee_type":     update.EmployeeType,
			"employee_priority": update.EmployeePriority})
	return err
}

const getEmployeePriorityQuery = "SELECT employee_priority FROM employees WHERE username = $1"

func (pg *postgres) GetEmployeePriority(ctx context.Context, username string) (int, error) {
	row := pg.db.QueryRow(ctx, getEmployeePriorityQuery, username)
	var p int
	if err := row.Scan(&p); err != nil {
		return -1, err
	}
	return p, nil
}
