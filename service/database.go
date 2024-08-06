package service

import (
	"github.com/gin-gonic/gin"
	"time"
	"users/domain"
	"users/ecode"
)

type DataBaseRepository interface {
	AddUser(c *gin.Context, user *domain.Users) (int64, error)
	GetUserByID(c *gin.Context, id int64) (*domain.Users, error)
	GetUsersByDateAndAgeRange(c *gin.Context, startDate, endDate *time.Time, minAge, maxAge *int) ([]domain.Users, int64, error)
}

type DataBaseWorker struct {
	repo DataBaseRepository
}

func NewDataBaseWorker(repo DataBaseRepository) *DataBaseWorker {
	return &DataBaseWorker{
		repo: repo,
	}
}

func (dw *DataBaseWorker) Add(c *gin.Context, user *domain.Users) (int64, error) { //
	userID, err := dw.repo.AddUser(c, user)
	if err != nil {
		return 0, ecode.ErrWriteDB
	}

	return userID, nil
}

func (dw *DataBaseWorker) GetUserByID(c *gin.Context, id int64) (*domain.Users, error) {
	return dw.repo.GetUserByID(c, id)
}

func (dw *DataBaseWorker) GetUsersByDateAndAgeRange(c *gin.Context, startDate, endDate *time.Time, minAge, maxAge *int) ([]domain.Users, int64, error) {
	return dw.repo.GetUsersByDateAndAgeRange(c, startDate, endDate, minAge, maxAge)
}
