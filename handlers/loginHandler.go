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

func Login(c echo.Context) error {
	u := new(user.User)

	err := c.Bind(u)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if !db.IsRegistered(u.ID, u.Role) {
		return echo.NewHTTPError(http.StatusNotFound, "You Have Not Registered Yet")
	}

	hash, err := db.Password(u.ID, u.Role)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if !isCorrect(u.Password, hash) {
		return echo.NewHTTPError(http.StatusBadRequest, "entered password is incorrect")
	}

	ts, err := jwt.CreateToken(u.ID, u.Role)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err)
	}

	saveErr := jwt.CreateAuth(u.ID, u.Role, ts)
	if saveErr != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err)
	}

	response := JsonResponses{
		"access_token":  ts.AccessToken,
		"refresh_token": ts.RefreshToken,
	}
	return c.JSON(http.StatusOK, response)
}

func isCorrect(password, hash string) bool {
	bytes := sha512.Sum512([]byte(password))
	b := fmt.Sprintf("%x", bytes)
	if b != hash {
		return false
	}
	return true
}
