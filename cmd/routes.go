package main

import (
	"github.com/labstack/echo/v4"
)

func CreateRoutes(e *echo.Echo) {

	customerGroup := e.Group("/customer")
	customerGroup.POST("", createCustomer) //signup
	//customerGroup.POST("/login", customerLogin)
	customerGroup.GET("/:id", getCustomer)    //view my account
	customerGroup.PUT("/:id", updateCustomer) //update my account
	//customerGroup.DELETE("/:id", deleteCustomer) //delete my account
	//customerGroup.GET("/job", getCustomerJobs)
	//customerGroup.POST("/job", createJobByCustomer)
	//customerGroup.PUT("/job/:id", updateJobByCustomer)

	// Group for employee routes
	employeeGroup := e.Group("/employee")
	//employeeGroup.POST("/login", employeeLogin)  // Login
	employeeGroup.GET("/:id", getEmployee)    // Employee views their own
	employeeGroup.PUT("/:id", updateEmployee) // Update my account
	employeeGroup.GET("/jobs", listJobs)      // View list of jobs by status (pending, confirmed, all)
	//employeeGroup.POST("/jobs/requestJobAssign/:job_id", requstAssign)

	// Group for manager routes
	managerGroup := e.Group("/manager")
	//managerGroup.Use(managerMiddleware)				// Add a middleware for manager authentication
	managerGroup.GET("/employee", listEmployees)   // Manager view employees
	managerGroup.POST("/employee", createEmployee) // Manager adds new employee
	//managerGroup.POST("/job", createJob)           // Manager creates a job, needed for cases where a customer call in or a job is recieved from Uhaul for example
	//managerGroup.PUT("/job/:id", updateJob)        // Manager makes changes to a job or confirms a job
}
