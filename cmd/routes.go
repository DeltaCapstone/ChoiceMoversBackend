package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func CreateRoutes(e *echo.Echo) {
	//e.GET("/", hello)
	e.POST("/customer", createCustomer) //signup
	e.POST("/customer/login", customerLogin)
	e.GET("/customer/:id", getCustomer)    //view my account
	e.PUT("/customer/:id", updateCustomer) //updateAccount
	e.DELETE("/customer/:id", deleteCustomer)

	//employee portal
	//auth group
	e.POST("/employee/login", employeeLogin)
	e.GET("/employee/:id", getEmployeeInfo) //employee views thier own
	e.PUT("/employee/:id", updateEmployee)  //updateAccount
	e.GET("/jobs/:status", listJobs)        //veiw list of jobs by status (pending, confirmed, all)
	e.POST("/jobs/requestJobAssign/:job_id", requstAssign)

	//admin or manager auth group
	e.GET("/employee", listEmployees)
	e.POST("/employee", createEmployee) //manager adds new employee?
	e.POST("/job", createJob)
	e.PUT("/job/:id", updateJob)

}

func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
