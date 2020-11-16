package main

import (
	"github.com/atrovanProject/handlers"
	_ "github.com/atrovanProject/jwt"
	"github.com/labstack/echo"
)

func main() {
	e := echo.New()
	e.GET("/", handlers.Root)

	r := e.Group("/register")
	r.POST("", handlers.Register)

	l := e.Group("/login")
	l.POST("", handlers.Login)

	add := e.Group("/courses/add")
	add.POST("", handlers.AddCourse)

	c := e.Group("/courses/list")
	c.GET("", handlers.GetCourses)

	res := e.Group("/courses/reserve")
	res.POST("", handlers.Reserve)

	mycourse := e.Group("/courses/student/reserved")
	mycourse.GET("", handlers.ListReserved)
	e.Logger.Fatal(e.Start(":8080"))
}
