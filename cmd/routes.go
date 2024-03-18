package main

import (
	"github.com/DeltaCapstone/ChoiceMoversBackend/token"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func CreateRoutes(e *echo.Echo) {

	e.POST("", createCustomer) //signup
	e.POST("/login", customerLogin)
	//e.POST("/getEstimate",createEstimate)
	e.POST("/renewAccess", renewAccessToken)

	customerGroup := e.Group("/customer")
	customerGroup.Use(echojwt.WithConfig(token.Config), customerMiddleware)
	customerGroup.GET("/profile", getCustomer)    //view my account
	customerGroup.PUT("/profile", updateCustomer) //update my account
	//customerGroup.DELETE("/:username", deleteCustomer) //delete my account
	//customerGroup.GET("/job", getCustomerJobs)
	customerGroup.POST("/job", createJobByCustomer)
	//customerGroup.PUT("/job/:job_id", updateJobByCustomer)

	e.POST("/portal/login", employeeLogin) // Login

	// Group for employee routes
	employeeGroup := e.Group("/employee")
	employeeGroup.Use(echojwt.WithConfig(token.Config), employeeMiddleware)
	employeeGroup.GET("/profile", viewMyEmployeeProfile) // Employee views their own
	employeeGroup.PUT("/profile", updateEmployee)        // Update my account
	employeeGroup.GET("/jobs", listJobs)                 // View list of jobs by status (?status= pending, confirmed, all)
	//employeeGroup.POST("/jobs/requestJobAssign/:job_id", requstAssign)

	// Group for manager routes
	managerGroup := e.Group("/manager")
	managerGroup.Use(echojwt.WithConfig(token.Config), managerMiddleware) // Add a middleware for manager authentication
	managerGroup.GET("/employee", listEmployees)                          // Manager view employees
	managerGroup.POST("/employee", createEmployee)
	managerGroup.GET("/employee/:username", viewSomeEmployee)           // Manager views employee info
	managerGroup.DELETE("/employee/:username", deleteEmployee)          // Manager adds new employee
	managerGroup.PUT("/employee/:username", updateEmployeeTypePriority) //manager makes changes  to employee
	//managerGroup.POST("/job", createJob)           // Manager creates a job, needed for cases where a customer call in or a job is recieved from Uhaul for example
	//managerGroup.PUT("/job/:job_id", updateJob)        // Manager makes changes to a job or confirms a job
}
