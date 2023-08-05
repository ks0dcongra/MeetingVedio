package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Id       int
	Name     string
	Email    string
	Password string
}

type MongoUser struct {
	Id       int `json:"id"`
	Name     string `json:"name"`
	Password    string `json:"password"`
	Age int `json:"age"`
	Email string `json:"email"`
}

