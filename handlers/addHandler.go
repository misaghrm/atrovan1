package handlers

import (
	"github.com/atrovanProject/db"
	"github.com/atrovanProject/jwt"
	"github.com/atrovanProject/user"
	"github.com/labstack/echo"
	"net/http"
)

func AddCourse(c echo.Context) error {

	err := jwt.TokenValid(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	tokenAuth, err := jwt.ExtractTokenMetadata(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	tid, role, err := jwt.FetchAuth(tokenAuth)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	if role != "teacher" {
		return echo.NewHTTPError(http.StatusUnauthorized, "just teachers can add course")
	}

	u := new(user.Course)

	err = c.Bind(u)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	u.Tid = tid
	u.Tname = db.GetTeacherName(tid)

	if ok, err := db.HasConflict(u.Tid, u.Weekday, u.Time); ok {
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		return echo.NewHTTPError(http.StatusBadRequest, "time conflicted")
	}

	u.Cid, err = db.GetCid(u.Cname)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	u.GenerateUuid()

	err = db.AddCourse(u)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.NoContent(http.StatusOK)
}
