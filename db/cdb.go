package db

import (
	"crypto/rand"
	"github.com/atrovanProject/user"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func AddCourse(c *user.Course) error {
	db, err := gorm.Open(postgres.Open(dbUrl), &gorm.Config{})
	if err != nil {
		return err
	}

	err = db.Table(cTable).Create(&c).Error
	if err != nil {
		return err
	}
	return nil
}

func HasConflict(tid string, weekday, time int) (bool, error) {
	var c user.Course
	db, err := gorm.Open(postgres.Open(dbUrl), &gorm.Config{})
	if err != nil {
		return true, err
	}

	err = db.Table(cTable).Where("tid = ? and weekday = ? and time = ?", tid, weekday, time).Find(&c).Error
	if err != nil {
		return false, err
	}
	if c.Weekday == weekday && c.Time == time {
		return true, nil
	}
	return false, nil
}

func GetCid(cname string) (string, error) {
	var c user.Course
	db, err := gorm.Open(postgres.Open(dbUrl), &gorm.Config{})
	if err != nil {
		return "", err
	}

	err = db.Table(cTable).Select("cid").Where("cname = ?", cname).First(&c).Error
	if err != nil {
		if c.Cid == "" {
			c.Cid, err = generateId()
			if err != nil {
				return "", err
			}
			return c.Cid, nil
		}
	}
	return c.Cid, nil
}

func generateId() (string, error) {
	const (
		otpChars = "1234567890"
		prefix   = "11"
		length   = 4
	)

	buffer := make([]byte, length)
	_, err := rand.Read(buffer)
	if err != nil {
		return "", err
	}

	otpCharsLength := len(otpChars)
	for i := 0; i < length; i++ {
		buffer[i] = otpChars[int(buffer[i])%otpCharsLength]
	}

	return prefix + string(buffer), nil
}

func GetCourses() []user.Course {
	var c []user.Course
	db, err := gorm.Open(postgres.Open(dbUrl), &gorm.Config{})
	if err != nil {
		return nil
	}

	err = db.Table(cTable).Where("reserved <= cap").Find(&c).Error
	if err != nil {
		return nil
	}
	return c
}
