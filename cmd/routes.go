package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func CreateRoutes(e *echo.Echo) {
	e.GET("/", hello)
	e.GET("/customer", getCustomers)
	//e.GET("/customer/:id",getCustomer)
	e.POST("/customer", CreateCustomer)
	e.PUT("/customer", UpdateCustomer)

	e.GET("/employee", getEmployees)
	//e.GET("/employee/:id",getCustomer)
	e.POST("/employee", CreateEmployee)
	e.PUT("/employee", UpdateEmployee)

}

func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
