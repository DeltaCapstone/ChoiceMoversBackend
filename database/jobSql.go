package DB

import (
	"context"

	"github.com/DeltaCapstone/ChoiceMoversBackend/models"
	"github.com/jackc/pgx/v5"
)

////////////////////////////////////////////////
//Jobs

const listJobsQuery = `SELECT 
estimate_id, customer_username,load_addr_id,unload_addr_id,start_time,end_time,rooms,special,
small_items,medium_items,large_items,boxes,item_load,flight_mult,pack,unpack,load,unload,
clean,need_truck,number_workers,
dist_to_job,dist_move,
estimated_man_hours,
estimated_rate,
estimated_cost,
job_id,
man_hours,
rate,
cost,
finalized,
actual_man_hours,
final_cost,
amount_payed,
notes   
FROM estimates NATURAL JOIN jobs WHERE start_time >= @start AND start_time <= @end`

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
		var ej models.EstimateJobJoin
		if err := scanStruct(rows, &ej); err != nil {
			return nil, err
		}
		var er models.EstimateResponse
		er.MakeFromJoin(ej)

		if ej.LoadAddrID != 0 {
			er.LoadAddr, _ = getAddr(ctx, ej.LoadAddrID)
		}
		if ej.UnloadAddrID != 0 {
			er.UnloadAddr, _ = getAddr(ctx, ej.UnloadAddrID)
		}
		er.Customer, _ = pg.GetCustomerByUserName(ctx, ej.CustomerUsername)
		var jr models.JobResponse
		jr.MakeFromJoin(ej)
		jr.EstimateResponse = er
		jr.AssignedEmp, _ = GetAssignedEmployees(ctx, jr.JobID)
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

const assignedEmpsQuery = `SELECT username, first_name,last_name,email,phone_primary,phone_other,employee_type, employee_priority, manager_override
	FROM employee_jobs JOIN employee ON employee_jobs.employee_username = employees.username WHERE job_id = $1`

func GetAssignedEmployees(ctx context.Context, jobID int) ([]models.AssignedEmployee, error) {
	var employees []models.AssignedEmployee
	var rows pgx.Rows
	var err error

	rows, err = PgInstance.db.Query(ctx, assignedEmpsQuery, jobID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var employee models.AssignedEmployee
		if err := rows.Scan(
			&employee.UserName,
			&employee.FirstName,
			&employee.LastName,
			&employee.Email,
			&employee.PhonePrimary,
			&employee.PhoneOther,
			&employee.EmployeeType,
			&employee.EmployeePriority,
			&employee.ManagerAssigned); err != nil {
			return nil, err
		}
		employees = append(employees, employee)
	}
	return employees, nil
}

const numWorkersQuery = "SELECT number_workers FROM estimates NATURAL JOIN jobs WHERE job_id = $1"

func (pg *postgres) GetNumWorksForJob(ctx context.Context, jobID int) (int, error) {
	row := pg.db.QueryRow(ctx, numWorkersQuery, jobID)
	var n int
	if err := row.Scan(&n); err != nil {
		return 0, err
	}
	return n, nil
}

const addEmployeeToJobQuery = `
INSERT INTO employee_jobs
(username, job_id,manager_override)
VALUES
(@username,@job_id,@manager_override)`

func (pg *postgres) AddEmployeeToJob(ctx context.Context, username string, jobId int, managerAssigned bool) error {
	_, err := pg.db.Exec(ctx, addEmployeeToJobQuery,
		pgx.NamedArgs{
			"username":         username,
			"job_id":           jobId,
			"manager_override": managerAssigned})
	return err
}

const removeEmployeeFromJobQuery = `
DELETE FROM employee_jobs
WHERE username=@username AND job_id=@job_id`

func (pg *postgres) RemoveEmployeeFromJob(ctx context.Context, username string, jobId int) error {
	_, err := pg.db.Exec(ctx, removeEmployeeFromJobQuery,
		pgx.NamedArgs{
			"username": username,
			"job_id":   jobId})
	return err
}
