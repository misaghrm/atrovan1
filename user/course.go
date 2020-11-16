package user

import "github.com/twinj/uuid"

type Course struct {
	Cid      string `json:"cid"`
	Cname    string `json:"cname"`
	Tname    string `json:"-"`
	Tid      string `json:"-"`
	Weight   int    `json:"weight"`
	Weekday  int    `json:"weekday"`
	Time     int    `json:"time"`
	Cap      int    `json:"cap"`
	Reserved int    `json:"reserved"`
	Uuid     string `json:"uuid"`
}

func (c *Course) GenerateUuid() {
	c.Uuid = uuid.NewV4().String()
}
