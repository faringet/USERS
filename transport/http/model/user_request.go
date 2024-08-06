package model

import (
	"gorm.io/gorm"
	"time"
	"users/domain"
)

type UserRequestAdd struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Age       int    `json:"age"`
}

func (u UserRequestAdd) MapToDomain() domain.Users {
	return domain.Users{
		Model:         gorm.Model{},
		FirstName:     u.FirstName,
		LastName:      u.LastName,
		Age:           u.Age,
		RecordingDate: time.Now(),
	}
}
