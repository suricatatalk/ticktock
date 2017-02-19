package user

import uuid "github.com/satori/go.uuid"

type User struct {
	ID        string `bson:"_id" json:"id"`
	Firstname string `bson:"firstName" json:"firstName"`
	Lastname  string `bson:"lastName" json:"lastName"`
	Accounts  []string
	Email     string `bson:"email" json:"email"`
}

type Users []User

func NewUser() *User {
	return &User{
		ID: uuid.NewV4().String(),
	}
}
