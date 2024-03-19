package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

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

func managerMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get("user").(*jwt.Token)
		claims := user.Claims.(*token.JwtCustomClaims)
		role := claims.Role
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)
		//return c.String(http.StatusFound, fmt.Sprintf("your role is %v", role))
		if role != "Manager" {
			return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
		}
		return next(c)
	}
}

func employeeMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get("user").(*jwt.Token)
		claims := user.Claims.(*token.JwtCustomClaims)
		role := claims.Role
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)
		//return c.String(http.StatusFound, fmt.Sprintf("your role is %v", role))
		if (role != "Full-time") && (role != "Part-time") && (role != "Manager") {
			return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
		}
		return next(c)
	}
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////
//Employee

//TODO: Redo error handling to get rid of of al lthe sprintf's

// //////////////////////////////////////
// Manager routes
// /////////////////////////////////////
func listEmployees(c echo.Context) error {
	//id := c.QueryParam("id")
	users, err := DB.PgInstance.GetEmployeeList(c.Request().Context())
	if err != nil {
		zap.L().Sugar().Errorf("Error querying db for employess: ", err.Error())
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Error retrieving data: %v", err))
	}
	if users == nil {
		zap.L().Sugar().Errorf("No Employees Found: ")
		return c.String(http.StatusNotFound, "No no employees found.")
	}
	return c.JSON(http.StatusOK, users)
}

func deleteEmployee(c echo.Context) error {
	username := c.Param("username")

	zap.L().Debug("deleteEmployee: ", zap.Any("Employee username", username))

	err := DB.PgInstance.DeleteEmployeeByUsername(c.Request().Context(), username)
	if err != nil {
		zap.L().Sugar().Errorf("Failed to delete employee: ", err.Error())
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Error deleting data: %v", err))
	}
	return c.NoContent(http.StatusOK)
}

// config
const empSignupURL = "www.choicemovers.com/portal?token="
const signupMessage = "<h3> use the link to create your Employee Account</h3>"

func addEmployee(c echo.Context) error {
	e := c.QueryParam("email")
	et, ok := models.IsValidEmployeeType(c.QueryParam("type"))
	if !ok {
		return c.String(http.StatusBadRequest, "not a valid employee type")
	}
	p, err := strconv.Atoi(c.QueryParam("priority"))
	if err != nil {
		return c.String(http.StatusBadRequest, "could not parse priority")
	}

	//make a token for the link and to store
	t, claims, err := token.MakeEmployeeSignupToken(e)
	if err != nil {
		return err
	}
	//store
	newEmployee := models.EmployeeSignup{
		Id:               claims.TokenID,
		Email:            e,
		EmployeeType:     et,
		EmployeePriority: p,
		SignupToken:      t,
		ExpiresAt:        claims.ExpiresAt.Time,
		Used:             false,
	}
	DB.PgInstance.AddEmployeeSignup(c.Request().Context(), newEmployee)
	//make url

	url := fmt.Sprintf("%v", empSignupURL+t)
	link := fmt.Sprintf(`<p><a href="%s">Signup Link</a></p>`, url)
	//email it
	body := signupMessage + link
	err = mailer.SendEmail("new employee link", body, []string{e}, nil, nil, nil)
	if err != nil {
		zap.L().Sugar().Errorf("Failed send email: ", err.Error())
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Failed to send email: %v", err))
	}

	return c.String(http.StatusOK, fmt.Sprintf("Email sent to: %v\n", e))
}

func createEmployee(c echo.Context) error {
	//check token
	token := c.QueryParam("token")
	claims, err := VerifyEmployeeSignupToken(token)
	if err != nil {
		zap.L().Sugar().Errorf("Could not parse signup token in url or invalid token: ", err.Error())
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("could not parse url token or invalid token: %v", err.Error()))
	}
	//check db
	su, err := DB.PgInstance.GetEmployeeSignup(c.Request().Context(), claims.TokenID)
	if err != nil {
		zap.L().Sugar().Errorf("Signup token does not exist: ", err.Error())
		return c.JSON(http.StatusUnauthorized, "signup token does not exist")
	}
	//make sure everything matches and token/link hasn't alreayd been used
	if su.Email != claims.Email ||
		su.Id != claims.TokenID ||
		su.SignupToken != token ||
		su.Used {
		zap.L().Sugar().Errorf("tokens do not match or this token is already used.")
		return c.JSON(http.StatusBadRequest, "token does not match stored parameters or has already been used.")
	}
	// set the signup 'used' field in db to true
	err = DB.PgInstance.UseEmployeeSignup(c.Request().Context(), claims.TokenID)
	if err != nil {
		zap.L().Sugar().Errorf("could not update employee signup entry: ", err.Error())
		return c.JSON(http.StatusUnauthorized, "could not update employee signup entry.")
	}

	var newEmployee models.CreateEmployeeRequest
	// attempt at binding incoming json to a newUser
	if err := c.Bind(&newEmployee); err != nil {
		zap.L().Sugar().Errorf("Incorrect data format for creating employee: ", err.Error())
		return c.JSON(http.StatusBadRequest, echo.Map{"Bind error": "Invalid user data"})
	}

	if newEmployee.Email != claims.Email {
		return c.JSON(http.StatusUnauthorized, "email entered does not match link recipiant")
	}

	zap.L().Debug("createEmployee", zap.Any("Employee", newEmployee))

	//validate password

	//replace plaintext password with hash
	hashedPassword, _ := utils.HashPassword(newEmployee.PasswordPlain)

	args := models.CreateEmployeeParams{
		UserName:         newEmployee.UserName,
		PasswordHash:     hashedPassword,
		FirstName:        newEmployee.FirstName,
		LastName:         newEmployee.LastName,
		Email:            newEmployee.Email,
		PhonePrimary:     newEmployee.PhonePrimary,
		PhoneOther:       newEmployee.PhoneOther,
		EmployeeType:     newEmployee.EmployeeType,
		EmployeePriority: newEmployee.EmployeePriority,
	}

	// validation stuff probably needed

	err = DB.PgInstance.CreateEmployee(c.Request().Context(), args)
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
		zap.L().Sugar().Errorf("Error adding employee to db: ", err.Error())
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("Failed to create employee: %v", err))
	}

	return c.JSON(http.StatusCreated, echo.Map{"username": newEmployee.UserName})
}

