package be

import (
	"github.com/coocood/qbs"
)

type User struct {
	Id        int64  `json:"id" qbs:"pk"`
	Email     string `json:"email"`
	FirstName string `json:"first-name"`
	LastName  string `json:"last-name"`
	Staff     bool   `json:"staff"`
}

func CreateUser(email string, firstName string, lastName string, staff bool, db *qbs.Qbs) (*User, error) {
	user := new(User)
	user.Email = email
	user.FirstName = firstName
	user.LastName = lastName
	user.Staff = staff
	_, err := db.Save(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func UpdateUser(user *User, db *qbs.Qbs) error {
	_, err := db.Save(user)
	if err != nil {
		return err
	}
	return nil
}

func FindUser(id int64, db *qbs.Qbs) (*User, error) {
	user := new(User)
	user.Id = id
	err := db.Find(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}
