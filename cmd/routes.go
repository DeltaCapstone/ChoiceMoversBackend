package main

import (
	"github.com/labstack/echo/v4"
)

func CreateRoutes(e *echo.Echo) {

	e.POST("", createCustomer) //signup
	e.POST("/login", customerLogin)

	customerGroup := e.Group("/customer")
	//customerGroup.Use(customerAuthMiddleware)
	customerGroup.GET("/:username", getCustomer) //view my account
	customerGroup.PUT("/", updateCustomer)       //update my account
	//customerGroup.DELETE("/:username", deleteCustomer) //delete my account
	//customerGroup.GET("/job", getCustomerJobs)
	//customerGroup.POST("/job", createJobByCustomer)
	//customerGroup.PUT("/job/:job_id", updateJobByCustomer)

	//e.POST("/login", employeeLogin)  // Login

	// Group for employee routes
	employeeGroup := e.Group("/employee")
	//employeeGroup.Use(employeeAuthMiddleware)
	employeeGroup.GET("/:username", getEmployee) // Employee views their own
	employeeGroup.PUT("/", updateEmployee)       // Update my account
	employeeGroup.GET("/jobs", listJobs)         // View list of jobs by status (?status= pending, confirmed, all)
	//need to figure out how to limit query options for employees vs managers
	//employeeGroup.POST("/jobs/requestJobAssign/:job_id", requstAssign)

	// Group for manager routes
	managerGroup := e.Group("/manager")
	//managerGroup.Use(managerMiddleware)				// Add a middleware for manager authentication
	managerGroup.GET("/employee", listEmployees)               // Manager view employees
	managerGroup.POST("/employee", createEmployee)             // Manager adds new employee
	managerGroup.DELETE("/employee/:username", deleteEmployee) // Manager adds new employee
	//managerGroup.POST("/job", createJob)           // Manager creates a job, needed for cases where a customer call in or a job is recieved from Uhaul for example
	//managerGroup.PUT("/job/:job_id", updateJob)        // Manager makes changes to a job or confirms a job
}
