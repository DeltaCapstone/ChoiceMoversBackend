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
//Customer

type CreateCustomerRequest struct {
	UserName      string        `db:"username" json:"userName"`
	PasswordPlain string        `db:"password_plain" json:"passwordPlain"`
	FirstName     string        `db:"first_name" json:"firstName"`
	LastName      string        `db:"last_name" json:"lastName"`
	Email         string        `db:"email" json:"email"`
	PhonePrimary  pgtype.Text   `db:"phone_primary" json:"phonePrimary"`
	PhoneOther    []pgtype.Text `db:"phone_other" json:"phoneOther"`
}

// accountType must match account types ENUM in db
func getCustomer(c echo.Context) error {
	id := c.QueryParam("id")
	users, err := DB.PgInstance.GetCustomerById(c.Request().Context(), id)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Error retrieving data: %v", err))
	}
	if users == nil {
		return c.String(http.StatusNotFound, fmt.Sprintf("No user found with id: %v", id))
	}

	return c.JSON(http.StatusOK, users)
}

// POST handler to create a new user
func createCustomer(c echo.Context) error {
	var newCustomer CreateCustomerRequest
	// attempt at binding incoming json to a newUser
	if err := c.Bind(&newCustomer); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid user data"})
	}
	//validate password

	//replace plaintext password with hash
	hashedPassword, _ := utils.HashPassword(newCustomer.PasswordPlain)

	args := DB.CreateCustomerParams{
		UserName:     newCustomer.UserName,
		PasswordHash: hashedPassword,
		FirstName:    newCustomer.FirstName,
		LastName:     newCustomer.LastName,
		Email:        newCustomer.Email,
		PhonePrimary: newCustomer.PhonePrimary,
		PhoneOther:   newCustomer.PhoneOther,
	}

	// validation stuff probably needed

	user, err := DB.PgInstance.CreateCustomer(c.Request().Context(), args)
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

func updateCustomer(c echo.Context) error {
	var updatedCustomer DB.Customer
	// binding request
	if err := c.Bind(&updatedCustomer); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input")
	}

	// update operation
	err := DB.PgInstance.UpdateCustomer(c.Request().Context(), updatedCustomer)
	if err != nil {
		// return internal server error if update fails
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update customer")
	}

	return c.JSON(http.StatusOK, "Customer updated")

}
