package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
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
// Employee Management routes
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

const signupMessage = "<h4> Use the link to create your Employee Account</h4> <p>This link is good for 15 minutes</p> "

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

	url := fmt.Sprintf("%v", utils.ServerConfig.EmpSignupURL+"?token="+t)
	link := fmt.Sprintf(`<p><a href="%s">Signup Link</a></p>`, url)
	//email it
	body := signupMessage + link

	//for dev, just send the url in the response
	if utils.ServerConfig.Environment == "development" {
		return c.JSON(http.StatusCreated, echo.Map{"url": url})
	}

	err = mailer.SendEmail("new employee link", body, []string{e}, nil, nil, nil)
	if err != nil {
		zap.L().Sugar().Errorf("Failed send email: ", err.Error())
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Failed to send email: %v", err))
	}
	return c.JSON(http.StatusCreated, echo.Map{"Email sent to": e})

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

func updateEmployeeTypePriority(c echo.Context) error {
	var updatedEmployee models.UpdateEmployeeTypePriorityParams
	if err := c.Bind(&updatedEmployee); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input")
	}
	zap.L().Debug("updateEmployee: ", zap.Any("Updated employee", updatedEmployee))

	if err := DB.PgInstance.UpdateEmployeeTypePriority(c.Request().Context(), updatedEmployee); err != nil {
		zap.L().Sugar().Errorf("Failed to update employee in db: ", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update employee")
	}
	return c.JSON(http.StatusOK, "Employee updated")
}

// wrapper for getEmployee when used with manager/employee/:username
func viewSomeEmployee(c echo.Context) error {
	return getEmployee(c, c.Param("username"))
}

// /////////////////////////////////////////
// Self routes
// /////////////////////////////////////////

