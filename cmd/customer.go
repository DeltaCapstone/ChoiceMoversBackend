package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	DB "github.com/DeltaCapstone/ChoiceMoversBackend/database"
	"github.com/DeltaCapstone/ChoiceMoversBackend/utils"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

///////////////////////////////////////////////////////////////////////////////////////////////////////////
//Customer

//TODO: Redo error handling to get rid of of al lthe sprintf's

type CreateCustomerRequest struct {
	UserName      string        `db:"username" json:"userName"`
	PasswordPlain string        `db:"password_plain" json:"passwordPlain"`
	FirstName     string        `db:"first_name" json:"firstName"`
	LastName      string        `db:"last_name" json:"lastName"`
	Email         string        `db:"email" json:"email"`
	PhonePrimary  pgtype.Text   `db:"phone_primary" json:"phonePrimary"`
	PhoneOther    []pgtype.Text `db:"phone_other" json:"phoneOther"`
}

type CustomerLoginRequest struct {
	UserName      string `db:"username" json:"userName"`
	PasswordPlain string `db:"password_plain" json:"passwordPlain"`
}

func getCustomer(c echo.Context) error {
	id := c.Param("id")
	ID, err := strconv.Atoi(id)
	if err != nil {
		return fmt.Errorf("id is not an integer: %v", err)
	}
	user, err := DB.PgInstance.GetCustomerById(c.Request().Context(), ID)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Error retrieving data: %v", err))
	}
	if user.UserName == "" {
		return c.String(http.StatusNotFound, fmt.Sprintf("No user found with id: %v", id))
	}
	return c.JSON(http.StatusOK, user)
}

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

func customerLogin(c echo.Context) error {
	var customerLogin CustomerLoginRequest

	// bind request data to the CustomerLoginRequest struct
	if err := c.Bind(&customerLogin); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input")
	}

	// Get the customer with the username that was submitted
	user, err := DB.PgInstance.GetCustomerByUserName(c.Request().Context(), customerLogin.UserName)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Error retrieving data: %v", err))
	}

	// Check that the user exists
	if user.UserName == "" {
		return c.String(http.StatusNotFound, fmt.Sprintf("No user found with username: %v", customerLogin.UserName))
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(customerLogin.PasswordPlain))
	if err != nil {
		return c.String(http.StatusNotFound, fmt.Sprintf("Incorrect password for user with username: %v", customerLogin.UserName))
	}

	return c.JSON(http.StatusOK, "Login Success")
}
