package main

import (
	"github.com/labstack/echo/v4"
)

func CreateRoutes(e *echo.Echo) {

	e.POST("/customer", createCustomer) //signup
	e.POST("/customer/login", customerLogin)
	e.GET("/customer/:id", getCustomer)       //view my account
	e.PUT("/customer/:id", updateCustomer)    //update my account
	e.DELETE("/customer/:id", deleteCustomer) //delete my account

	//employee portal
	//auth group
	e.POST("/employee/login", employeeLogin)
	e.GET("/employee/:id", getEmployeeInfo) //employee views thier own
	e.PUT("/employee/:id", updateEmployee)  //update my ccount
	e.GET("/jobs/:status", listJobs)        //veiw list of jobs by status (pending, confirmed, all)
	e.POST("/jobs/requestJobAssign/:job_id", requstAssign)

	//admin or manager auth group
	e.GET("/employee", listEmployees)   //manager view employees
	e.POST("/employee", createEmployee) //manager adds new employee?
	e.POST("/job", createJob)           //Manager should have the ability to manual create jobs, like if someone calls in
	e.PUT("/job/:id", updateJob)        //manger makes changes to a job or confirms a job

}
