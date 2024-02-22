package main

import (
	"net/http"
	//"golang.org/x/crypto/bcrypt"

	DB "github.com/DeltaCapstone/ChoiceMoversBackend/database"
	"github.com/labstack/echo/v4"
)

func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

func getUser(c echo.Context) error {
	value, err := DB.PgInstance.GetName(c.Request().Context(), "dakota")
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error()+" Error retrieving data")
	}
	return c.String(http.StatusOK, "dakota's id is: "+value)
}

/*
REMINDER!
use bcrypt for password hashing when
create user
login
*/