// /////////////////////////////////////////
// Self routes
// /////////////////////////////////////////

// wrapper for getEmployee when used with employee/profile
func viewMyEmployeeProfile(c echo.Context) error {
	return getEmployee(c, c.Get("username").(string))
}

// wrapper for getEmployee when used with manager/employee/:username
func viewSomeEmployee(c echo.Context) error {
	return getEmployee(c, c.Param("username"))
}

func getEmployee(c echo.Context, username string) error {

	zap.L().Debug("getEmployee: ", zap.Any("Employee username", username))

	user, err := DB.PgInstance.GetEmployeeByUsername(c.Request().Context(), username)
	if err != nil {
		zap.L().Sugar().Errorf("Failed to get employee: ", err.Error())
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Error retrieving data: %v", err))
	}
	if user.UserName == "" {
		return c.String(http.StatusNotFound, fmt.Sprintf("No user found with username: %v", username))
	}
	return c.JSON(http.StatusOK, user)
}

func updateEmployee(c echo.Context) error {
	var updatedEmployee models.UpdateEmployeeParams

	// binding json to employee
	if err := c.Bind(&updatedEmployee); err != nil {
		zap.L().Sugar().Errorf("Failed to update employee: ", err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input")
	}

	zap.L().Debug("updateEmployee: ", zap.Any("Updated employee", updatedEmployee))

	if c.Get("username") != updatedEmployee.UserName {
		zap.L().Sugar().Errorf("Token username does not match updateEmployeeParams. ")
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input: username doesnt match")
	}

	// update operation
	err := DB.PgInstance.UpdateEmployee(c.Request().Context(), updatedEmployee)
	if err != nil {
		zap.L().Sugar().Errorf("Failed to update employee in db: ", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update employee")
	}

	return c.JSON(http.StatusOK, "Employee updated")
}

func updateEmployeeTypePriority(c echo.Context) error {
	var updatedEmployee models.UpdateEmployeeTypePriorityParams
	if err := c.Bind(&updatedEmployee); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input")
	}
	zap.L().Debug("updateEmployee: ", zap.Any("Updated employee", updatedEmployee))

	//db querry
	return c.JSON(http.StatusOK, "Employee updated")
}

func employeeLogin(c echo.Context) error {
	var employeeLogin models.EmployeeLoginRequest

	// bind request data to the CustomerLoginRequest struct
	if err := c.Bind(&employeeLogin); err != nil {
		zap.L().Sugar().Errorf("Invalid loggin request format: ", err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input")
	}

	// Get the customer with the username that was submitted
	hash, err := DB.PgInstance.GetEmployeeCredentials(c.Request().Context(), employeeLogin.UserName)
	if err != nil {
		zap.L().Sugar().Errorf("Could not retrieve credentials for comparison: ", err.Error())
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Error retrieving data: %v", err))
	}

	if hash == "" {
		zap.L().Sugar().Errorf("could not find employee with that username. ")
		return c.String(http.StatusNotFound, fmt.Sprintf("No user found with username: %v", employeeLogin.UserName))
		//return echo.ErrUnauthorized
	}
	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(employeeLogin.PasswordPlain))
	if err != nil {
		zap.L().Sugar().Errorf("Wrong password supplied: ", err.Error())
		return c.String(http.StatusNotFound, fmt.Sprintf("Incorrect password for user with username: %v ", employeeLogin.UserName))
	}

	role, err := DB.PgInstance.GetEmployeeRole(c.Request().Context(), employeeLogin.UserName)
	if err != nil {
		zap.L().Sugar().Errorf("Could not retrieve role: ", err.Error())
		return c.String(http.StatusInternalServerError, "Could not deterimine role")
	}

	accessToken, accessClaims, err := token.MakeAccessToken(employeeLogin.UserName, role)
	if err != nil {
		zap.L().Sugar().Errorf("problem making access token: ", err.Error())
		return c.String(http.StatusInternalServerError, "Error creating access token")
	}

	refreshToken, refreshClaims, err := token.MakeAccessToken(employeeLogin.UserName, role)
	if err != nil {
		zap.L().Sugar().Errorf("problem making refresh token: ", err.Error())
		return c.String(http.StatusInternalServerError, "Error creating refresh token")
	}

	sessionId, err := DB.PgInstance.CreateSession(c.Request().Context(), models.CreateSessionParams{
		ID:           refreshClaims.TokenID,
		Username:     employeeLogin.UserName,
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

	rsp := models.EmployeeLoginResponse{
		SessionId:             sessionId,
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessClaims.ExpiresAt.Time,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshClaims.ExpiresAt.Time,
		Username:              employeeLogin.UserName,
	}

	return c.JSON(http.StatusOK, rsp)
}
