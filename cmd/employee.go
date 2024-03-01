package main

import (
	"errors"
	"fmt"
	"net/http"

	DB "github.com/DeltaCapstone/ChoiceMoversBackend/database"
	"github.com/DeltaCapstone/ChoiceMoversBackend/utils"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
)

///////////////////////////////////////////////////////////////////////////////////////////////////////////
//Employee

//TODO: Redo error handling to get rid of of al lthe sprintf's

type CreateEmployeeRequest struct {
	UserName      string        `db:"username" json:"userName"`
	PasswordPlain string        `db:"password_plain" json:"passwordPlain"`
	FirstName     string        `db:"first_name" json:"firstName"`
	LastName      string        `db:"last_name" json:"lastName"`
	Email         string        `db:"email" json:"email"`
	PhonePrimary  pgtype.Text   `db:"phone_primary" json:"phonePrimary"`
	PhoneOther    []pgtype.Text `db:"phone_other" json:"phoneOther"`
	EmployeeType  string        `db:"employee_type" json:"employeeType"`
}

type EmployeeLoginRequest struct {
	UserName      string `db:"username" json:"userName"`
	PasswordPlain string `db:"password_plain" json:"passwordPlain"`
}

func listEmployees(c echo.Context) error {
	//id := c.QueryParam("id")
	users, err := DB.PgInstance.GetEmployeeList(c.Request().Context())
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Error retrieving data: %v", err))
	}
	if users == nil {
		return c.String(http.StatusNotFound, "No no employees found.")
	}
	return c.JSON(http.StatusOK, users)
}

func getEmployee(c echo.Context) error {
	username := c.Param("username")
	user, err := DB.PgInstance.GetEmployeeByUsername(c.Request().Context(), username)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Error retrieving data: %v", err))
	}
	if user.UserName == "" {
		return c.String(http.StatusNotFound, fmt.Sprintf("No user found with id: %v", username))
	}
	return c.JSON(http.StatusOK, user)
}

func createEmployee(c echo.Context) error {
	var newEmployee CreateEmployeeRequest
	// attempt at binding incoming json to a newUser
	if err := c.Bind(&newEmployee); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"Bind error": "Invalid user data"})
	}
	//validate password

	//replace plaintext password with hash
	hashedPassword, _ := utils.HashPassword(newEmployee.PasswordPlain)

	args := DB.CreateEmployeeParams{
		UserName:     newEmployee.UserName,
		PasswordHash: hashedPassword,
		FirstName:    newEmployee.FirstName,
		LastName:     newEmployee.LastName,
		Email:        newEmployee.Email,
		PhonePrimary: newEmployee.PhonePrimary,
		PhoneOther:   newEmployee.PhoneOther,
	}

	// validation stuff probably needed

	user, err := DB.PgInstance.CreateEmployee(c.Request().Context(), args)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case pgerrcode.UniqueViolation:
				fallthrough
			case pgerrcode.NotNullViolation:
				return c.JSON(http.StatusConflict, fmt.Sprintf("username or email already in use: %v", err))
			}
		}
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("Failed to create user: %v", err))
	}

	return c.JSON(http.StatusCreated, echo.Map{"username": user})
}

func updateEmployee(c echo.Context) error {
	var updatedEmployee DB.Employee

	// binding json to employee
	if err := c.Bind(&updatedEmployee); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input")
	}

	// update operation
	err := DB.PgInstance.UpdateEmployee(c.Request().Context(), updatedEmployee)
	if err != nil {
		// return internal server error if update fails
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update employee")
	}

	return c.JSON(http.StatusOK, "Employee updated")
}
