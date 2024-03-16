package DB

import (
	"context"

	"github.com/DeltaCapstone/ChoiceMoversBackend/models"
	"github.com/jackc/pgx/v5"
)

////////////////////////////////////////////////
//Jobs

const listJobsQuery = `SELECT * FROM jobs WHERE start_time >= @start AND start_time <= @end`

// TODO: Figure out error handling for address errors
func (pg *postgres) GetJobsByStatusAndRange(ctx context.Context, status string, start string, end string) ([]models.JobResponse, error) {
	var jobs []models.JobResponse
	var query string
	switch status {
	case "all":
		query = listJobsQuery
	case "finalized":
		query = listJobsQuery + " AND finalized = true"
	}

	rows, err := pg.db.Query(ctx, query, pgx.NamedArgs{"start": start, "end": end})

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var j models.Job
		if err := scanStructfromRows(rows, &j); err != nil {
			return nil, err
		}
		er, err := pg.getEstimateForJob(ctx, j.EstimateID)
		if err != nil {
			return nil, err
		}
		var jr models.JobResponse
		jr.MakeFromJob(j)
		jr.EstimateResponse = er
		jr.AssignedEmp, _ = getAssignedEmployees(ctx, jr.JobID)
		jobs = append(jobs, jr)
	}
	return jobs, nil
}

const addrQuery = "SELECT * FROM addresses WHERE address_id = $1"

func getAddr(ctx context.Context, addrID int) (models.Address, error) {
	var a models.Address
	row := PgInstance.db.QueryRow(ctx, addrQuery, addrID)
	err := row.Scan(
		&a.AddressID,
		&a.Street,
		&a.City,
		&a.State,
		&a.Zip,
		&a.ResType,
		&a.SquareFeet,
		&a.Flights,
		&a.AptNum,
	)
	if err != nil {
		return a, err
	}
	return a, nil
}

const assignedEmpsQuery = `SELECT username, first_name,last_name,email,phone_primary,phone_other,employee_type, employee_priority
	FROM employee_jobs JOIN employee ON employee_jobs.employee_username = employees.username WHERE job_id = $1`

func getAssignedEmployees(ctx context.Context, jobID int) ([]models.GetEmployeeResponse, error) {
	var employees []models.GetEmployeeResponse
	var rows pgx.Rows
	var err error

	rows, err = PgInstance.db.Query(ctx, assignedEmpsQuery, jobID)

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
			&employee.PhoneOther,
			&employee.EmployeeType,
			&employee.EmployeePriority); err != nil {
			return nil, err
		}
		employees = append(employees, employee)
	}
	return employees, nil
}

const getEstimateByIDQuery = `SELECT * FROM estimates where estimate_id = $1`

func (pg *postgres) getEstimateForJob(ctx context.Context, estId int) (models.EstimateResponse, error) {
	var (
		LoadAddrID   int
		UnloadAddrID int
		Customer     string
		er           models.EstimateResponse
	)
	row := pg.db.QueryRow(ctx, getEstimateByIDQuery, estId)

	var e models.Estimate
	if err := scanStruct(row, &e); err != nil {
		return er, err
	}

	er.MakeFromEstimate(e)
	if LoadAddrID != 0 {
		er.LoadAddr, _ = getAddr(ctx, e.LoadAddrID)
	}
	if UnloadAddrID != 0 {
		er.UnloadAddr, _ = getAddr(ctx, e.UnloadAddrID)
	}
	er.Customer, _ = pg.GetCustomerByUserName(ctx, Customer)

	return er, nil
}
