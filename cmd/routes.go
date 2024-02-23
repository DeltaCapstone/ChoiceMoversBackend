package main

import "github.com/labstack/echo/v4"

func CreateRoutes(e *echo.Echo) {
	e.GET("/", hello)
	e.GET("/users", getUsers)
	e.POST("/users", CreateUser)
}
