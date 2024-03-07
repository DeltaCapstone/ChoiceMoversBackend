package main

import (
	//"errors"
	"errors"
	"fmt"
	"net/http"

	DB "github.com/DeltaCapstone/ChoiceMoversBackend/database"
	models "github.com/DeltaCapstone/ChoiceMoversBackend/models"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"

	//"github.com/jackc/pgerrcode"
	//"github.com/jackc/pgx/v5/pgconn"
	"github.com/labstack/echo/v4"
)

//TODO: Redo error handling to get rid of of al lthe sprintf's

func listJobs(c echo.Context) error {
	startDate := c.QueryParam("start")
	endDate := c.QueryParam("end")
	var status string
	if c.Get("role").(string) == "Manager" {
		status = "all"
	} else {
		status = "finalized"
	}

	jobs, err := DB.PgInstance.GetJobsByStatusAndRange(c.Request().Context(), status, startDate, endDate)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Error retrieving data: %v", err))
	}
	if jobs == nil {
		return c.String(http.StatusNotFound, fmt.Sprintf("No jobs found with status: %v", status))
	}
	return c.JSON(http.StatusOK, jobs)
}

func createAddress(address models.Address, c echo.Context) (int, error) {
	addr_id, err := DB.PgInstance.CreateAddress(c.Request().Context(), address)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case pgerrcode.UniqueViolation:
				fallthrough
			case pgerrcode.NotNullViolation:
				return 0, err
			}
		}
		return 0, err
	}
	return addr_id, nil
}

// Calculates how many hours a job should take
func jobHours(jobRequest models.CreateJobRequest) (pgtype.Interval, error) {

	return pgtype.Interval{}, nil
}

// Calculate the milage of a job
func jobMilage(jobRequest models.CreateJobRequest) (int, error) {

	return 0, nil
}

// Calculate the total cost of a job
func jobCost(jobRequest models.CreateJobRequest, hours pgtype.Interval, milage int) (string, error) {

	return "", nil
}

// Job POST Route to create a job
func createJobByCustomer(c echo.Context) error {
	var jobRequest models.CreateJobRequest
	// attempt at binding incoming json to a jobRequest
	if err := c.Bind(&jobRequest); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid job request data"})
	}

	loadAddrID, err := createAddress(jobRequest.LoadAddr, c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid job request data"})
	}

	unloadAddrID, err := createAddress(jobRequest.UnloadAddr, c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid job request data"})
	}

	// Calculate Labor Hours
	hours, err := jobHours(jobRequest)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid job request data"})
	}

	// Calculate the milage
	milage, err := jobMilage(jobRequest)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid job request data"})
	}

	// Calculate the cost of the job
	cost, err := jobCost(jobRequest, hours, milage)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid job request data"})
	}

	args := models.Job{
		Customer:   jobRequest.UserName,
		LoadAddr:   loadAddrID,
		UnloadAddr: unloadAddrID,
		StartTime:  jobRequest.StartTime,
		HoursLabor: hours,
		Finalized:  false,
		Rooms:      jobRequest.Rooms,
		Pack:       jobRequest.Pack,
		Unpack:     jobRequest.Unpack,
		Load:       jobRequest.Load,
		Unload:     jobRequest.Unload,
		Clean:      jobRequest.Clean,
		Milage:     milage,
		Cost:       cost,
	}

	job_id, err := DB.PgInstance.CreateJob(c.Request().Context(), args)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case pgerrcode.UniqueViolation:
				fallthrough
			case pgerrcode.NotNullViolation:
				return c.JSON(http.StatusConflict, fmt.Sprintf("Duplicate job: %v", err))
			}
		}
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("Failed to create job: %v", err))
	}

	return c.JSON(http.StatusCreated, echo.Map{"job id": job_id})
}
