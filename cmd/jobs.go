package main

import (
	//"errors"

	"errors"
	"fmt"
	"net/http"
	"time"

	DB "github.com/DeltaCapstone/ChoiceMoversBackend/database"
	"github.com/DeltaCapstone/ChoiceMoversBackend/models"
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
		Notes:      "",
	}

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

func updateJob(c echo.Context) error {
	var updatedJobRequest models.JobResponse
	if err := c.Bind(&updatedJobRequest); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err})
	}

	oldJob, err := DB.PgInstance.GetJobByID(c.Request().Context(), updatedJobRequest.JobID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "cannot find job by id", "error": err})
	}

	oldJob.JobID = updatedJobRequest.JobID

	var updatedJob models.Job

	updatedJob.JobID = updatedJobRequest.JobID

	if oldJob.Finalized {
		return c.JSON(http.StatusConflict, echo.Map{"message": "cannot modify a finalized job"})
	}

	if updatedJobRequest.AmountPaid < oldJob.AmountPaid {
		return c.JSON(http.StatusConflict, echo.Map{"message": "new paid amount lower than previous amount"})
	}

	if updatedJobRequest.Cost == 0 {
		updatedJob.Cost = oldJob.Cost
	} else {
		updatedJob.Cost = updatedJobRequest.Cost
	}

	if updatedJobRequest.ManHours == "" {
		updatedJob.ManHours = oldJob.ManHours
	} else {
		duration, _ := time.ParseDuration(updatedJobRequest.ManHours)
		var interval = pgtype.Interval{
			Microseconds: duration.Microseconds(),
			Valid:        true,
		}
		updatedJob.ManHours = interval
	}

	if updatedJobRequest.Rate == 0 {
		updatedJob.Rate = oldJob.Rate
	} else {
		updatedJob.Rate = updatedJobRequest.Rate
	}

	if updatedJobRequest.FinalCost == 0 {
		updatedJob.FinalCost = oldJob.FinalCost
	} else {
		updatedJob.FinalCost = updatedJobRequest.FinalCost
	}

	if updatedJobRequest.ActualManHours == "" {
		updatedJob.ActualManHours = oldJob.ActualManHours
	} else {
		duration, _ := time.ParseDuration(updatedJobRequest.ActualManHours)
		var interval = pgtype.Interval{
			Microseconds: duration.Microseconds(),
			Valid:        true,
		}
		updatedJob.ActualManHours = interval
	}

	if len(updatedJobRequest.Notes) == 0 {
		updatedJob.Notes = oldJob.Notes
	} else {
		updatedJob.Notes = updatedJobRequest.Notes
	}

	err = DB.PgInstance.UpdateJob(c.Request().Context(), updatedJob)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "cannot update job", "error": err})
	}

	newJob, err := DB.PgInstance.GetJobByID(c.Request().Context(), updatedJobRequest.JobID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "cannot find updated job", "error": err})
	}

	return c.JSON(http.StatusOK, echo.Map{"oldJob": oldJob, "updatedJob": newJob})
}

func getCustomerJobs(c echo.Context) error {
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
