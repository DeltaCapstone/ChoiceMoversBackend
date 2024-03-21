package main

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	DB "github.com/DeltaCapstone/ChoiceMoversBackend/database"
	"github.com/DeltaCapstone/ChoiceMoversBackend/mailer"
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
		c.Set("username", claims.Username)
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

func changeCustomerPassword(c echo.Context) error {
	var updatedCustomer models.UpdateCustomerPasswordRequest
	// binding request
	if err := c.Bind(&updatedCustomer); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input")
	}
	zap.L().Debug("updateCustomer: ", zap.Any("Customer password change request", updatedCustomer))

	if c.Get("username") != updatedCustomer.UserName {
		zap.L().Sugar().Errorf("Token username does not match updateCustomerPasswordRequest. ")
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input: username doesnt match")
	}

	storedHash, err := DB.PgInstance.GetCustomerCredentials(c.Request().Context(), updatedCustomer.UserName)
	if err != nil {
		zap.L().Sugar().Errorf("Error retrieving old password: ", err.Error())
		return c.String(http.StatusUnauthorized, "Something went wrong")
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(updatedCustomer.PasswordOld))
	if err != nil {
		zap.L().Sugar().Errorf("Wrong password supplied: ", err.Error())
		return c.String(http.StatusUnauthorized, fmt.Sprintf("Incorrect password for user with username: %v ", updatedCustomer.UserName))
	}

	hash, _ := utils.HashPassword(updatedCustomer.PasswordNew)

	if err := DB.PgInstance.UpdateCustomerPassword(c.Request().Context(), updatedCustomer.UserName, hash); err != nil {
		zap.L().Sugar().Errorf("Failed to update customer in db: ", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update password")
	}

	return c.JSON(http.StatusOK, "Password updated")
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

func sendResetCodeCustomer(c echo.Context) error {
	role := "customer"
	//get username from request
	username := c.QueryParam("username") //should do this differently
	//verify user in db
	customer, err := DB.PgInstance.GetCustomerByUserName(c.Request().Context(), username)
	if err != nil {
		zap.L().Sugar().Errorf("No customer with that username: ", err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input")
	}
	//create a password reset code, and write to db
	code, _ := utils.GenerateRandomCode(6)
	newReset := models.PasswordReset{
		Code:      code,
		Username:  customer.UserName,
		Email:     customer.Email,
		Role:      role,
		ExpiresAt: time.Now().Add(utils.ServerConfig.PasswordResetDuration),
	}
	_, err = DB.PgInstance.CreatePasswordReset(c.Request().Context(), newReset)
	if err != nil {
		zap.L().Sugar().Errorf("Error creating password reset: ", err.Error())
		return c.JSON(http.StatusInternalServerError, "Something went wrong.")
	}

	//for dev, just send the code in the response
	if utils.ServerConfig.Environment == "development" {
		return c.JSON(http.StatusCreated, echo.Map{"code": code})
	}

	//email it
	body := fmt.Sprintf("<p> Use this code to reset your password on chioce movers: <b>%s</b></p>", code)
	err = mailer.SendEmail("password reset code", body, []string{customer.Email}, nil, nil, nil)
	if err != nil {
		zap.L().Sugar().Errorf("Failed send email: ", err.Error())
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Failed to send email: %v", err))
	}
	return c.String(http.StatusCreated, "Code sent to your email")
}

func resetPasswordCustomer(c echo.Context) error {
	var pwrr models.PasswordResetRequest
	if err := c.Bind(&pwrr); err != nil {
		zap.L().Sugar().Errorf("Invalid password reset request: ", err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input")
	}
	pwr, err := DB.PgInstance.GetPasswordReset(c.Request().Context(), pwrr.Code)
	if err != nil {
		zap.L().Sugar().Errorf("Invalid password reset code, DNE: ", err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, "Something went wrong.")
	}
	if time.Now().After(pwr.ExpiresAt) {
		zap.L().Sugar().Errorf("Invalid password reset code, expired.")
		return echo.NewHTTPError(http.StatusBadRequest, "Something went wrong.")
	}
	if pwr.Role != "customer" {
		zap.L().Sugar().Errorf("Invalid password reset code, roles do not match. ")
		return echo.NewHTTPError(http.StatusBadRequest, "Something went wrong.")
	}

	if err := DB.PgInstance.DeletePasswordReset(c.Request().Context(), pwr.Code); err != nil {
		zap.L().Sugar().Errorf("Couldn't delete the password reset from DB: ", err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, "Something went wrong.")
	}

	newHash, _ := utils.HashPassword(pwrr.NewPW)
	if err := DB.PgInstance.UpdateCustomerPassword(c.Request().Context(), pwr.Username, newHash); err != nil {
		zap.L().Sugar().Errorf("Error updating password in database: ", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, "Something went wrong.")
	}

	return c.JSON(http.StatusOK, "Password updated")
}
