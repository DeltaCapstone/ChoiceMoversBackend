package main

import (
	//"errors"
	"fmt"
	"net/http"

	DB "github.com/DeltaCapstone/ChoiceMoversBackend/database"
	//"github.com/jackc/pgerrcode"
	//"github.com/jackc/pgx/v5/pgconn"
	"github.com/labstack/echo/v4"
)

//TODO: Redo error handling to get rid of of al lthe sprintf's

func listJobs(c echo.Context) error {
	status := c.QueryParam("status")
	if status == "" {
		status = "all"
	}
	jobs, err := DB.PgInstance.GetJobsByStatus(c.Request().Context(), status)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Error retrieving data: %v", err))
	}
	if jobs == nil {
		return c.String(http.StatusNotFound, fmt.Sprintf("No jobs found with status: %v", status))
	}
	return c.JSON(http.StatusOK, jobs)
}
