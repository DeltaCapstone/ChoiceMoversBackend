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

func getUsers(c echo.Context) error {
	users, err := DB.PgInstance.GetUsers(c.Request().Context())
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
