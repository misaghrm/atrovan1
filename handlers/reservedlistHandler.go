package handlers

import (
	"github.com/atrovanProject/db"
	"github.com/atrovanProject/jwt"
	"github.com/labstack/echo"
	"net/http"
)

func ListReserved(c echo.Context) error {
	err := jwt.TokenValid(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	tokenAuth, err := jwt.ExtractTokenMetadata(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	sid, role, err := jwt.FetchAuth(tokenAuth)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	if role != "student" {
		return echo.NewHTTPError(http.StatusUnauthorized, "just students can reserve course")
	}

	list, err := db.GetReservedList(sid)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, list)
}
