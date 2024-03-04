package main

import (
	"errors"
	"fmt"
	"net/http"

	DB "github.com/DeltaCapstone/ChoiceMoversBackend/database"
	models "github.com/DeltaCapstone/ChoiceMoversBackend/models"
	"github.com/DeltaCapstone/ChoiceMoversBackend/utils"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

func managerMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if c.Get("role") != "Manager" {
			return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
		}
		return next(c)
	}
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////
//Employee

//TODO: Redo error handling to get rid of of al lthe sprintf's

func listEmployees(c echo.Context) error {
	//id := c.QueryParam("id")
	users, err := DB.PgInstance.GetEmployeeList(c.Request().Context())
	if err != nil {
		zap.L().Sugar().Errorf("Failed to list employees: ", err.Error())
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Error retrieving data: %v", err))
	}
	if users == nil {
		return c.String(http.StatusNotFound, "No no employees found.")
	}
	return c.JSON(http.StatusOK, users)
}

func deleteEmployee(c echo.Context) error {
	username := c.Param("username")

	zap.L().Debug("deleteEmployee: ", zap.Any("Employee username", username))

	err := DB.PgInstance.DeleteEmployeeByUsername(c.Request().Context(), username)
	if err != nil {
		zap.L().Sugar().Errorf("Failed to delete employee: ", err.Error())
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Error deleting data: %v", err))
	}
	return c.NoContent(http.StatusNoContent)
}

func createEmployee(c echo.Context) error {
	var newEmployee models.CreateEmployeeRequest
	// attempt at binding incoming json to a newUser
	if err := c.Bind(&newEmployee); err != nil {
		zap.L().Sugar().Errorf("Failed to create employee: ", err.Error())
		return c.JSON(http.StatusBadRequest, echo.Map{"Bind error": "Invalid user data"})
	}

	zap.L().Debug("createEmployee", zap.Any("Employee", newEmployee))

	//validate password

	//replace plaintext password with hash
	hashedPassword, _ := utils.HashPassword(newEmployee.PasswordPlain)

	args := models.CreateEmployeeParams{
		UserName:     newEmployee.UserName,
		PasswordHash: hashedPassword,
		FirstName:    newEmployee.FirstName,
		LastName:     newEmployee.LastName,
		Email:        newEmployee.Email,
		PhonePrimary: newEmployee.PhonePrimary,
		PhoneOther:   newEmployee.PhoneOther,
		EmployeeType: newEmployee.EmployeeType,
	}

	// validation stuff probably needed

	err := DB.PgInstance.CreateEmployee(c.Request().Context(), args)
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
		zap.L().Sugar().Errorf("Failed to create employee: ", err.Error())
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("Failed to create employee: %v", err))
	}

	return c.JSON(http.StatusCreated, echo.Map{"username": newEmployee.UserName})
}

func getEmployee(c echo.Context) error {
	username := c.Param("username")

	zap.L().Debug("getEmployee: ", zap.Any("Employee username", username))

	user, err := DB.PgInstance.GetEmployeeByUsername(c.Request().Context(), username)
	if err != nil {
		zap.L().Sugar().Errorf("Failed to get employee: ", err.Error())
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Error retrieving data: %v", err))
	}
	if user.UserName == "" {
		return c.String(http.StatusNotFound, fmt.Sprintf("No user found with id: %v", username))
	}
	return c.JSON(http.StatusOK, user)
}

func updateEmployee(c echo.Context) error {
	var updatedEmployee models.Employee

	// binding json to employee
	if err := c.Bind(&updatedEmployee); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input")
	}

	zap.L().Debug("updateEmployee: ", zap.Any("Updated employee", updatedEmployee))

	//verify username on token matches username in struct

	// update operation
	err := DB.PgInstance.UpdateEmployee(c.Request().Context(), updatedEmployee)
	if err != nil {
		zap.L().Sugar().Errorf("Failed to update employee: ", err.Error())
		// return internal server error if update fails
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update employee")
	}

	return c.JSON(http.StatusOK, "Employee updated")
}

func employeeLogin(c echo.Context) error {
	var employeeLogin models.EmployeeLoginRequest

	// bind request data to the CustomerLoginRequest struct
	if err := c.Bind(&employeeLogin); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input")
	}

	// Get the customer with the username that was submitted
	hash, err := DB.PgInstance.GetEmployeeHashByUserName(c.Request().Context(), employeeLogin.UserName)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Error retrieving data: %v", err))
	}

	// Check that the user exists

	if hash == "" {
		return c.String(http.StatusNotFound, fmt.Sprintf("No user found with username: %v", employeeLogin.UserName))
	}
	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(employeeLogin.PasswordPlain))
	if err != nil {
		return c.String(http.StatusNotFound, fmt.Sprintf("Incorrect password for user with username: %v ", employeeLogin.UserName))
	}

	/*
		signedToken := MakeToken(employeeLogin.UserName, "Customer")


		return c.JSON(http.StatusOK, echo.Map{
			"token": signedToken,
		})
	*/

	return c.JSON(http.StatusOK, "Login Success")
}
