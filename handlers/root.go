package handlers

import (
	"github.com/labstack/echo"
	"net/http"
)

func Root(c echo.Context) error {
	return c.String(http.StatusOK, "Running API v1")
}
