package handlers

import (
	"crypto/sha512"
	"fmt"
	"github.com/atrovanProject/db"
	"github.com/atrovanProject/jwt"
	"github.com/atrovanProject/user"
	"github.com/labstack/echo"
	"net/http"
)

type JsonResponses map[string]interface{}

func Register(c echo.Context) error {
	u := new(user.User)

	err := c.Bind(u)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	if err := u.GenerateID(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	u.Password = HashPassword(u.Password)

	err = db.Register(u)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	ts, err := jwt.CreateToken(u.ID, u.Role)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	}

	saveErr := jwt.CreateAuth(u.ID, u.Role, ts)
	if saveErr != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, saveErr.Error())
	}

	response := JsonResponses{
		"access_token":  ts.AccessToken,
		"refresh_token": ts.RefreshToken,
	}
	c.JSON(http.StatusOK, response)
	return nil
}

func HashPassword(password string) string {
	bytes := sha512.Sum512([]byte(password))
	hash := fmt.Sprintf("%x", bytes)
	return hash
}
