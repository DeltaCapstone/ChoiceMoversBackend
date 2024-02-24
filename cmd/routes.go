package main

import "github.com/labstack/echo/v4"

func CreateRoutes(e *echo.Echo) {
	e.GET("/", hello)
	e.GET("/customer", getCustomers)
	e.POST("/customer", CreateCustomer)
	e.GET("/employee", getEmployees)
	e.POST("/employee", CreateEmployee)
}
