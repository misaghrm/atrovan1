package db

import (
	"errors"
	"github.com/atrovanProject/user"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Reserved struct {
	Sid     string `json:"sid"`
	Uuid    string `json:"uuid"`
	Tname   string `json:"Tname"`
	Cname   string `json:"cname"`
	Weekday int    `json:"weekday"`
	Time    int    `json:"time"`
}

func AddReserve(r *Reserved) error {
	var c user.Course
	var d Reserved
	db, err := gorm.Open(postgres.Open(dbUrl), &gorm.Config{})
	if err != nil {
		return err
	}

	err = db.Table(cTable).Select("cap", "reserved", "weekday", "time", "tname", "cname").Where("uuid = ?", r.Uuid).Find(&c).Error
	err = db.Table(rTable).Select("weekday", "time").Where("uuid = ? and weekday = ? and time = ?", r.Uuid, c.Weekday, c.Time).Find(&d).Error
	if d.Weekday == c.Weekday && d.Time == c.Time {
		return errors.New("there is a time conflict")
	}
	if c.Reserved >= c.Cap {
		return errors.New("capacity is full")
	}
	c.Reserved++
	err = db.Table(cTable).Where("uuid = ?", r.Uuid).UpdateColumn("reserved", c.Reserved).Error
	if err != nil {
		return err
	}

	r.Tname = c.Tname
	r.Cname = c.Cname
	r.Weekday = c.Weekday
	r.Time = c.Time

	err = db.Table(rTable).Create(&r).Error
	if err != nil {
		return err
	}
	return nil
}

func GetReservedList(sid string) ([]Reserved, error) {
	var r []Reserved
	db, err := gorm.Open(postgres.Open(dbUrl), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.Table(rTable).Where("sid = ?", sid).Find(&r).Error
	if err != nil {
		return nil, err
	}
	return r, nil
}
