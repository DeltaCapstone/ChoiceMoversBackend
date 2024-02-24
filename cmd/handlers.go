package main

import (
	"fmt"
	"net/http"

	//"golang.org/x/crypto/bcrypt"

	DB "github.com/DeltaCapstone/ChoiceMoversBackend/database"
	"github.com/labstack/echo/v4"
)

func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

// accountType must match account types ENUM in db
func getCustomers(c echo.Context) error {
	id := c.QueryParam("id")
	users, err := DB.PgInstance.GetCustomers(c.Request().Context(), id)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Error retrieving data: %v", err))
	}
	//fmt.Print(users)

	return c.JSON(http.StatusOK, users)
}

// this is here temporarily cause I didnt want to mess with imports

// POST handler to create a new user
func CreateCustomer(c echo.Context) error {
	var newCustomer DB.Customer
	// attempt at binding incoming json to a newUser
	if err := c.Bind(&newCustomer); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid user data"})
	}

	// validation stuff probably needed

	userID, err := DB.PgInstance.CreateCustomer(c.Request().Context(), newCustomer)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to create user"})
	}

	return c.JSON(http.StatusCreated, echo.Map{"ID": userID})
}

func getEmployees(c echo.Context) error {
	id := c.QueryParam("id")
	users, err := DB.PgInstance.GetEmployees(c.Request().Context(), id)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Error retrieving data: %v", err))
	}
	//fmt.Print(users)

	return c.JSON(http.StatusOK, users)
}

func CreateEmployee(c echo.Context) error {
	var newEmployee DB.Employee
	// attempt at binding incoming json to a newUser
	if err := c.Bind(&newEmployee); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid user data"})
	}

	// validation stuff probably needed

	userID, err := DB.PgInstance.CreateEmployee(c.Request().Context(), newEmployee)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to create user"})
	}

	return c.JSON(http.StatusCreated, echo.Map{"ID": userID})
}

/*
REMINDER!
use bcrypt for password hashing when
create user
login
*/