// wrapper for getEmployee when used with employee/profile
func viewMyEmployeeProfile(c echo.Context) error {
	return getEmployee(c, c.Get("username").(string))
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

func changeEmployeePassword(c echo.Context) error {
	var updatedEmployee models.UpdateEmployeePasswordRequest
	// binding request
	if err := c.Bind(&updatedEmployee); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input")
	}
	zap.L().Debug("updateEmployee: ", zap.Any("Employee password change request", updatedEmployee))

	if c.Get("username") != updatedEmployee.UserName {
		zap.L().Sugar().Errorf("Token username does not match updateEmployeePasswordRequest. ")
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input: username doesnt match")
	}

	storedHash, err := DB.PgInstance.GetEmployeeCredentials(c.Request().Context(), updatedEmployee.UserName)
	if err != nil {
		zap.L().Sugar().Errorf("Error retrieving old password: ", err.Error())
		return c.String(http.StatusUnauthorized, "Something went wrong")
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(updatedEmployee.PasswordOld))
	if err != nil {
		zap.L().Sugar().Errorf("Wrong password supplied: ", err.Error())
		return c.String(http.StatusUnauthorized, fmt.Sprintf("Incorrect password for user with username: %v ", updatedEmployee.UserName))
	}

	hash, _ := utils.HashPassword(updatedEmployee.PasswordNew)

	if err := DB.PgInstance.UpdateEmployeePassword(c.Request().Context(), updatedEmployee.UserName, hash); err != nil {
		zap.L().Sugar().Errorf("Failed to update employee in db: ", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update password")
	}

	return c.JSON(http.StatusOK, "Password updated")
}

func employeeLogin(c echo.Context) error {
	var employeeLogin models.EmployeeLoginRequest

	// bind request data to the employeeLoginRequest struct
	if err := c.Bind(&employeeLogin); err != nil {
		zap.L().Sugar().Errorf("Invalid loggin request format: ", err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input")
	}

	// Get the employee with the username that was submitted
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

func sendResetCodeEmployee(c echo.Context) error {
	//get username from request
	username := c.QueryParam("username") //should do this differently
	//verify user in db
	employee, err := DB.PgInstance.GetEmployeeByUsername(c.Request().Context(), username)
	if err != nil {
		zap.L().Sugar().Errorf("No employee with that username: ", err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input")
	}
	//create a password reset code, and write to db
	code, _ := utils.GenerateRandomCode(6)
	newReset := models.PasswordReset{
		Code:      code,
		Username:  employee.UserName,
		Email:     employee.Email,
		Role:      employee.EmployeeType,
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
	body := fmt.Sprintf("<p> Use this code to reset your password on chioce movers employee portal: <b>%s</b></p>", code)
	err = mailer.SendEmail("password reset code", body, []string{employee.Email}, nil, nil, nil)
	if err != nil {
		zap.L().Sugar().Errorf("Failed send email: ", err.Error())
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Failed to send email: %v", err))
	}
	return c.String(http.StatusCreated, "Code sent to your email")
}

func resetPasswordEmployee(c echo.Context) error {
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
	if pwr.Role != "Full-time" && pwr.Role != "Part-time" && pwr.Role != "Manager" {
		zap.L().Sugar().Errorf("Invalid password reset code, roles do not match. ")
		return echo.NewHTTPError(http.StatusBadRequest, "Something went wrong.")
	}

	if err := DB.PgInstance.DeletePasswordReset(c.Request().Context(), pwr.Code); err != nil {
		zap.L().Sugar().Errorf("Couldn't delete the password reset from DB: ", err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, "Something went wrong.")
	}

	newHash, _ := utils.HashPassword(pwrr.NewPW)
	if err := DB.PgInstance.UpdateEmployeePassword(c.Request().Context(), pwr.Username, newHash); err != nil {
		zap.L().Sugar().Errorf("Error updating password in database: ", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, "Something went wrong.")
	}

	return c.JSON(http.StatusOK, "Password updated")
}

func selfAssignToJob(c echo.Context) error {
	me := c.Get("username").(string)
	myPriority, err := DB.PgInstance.GetEmployeePriority(c.Request().Context(), me)
	if err != nil {
		zap.L().Sugar().Errorf("Error retriving user's employee priority from DB: ", err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, "Something went wrong.")
	}

	jobId, _ := strconv.Atoi(c.QueryParam("jobID"))
	n, err := DB.PgInstance.GetNumWorksForJob(c.Request().Context(), jobId)
	if err != nil {
		zap.L().Sugar().Errorf("Error retriving Number of Worker for Job: ", err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, "Something went wrong.")
	}

	assignedEmps, err := DB.GetAssignedEmployees(c.Request().Context(), jobId)
	if err != nil {
		zap.L().Sugar().Errorf("Error retriving list of assigned employees: ", err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, "Something went wrong.")
	}

	if n > len(assignedEmps) {
		DB.PgInstance.AddEmployeeToJob(c.Request().Context(), me, jobId, false)
	} else {
		min_priority_i := -1
		min_priority := myPriority
		for i, e := range assignedEmps {
			//priority is golf rules, smaller number = higher priority, ie min priority is actually max priority value
			//aslo start with own priority as the min
			if !e.ManagerAssigned && min_priority < e.EmployeePriority {
				min_priority = e.EmployeePriority
				min_priority_i = i
			}
		}
		//replace
		if min_priority_i != -1 {
			toBoot := assignedEmps[min_priority_i].UserName
			if err := DB.PgInstance.RemoveEmployeeFromJob(c.Request().Context(), toBoot, jobId); err != nil {
				zap.L().Sugar().Errorf("Error removing employee from this job in DB: ", err.Error())
				return echo.NewHTTPError(http.StatusBadRequest, "Something went wrong.")
			} else if err := DB.PgInstance.AddEmployeeToJob(c.Request().Context(), me, jobId, false); err != nil {
				zap.L().Sugar().Errorf("Error add user to job in DB: ", err.Error())
				return echo.NewHTTPError(http.StatusBadRequest, "Something went wrong.")
			}

		} else {
			return c.String(http.StatusAccepted, "Job is Full, and there is no one you are allowed to boot.")
		}
	}
	assignedEmps, err = DB.GetAssignedEmployees(c.Request().Context(), jobId)
	if err != nil {
		zap.L().Sugar().Errorf("Error retriving list of assigned employees: ", err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, "Something went wrong.")
	}
	return c.JSON(http.StatusCreated, assignedEmps)
}
