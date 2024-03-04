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
	"golang.org/x/crypto/bcrypt"
)

///////////////////////////////////////////////////////////////////////////////////////////////////////////
//Customer

//TODO: Redo error handling to get rid of of al lthe sprintf's

func getCustomer(c echo.Context) error {
	username := c.Param("username")

	user, err := DB.PgInstance.GetCustomerByUserName(c.Request().Context(), username)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Error retrieving data: %v", err))
	}
	if user.UserName == "" {
		return c.String(http.StatusNotFound, fmt.Sprintf("No user found with id: %v", username))
	}
	return c.JSON(http.StatusOK, user)
}

func createCustomer(c echo.Context) error {
	var newCustomer models.CreateCustomerRequest
	// attempt at binding incoming json to a newUser
	if err := c.Bind(&newCustomer); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid user data"})
	}
	//validate password

	//replace plaintext password with hash
	hashedPassword, _ := utils.HashPassword(newCustomer.PasswordPlain)

	args := models.CreateCustomerParams{
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
	var updatedCustomer models.Customer
	// binding request
	if err := c.Bind(&updatedCustomer); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input")
	}

	//verify username on token matches username in struct

	// update operation
	err := DB.PgInstance.UpdateCustomer(c.Request().Context(), updatedCustomer)
	if err != nil {
		// return internal server error if update fails
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update customer")
	}

	return c.JSON(http.StatusOK, "Customer updated")

}

func customerLogin(c echo.Context) error {
	var customerLogin models.CustomerLoginRequest

	// bind request data to the CustomerLoginRequest struct
	if err := c.Bind(&customerLogin); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input")
	}

	// Get the customer with the username that was submitted
	hash, err := DB.PgInstance.GetCustomerHashByUserName(c.Request().Context(), customerLogin.UserName)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Error retrieving data: %v", err))
	}

	// Check that the user exists

	if hash == "" {
		return c.String(http.StatusNotFound, fmt.Sprintf("No user found with username: %v", customerLogin.UserName))
	}
	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(customerLogin.PasswordPlain))
	if err != nil {
		return c.String(http.StatusNotFound, fmt.Sprintf("Incorrect password for user with username: %v ", customerLogin.UserName))
	}

	/*
		signedToken := MakeToken(customerLogin.UserName, "Customer")


		return c.JSON(http.StatusOK, echo.Map{
			"token": signedToken,
		})
	*/

	return c.JSON(http.StatusOK, "Login Success")
}
