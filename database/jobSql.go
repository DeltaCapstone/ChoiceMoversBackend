package DB

import (
	"context"
	"strconv"

	"github.com/DeltaCapstone/ChoiceMoversBackend/models"
	"github.com/DeltaCapstone/ChoiceMoversBackend/utils"
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
customer_notes,
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

const createEstimateQuery = `INSERT INTO estimates 
(customer_username, load_addr_id, unload_addr_id, start_time, end_time, rooms, special, small_items, medium_items, large_items, 
	boxes, item_load, flight_mult, pack, unpack, load, unload, clean, need_truck, number_workers, dist_to_job, dist_move,
	estimated_man_hours, estimated_rate, estimated_cost, customer_notes) VALUES 
(@customer_username, @load_addr_id, @unload_addr_id, @start_time, @end_time, @rooms, @special, @small_items, @medium_items, @large_items,
	@boxes, @item_load, @flight_mult, @pack, @unpack, @load, @unload, @clean, @need_truck, @number_workers, @dist_to_job, @dist_move,
	@estimated_man_hours, @estimated_rate, @estimated_cost, @customer_notes) RETURNING estimate_id`

func (pg *postgres) CreateEstimate(ctx context.Context, estimate models.Estimate) (int, error) {
	row := pg.db.QueryRow(ctx, createEstimateQuery, pgx.NamedArgs(utils.StructToMap(estimate, "db")))
	var u int
	err := row.Scan(&u)
	return u, err
}

const assignedEmpsQuery = `SELECT username, first_name,last_name,email,phone_primary,phone_other1,phone_other2,employee_type, employee_priority, manager_override
	FROM employee_jobs JOIN employees ON employee_jobs.employee_username = employees.username WHERE job_id = $1`

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
			&employee.PhoneOther1,
			&employee.PhoneOther2,
			&employee.EmployeeType,
			&employee.EmployeePriority,
			&employee.ManagerAssigned); err != nil {
			return nil, err
		}
		employees = append(employees, employee)
	}
	return employees, nil
}

const createAddress = `INSERT INTO addresses 
(street, city, state, zip, res_type, square_feet, flights, apt_num) VALUES 
(@street, @city, @state, @zip, @res_type, @square_feet, @flights, @apt_num) RETURNING address_id`

func (pg *postgres) CreateAddress(ctx context.Context, newJob models.Address) (int, error) {
	rows := pg.db.QueryRow(ctx, createAddress, pgx.NamedArgs(utils.StructToMap(newJob, "db")))
	var u string
	err := rows.Scan(&u)
	id, _ := strconv.Atoi(u)
	return id, err
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
(employee_username, job_id,manager_override)
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
WHERE employee_username=@username AND job_id=@job_id`

func (pg *postgres) RemoveEmployeeFromJob(ctx context.Context, username string, jobId int) error {
	_, err := pg.db.Exec(ctx, removeEmployeeFromJobQuery,
		pgx.NamedArgs{
			"username": username,
			"job_id":   jobId})
	return err
}

const getIsManagerAssignedQuery = `
SELECT manager_override FROM employee_jobs
WHERE employee_username=@username AND job_id=@job_id`

func (pg *postgres) GetIsManagerAssigned(ctx context.Context, username string, jobId int) (bool, error) {
	row := pg.db.QueryRow(ctx, getIsManagerAssignedQuery,
		pgx.NamedArgs{
			"username": username,
			"job_id":   jobId})
	var managerAssigned bool
	err := row.Scan(&managerAssigned)
	return managerAssigned, err
}

const getEstimateByIDQuery = `
SELECT customer_username, load_addr_id, unload_addr_id, start_time, end_time, rooms, special, small_items, medium_items, large_items, 
boxes, item_load, flight_mult, pack, unpack, load, unload, clean, need_truck, number_workers, dist_to_job, dist_move,
estimated_man_hours, estimated_rate, estimated_cost, customer_notes FROM estimates
WHERE estimate_id=$1`

func (pg *postgres) GetEstimateByID(ctx context.Context, estimateID int) (models.Estimate, error) {
	row := pg.db.QueryRow(ctx, getEstimateByIDQuery, estimateID)

	var res models.Estimate
	err := row.Scan(&res.CustomerUsername, &res.LoadAddrID, &res.UnloadAddrID, &res.StartTime, &res.EndTime, &res.Rooms, &res.Special, &res.Small, &res.Medium,
		&res.Large, &res.Boxes, &res.ItemLoad, &res.FlightMult, &res.Pack, &res.Unpack, &res.Load, &res.Unload, &res.Clean, &res.NeedTruck, &res.NumberWorkers,
		&res.DistToJob, &res.DistMove, &res.EstimateManHours, &res.EstimateRate, &res.EstimateCost, &res.CustomerNotes)
	return res, err
}

const createJobQuery = `INSERT INTO jobs 
(estimate_id, man_hours, rate, cost, finalized, final_cost, amount_payed, notes) VALUES 
(@estimate_id, @man_hours, @rate, @cost, @finalized, @final_cost, @amount_payed, @notes) RETURNING job_id`

func (pg *postgres) CreateJob(ctx context.Context, job models.Job) (int, error) {
	row := pg.db.QueryRow(ctx, createJobQuery, pgx.NamedArgs(utils.StructToMap(job, "db")))
	var u int
	err := row.Scan(&u)
	return u, err
}

const updateJobQuery = `UPDATE jobs 
SET man_hours = @man_hours, rate = @rate, cost = @cost, finalized = @finalized, 
final_cost = @final_cost, amount_payed = @amount_payed, notes = @notes
WHERE job_id=@job_id`

func (pg *postgres) UpdateJob(ctx context.Context, job models.Job) error {
	_, err := pg.db.Exec(ctx, updateJobQuery, pgx.NamedArgs(utils.StructToMap(job, "db")))
	return err
}

const getJobByIDQuery = `
SELECT estimate_id, man_hours, rate, cost, finalized, final_cost, amount_payed, notes FROM jobs
WHERE job_id=$1`

func (pg *postgres) GetJobByID(ctx context.Context, jobID int) (models.Job, error) {
	row := pg.db.QueryRow(ctx, getJobByIDQuery, jobID)

	var res models.Job
	err := row.Scan(&res.EstimateID, &res.ManHours, &res.Rate, &res.Cost, &res.Finalized, &res.FinalCost, &res.AmountPaid, &res.Notes)
	return res, err
}

const getJobsByUsernameQuery = `SELECT 
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
customer_notes,
notes   
FROM estimates NATURAL JOIN jobs WHERE customer_username=@username`

func (pg *postgres) GetJobsByUsername(ctx context.Context, username string) ([]models.JobResponse, error) {
	var jobs []models.JobResponse

	rows, err := pg.db.Query(ctx, getJobsByUsernameQuery, pgx.NamedArgs{"username": username})

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
		jobs = append(jobs, jr)
	}
	return jobs, nil
}
