package db

import (
	"errors"
	"github.com/atrovanProject/user"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
)

var (
	dbUrl  = "postgresql://misagh:13750620@127.0.0.1:5432/atrovan"
	sTable = "students"
	tTable = "teachers"
	cTable = "courses"
	rTable = "reserved"
)

func dbTable(role string) (string, error) {
	if role == "student" {
		return sTable, nil
	} else if role == "teacher" {
		return tTable, nil
	} else {
		return "", errors.New("role is undefined")
	}
}
func Register(u *user.User) error {
	dbTable, err := dbTable(u.Role)
	if err != nil {
		return err
	}

	db, err := gorm.Open(postgres.Open(dbUrl), &gorm.Config{})
	if err != nil {
		return err
	}

	err = db.Table(dbTable).Omit("role").Create(&u).Error
	if err != nil {
		return err
	}
	return nil
}

func IsRegistered(id, role string) bool {
	var u user.User
	dbTable, err := dbTable(role)
	if err != nil {
		return false
	}

	db, err := gorm.Open(postgres.Open(dbUrl), &gorm.Config{})
	if err != nil {
		os.Exit(1)
	}

	err = db.Table(dbTable).Where("id = ?", id).Find(&u).Error
	if err != nil {
		log.Fatalln(err)
		return false
	}

	if u.ID != id {
		return false
	}

	return true
}

func Password(id, role string) (string, error) {
	var u user.User
	dbTable, err := dbTable(role)
	if err != nil {
		return "", err
	}

	db, err := gorm.Open(postgres.Open(dbUrl), &gorm.Config{})
	if err != nil {
		return "", err
	}

	err = db.Table(dbTable).Select("password").Where("id = ?", id).Scan(&u).Error
	if err != nil {
		return "", err
	}
	return u.Password, nil
}

func GetTeacherName(id string) string {
	var u user.User
	db, err := gorm.Open(postgres.Open(dbUrl), &gorm.Config{})
	if err != nil {
		return ""
	}

	err = db.Table(tTable).Select("name", "family").Where("id = ?", id).Find(&u).Error
	if err != nil {
		return ""
	}
	return u.Name + " " + u.Family

}
