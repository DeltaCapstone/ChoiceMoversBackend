package main

import (
	"github.com/DeltaCapstone/ChoiceMoversBackend/token"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func CreateRoutes(e *echo.Echo) {

	e.POST("", createCustomer) //signup
	e.POST("/login", customerLogin)
	//e.POST("/createEstimate",createEstimate)
	e.POST("/passwordReset", sendResetCodeCustomer)
	e.PUT("/passwordReset", resetPasswordCustomer)
	e.POST("/portal", createEmployee)      //signup
	e.POST("/portal/login", employeeLogin) // Login
	e.POST("/portal/passwordReset", sendResetCodeEmployee)
	e.PUT("/portal/passwordReset", resetPasswordEmployee)
	e.POST("/renewAccess", renewAccessToken)

	customerGroup := e.Group("/customer")
	customerGroup.Use(echojwt.WithConfig(token.Config), customerMiddleware)
	customerGroup.GET("/profile", getCustomer)    //view my account
	customerGroup.PUT("/profile", updateCustomer) //update my account
	//customerGroup.DELETE("/", deleteCustomer) //delete my account
	//customerGroup.GET("/job", getCustomerJobs)
	//customerGroup.POST("/job", createJobByCustomer)
	//customerGroup.PUT("/job/:jobID", updateJobByCustomer)
	customerGroup.PUT("/password", changeCustomerPassword)

	// Group for employee routes
	employeeGroup := e.Group("/employee")
	employeeGroup.Use(echojwt.WithConfig(token.Config), employeeMiddleware)
	employeeGroup.GET("/profile", viewMyEmployeeProfile)    // Employee views their own
	employeeGroup.PUT("/profile", updateEmployee)           // Update my account, data in json
	employeeGroup.GET("/jobs", listJobs)                    // View list of jobs by status
	employeeGroup.POST("/jobs/selfAssign", selfAssignToJob) // Query param "jobID"
	employeeGroup.PUT("/password", changeEmployeePassword)

	// Group for manager routes
	managerGroup := e.Group("/manager")
	managerGroup.Use(echojwt.WithConfig(token.Config), managerMiddleware)
	managerGroup.GET("/employee", listEmployees)                        // Manager view employees
	managerGroup.POST("/employee", addEmployee)                         // Query param "email"
	managerGroup.GET("/employee/:username", viewSomeEmployee)           // Manager views employee info
	managerGroup.DELETE("/employee/:username", deleteEmployee)          // Manager adds new employee
	managerGroup.PUT("/employee/:username", updateEmployeeTypePriority) //manager makes changes  to employee
	//managerGroup.POST("/job", createJob)           // Manager creates a job, needed for cases where a customer call in or a job is recieved from Uhaul for example
	//managerGroup.PUT("/job/:jobID", updateJob)        // Manager makes changes to a job or confirms a job

}
