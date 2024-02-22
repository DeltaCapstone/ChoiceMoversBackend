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

/*
REMINDER!
use bcrypt for password hashing when
create user
login
*/
