package domain

import (
	"gorm.io/gorm"
	"time"
)

type Users struct {
	gorm.Model
	ID            int64
	FirstName     string
	LastName      string
	Age           int
	RecordingDate time.Time
}
