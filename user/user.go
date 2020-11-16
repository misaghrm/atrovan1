package user

import (
	"crypto/rand"
	"errors"
)

type User struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Family   string `json:"family"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

func (u *User) GenerateID() error {
	const (
		otpChars = "1234567890"
		sPrefix  = "99100"
		tPrefix  = "99200"
		length   = 5
	)

	buffer := make([]byte, length)
	_, err := rand.Read(buffer)
	if err != nil {
		return err
	}

	otpCharsLength := len(otpChars)
	for i := 0; i < length; i++ {
		buffer[i] = otpChars[int(buffer[i])%otpCharsLength]
	}

	switch u.Role {
	case "student":
		u.ID = sPrefix + string(buffer)
		return nil
	case "teacher":
		u.ID = tPrefix + string(buffer)
		return nil
	default:
		return errors.New("the role is not defined")
	}

}
