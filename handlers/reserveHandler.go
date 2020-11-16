package handlers

import (
	"github.com/atrovanProject/db"
	"github.com/atrovanProject/jwt"
	"github.com/labstack/echo"
	"net/http"
)

func Reserve(c echo.Context) error {
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

	u := new(db.Reserved)

	err = c.Bind(&u)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	u.Sid = sid

	err = db.AddReserve(u)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.NoContent(http.StatusOK)

}
