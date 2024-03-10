package main

import (
	"errors"
	"fmt"
	"net/http"

	DB "github.com/DeltaCapstone/ChoiceMoversBackend/database"
	models "github.com/DeltaCapstone/ChoiceMoversBackend/models"
	"github.com/DeltaCapstone/ChoiceMoversBackend/token"
	"github.com/DeltaCapstone/ChoiceMoversBackend/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

func customerMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get("user").(*jwt.Token)
		claims := user.Claims.(*token.JwtCustomClaims)
		role := claims.Role
		c.Set("username", claims.UserName)
		c.Set("role", claims.Role)
		//return c.String(http.StatusFound, fmt.Sprintf("your role is %v", role))

		if role != "Customer" {
			return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
		}
		return next(c)

	}
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////
//Customer

//TODO: Redo error handling to get rid of of al lthe sprintf's

func getCustomer(c echo.Context) error {
	username := c.Get("username").(string)

	user, err := DB.PgInstance.GetCustomerByUserName(c.Request().Context(), username)
	if err != nil {
		zap.L().Sugar().Errorf("Error querying db for that username ", err.Error())
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Error retrieving data: %v", err))
	}
	if user.UserName == "" {
		zap.L().Sugar().Errorf("User with that username does not exist ")
		return c.String(http.StatusNotFound, fmt.Sprintf("No user found with id: %v", username))
	}
	return c.JSON(http.StatusOK, user)
}

func createCustomer(c echo.Context) error {
	var newCustomer models.CreateCustomerRequest
	// attempt at binding incoming json to a newUser
	if err := c.Bind(&newCustomer); err != nil {
		zap.L().Sugar().Errorf("Incorrect data format to create customer: ", err.Error())
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
		zap.L().Sugar().Errorf("Error adding customer to db: ", err.Error())
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("Failed to create user: %v", err))
	}

	return c.JSON(http.StatusCreated, echo.Map{"username": user})
}

func updateCustomer(c echo.Context) error {
	var updatedCustomer models.UpdateCustomerParams
	// binding request
	if err := c.Bind(&updatedCustomer); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input")
	}

	zap.L().Debug("updateCustomer: ", zap.Any("Updated customer", updatedCustomer))

	if c.Get("username") != updatedCustomer.UserName {
		zap.L().Sugar().Errorf("Token username does not match updateCustomerParams. ")
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input: username doesnt match")
	}

	// update operation
	err := DB.PgInstance.UpdateCustomer(c.Request().Context(), updatedCustomer)
	if err != nil {
		zap.L().Sugar().Errorf("Failed to update customer in db: ", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update customer")
	}

	return c.JSON(http.StatusOK, "Customer updated")

}

func customerLogin(c echo.Context) error {
	var customerLogin models.CustomerLoginRequest

	// bind request data to the CustomerLoginRequest struct
	if err := c.Bind(&customerLogin); err != nil {
		zap.L().Sugar().Errorf("Invalid loggin request format: ", err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input")
	}

	// Get the customer with the username that was submitted
	hash, err := DB.PgInstance.GetCustomerCredentials(c.Request().Context(), customerLogin.UserName)
	if err != nil {
		zap.L().Sugar().Errorf("Could not retrieve credentials for comparison: ", err.Error())
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Error retrieving data: %v", err))
	}

	if hash == "" {
		zap.L().Sugar().Errorf("could not find customer with that username. ")
		return c.String(http.StatusNotFound, fmt.Sprintf("No user found with username: %v", customerLogin.UserName))
		//return echo.ErrUnauthorized
	}
	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(customerLogin.PasswordPlain))
	if err != nil {
		zap.L().Sugar().Errorf("Wrong password supplied: ", err.Error())
		return c.String(http.StatusNotFound, fmt.Sprintf("Incorrect password for user with username: %v ", customerLogin.UserName))
		//return echo.ErrUnauthorized
	}

	accessToken, accessClaims, err := token.MakeAccessToken(customerLogin.UserName, "Customer")
	if err != nil {
		zap.L().Sugar().Errorf("problem making access token: ", err.Error())
		return c.String(http.StatusInternalServerError, "Error creating access token")
	}

	refreshToken, refreshClaims, err := token.MakeAccessToken(customerLogin.UserName, "Customer")
	if err != nil {
		zap.L().Sugar().Errorf("problem making refresh token: ", err.Error())
		return c.String(http.StatusInternalServerError, "Error creating refresh token")
	}

	sessionId, err := DB.PgInstance.CreateSession(c.Request().Context(), models.CreateSessionParams{
		ID:           refreshClaims.TokenID,
		Username:     customerLogin.UserName,
		RefreshToken: refreshToken,
		UserAgent:    c.Request().UserAgent(),
		ClientIp:     c.RealIP(),
		IsBlocked:    false,
		ExpiresAt:    refreshClaims.ExpiresAt.Time,
	})
	if err != nil {
		zap.L().Sugar().Errorf("problem creating session: ", err.Error())
		return c.String(http.StatusInternalServerError, "Error creating session")
	}

	rsp := models.CustomerLoginResponse{
		SessionId:             sessionId,
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessClaims.ExpiresAt.Time,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshClaims.ExpiresAt.Time,
		Username:              customerLogin.UserName,
	}

	return c.JSON(http.StatusOK, rsp)

}
