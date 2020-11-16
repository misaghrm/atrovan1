package handlers

import (
	"github.com/atrovanProject/db"
	"github.com/labstack/echo"
	"net/http"
)

func GetCourses(c echo.Context) error {
	courses := db.GetCourses()
	_ = JsonResponses{"Available Courses": courses}
	return c.JSON(http.StatusOK, courses)
}
