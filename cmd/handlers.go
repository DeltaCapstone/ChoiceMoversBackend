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
func getUsers(c echo.Context) error {
	accountType := c.QueryParam("accountType")
	users, err := DB.PgInstance.GetUsers(c.Request().Context(), accountType)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Error retrieving data: %v", err))
	}
	//fmt.Print(users)

	return c.JSON(http.StatusOK, users)
}

// this is here temporarily cause I didnt want to mess with imports
type User struct {
	ID          int
	UserName    string
	AccountType string
	Email       string
}

// POST handler to create a new user
func CreateUser(c echo.Context) error {
	var newUser User
	// attempt at binding incoming json to a newUser
	if err := c.Bind(&newUser); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid user data"})
	}

	// validation stuff probably needed

	userID, err := DB.PgInstance.CreateUser(c.Request().Context(), newUser)
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
