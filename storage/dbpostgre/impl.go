package dbpostgre

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"time"
	"users/domain"
)

type DataBaseRepositoryImpl struct {
	postgreClient *gorm.DB
}

func NewDataBaseRepositoryImpl(postgreClient *gorm.DB) *DataBaseRepositoryImpl {
	return &DataBaseRepositoryImpl{
		postgreClient: postgreClient,
	}
}

func (dr *DataBaseRepositoryImpl) AddUser(c *gin.Context, user *domain.Users) (int64, error) {
	result := dr.postgreClient.Create(user)

	if result.Error != nil {
		return 0, result.Error
	}

	return user.ID, nil
}

func (dr *DataBaseRepositoryImpl) GetUserByID(c *gin.Context, id int64) (*domain.Users, error) {
	var user domain.Users
	result := dr.postgreClient.First(&user, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (dr *DataBaseRepositoryImpl) GetUsersByDateAndAgeRange(c *gin.Context, startDate, endDate *time.Time, minAge, maxAge *int) ([]domain.Users, int64, error) {
	var users []domain.Users
	query := dr.postgreClient.Model(&domain.Users{})

	if startDate != nil {
		query = query.Where("recording_date >= ?", *startDate)
	}
	if endDate != nil {
		query = query.Where("recording_date <= ?", *endDate)
	}
	if minAge != nil {
		query = query.Where("age >= ?", *minAge)
	}
	if maxAge != nil {
		query = query.Where("age <= ?", *maxAge)
	}

	result := query.Find(&users)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return users, result.RowsAffected, nil
}
