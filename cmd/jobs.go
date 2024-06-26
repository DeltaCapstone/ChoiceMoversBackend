package main

import (
	//"errors"

	"errors"
	"fmt"
	"net/http"

	DB "github.com/DeltaCapstone/ChoiceMoversBackend/database"
	"github.com/DeltaCapstone/ChoiceMoversBackend/models"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"

	//"github.com/jackc/pgerrcode"
	//"github.com/jackc/pgx/v5/pgconn"
	"github.com/labstack/echo/v4"
)

//TODO: Redo error handling to get rid of of al lthe sprintf's

// Lists jobs between a certain start and end date and based on status
// Query params:
// start - start date of the search
// end - end date of the search
func listJobs(c echo.Context) error {
	startDate := c.QueryParam("start")
	endDate := c.QueryParam("end")
	var status string

	// Managers should see all jobs, employees should only see finalized jobs
	if c.Get("role").(string) == "Manager" {
		status = "all"
	} else {
		status = "finalized"
	}

	// Query DB for jobs fitting criteria
	jobs, err := DB.PgInstance.GetJobsByStatusAndRange(c.Request().Context(), status, startDate, endDate)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Error retrieving data: %v", err))
	}
	if jobs == nil {
		return c.String(http.StatusNotFound, fmt.Sprintf("No jobs found with status: %v", status))
	}
	return c.JSON(http.StatusOK, jobs)
}

// Creates a new job in the database that pulls data from a given estimate
// Also links job to estimate via estimate id
func convertEstimateToJob(c echo.Context) error {
	var req models.ConvertEstimateToJob
	// attempt at binding incoming json to an estimate request
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err})
	}

	res, err := DB.PgInstance.GetEstimateByID(c.Request().Context(), req.EstimateID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err})
	}

	username := c.Get("username").(string)

	// Since this is a customer route, customers should only be able to turn their own estimates into jobs
	if username != res.CustomerUsername {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "not correct user"})
	}

	job := models.Job{
		EstimateID: req.EstimateID,
		ManHours:   res.EstimateManHours,
		Rate:       res.EstimateRate,
		Cost:       res.EstimateCost,

		Finalized:  false,
		FinalCost:  0,
		AmountPaid: 0,
		Notes:      res.CustomerNotes,
	}

	// Creates the job in the database
	id, err := DB.PgInstance.CreateJob(c.Request().Context(), job)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case pgerrcode.UniqueViolation:
				fallthrough
			case pgerrcode.NotNullViolation:
				return c.JSON(http.StatusConflict, fmt.Sprintf("Not Null violation: %v ----- Data: %v", err, job))
			}
		}
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("Failed to create job: %v", err))
	}

	job.JobID = id

	return c.JSON(http.StatusOK, echo.Map{"estimate": res, "job": job})
}

// Updates a job's values in the database
// Any field provided in the request body will remain unchanged
// Overwrites any fields provided with the new data
// Cannot change a finalized job
func updateJob(c echo.Context) error {
	var updatedJob models.Job
	// Pull all data from the request body into a job object
	if err := c.Bind(&updatedJob); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err})
	}

	// Find the old job being referenced
	oldJob, err := DB.PgInstance.GetJobByID(c.Request().Context(), updatedJob.JobID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "cannot find job by id", "error": err})
	}

	oldJob.JobID = updatedJob.JobID

	// If the job is already finalized, then it cannot be changed
	if oldJob.Finalized {
		return c.JSON(http.StatusConflict, echo.Map{"message": "cannot modify a finalized job"})
	}

	if updatedJob.Cost == 0 {
		updatedJob.Cost = oldJob.Cost
	}

	if updatedJob.ManHours == 0 {
		updatedJob.ManHours = oldJob.ManHours
	}

	if updatedJob.Rate == 0 {
		updatedJob.Rate = oldJob.Rate
	}

	if updatedJob.FinalCost == 0 {
		updatedJob.FinalCost = oldJob.FinalCost
	}

	if updatedJob.ActualManHours == 0 {
		updatedJob.ActualManHours = oldJob.ActualManHours
	}

	if len(updatedJob.Notes) == 0 {
		updatedJob.Notes = oldJob.Notes
	}

	// Update the job with the new data
	err = DB.PgInstance.UpdateJob(c.Request().Context(), updatedJob)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "cannot update job", "error": err})
	}

	// Query for the new job from the database to ensure the update happened correctly
	newJob, err := DB.PgInstance.GetJobByID(c.Request().Context(), updatedJob.JobID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "cannot find updated job", "error": err})
	}

	return c.JSON(http.StatusOK, echo.Map{"oldJob": oldJob, "updatedJob": newJob})
}

// Returns a list of all jobs owned by a customer with a given username
func getCustomerJobs(c echo.Context) error {
	// Get the username from the JWT
	username := c.Get("username").(string)

	jobs, err := DB.PgInstance.GetJobsByUsername(c.Request().Context(), username)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Error retrieving data: %v", err))
	}
	if jobs == nil {
		return c.String(http.StatusNotFound, fmt.Sprintf("No jobs found with username: %v", username))
	}
	return c.JSON(http.StatusOK, jobs)
}
