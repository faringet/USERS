package http

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"time"
	"users/domain"
	"users/transport/http/model"
)

type Database interface {
	Add(c *gin.Context, request *domain.Users) (int64, error)
	GetUserByID(c *gin.Context, id int64) (*domain.Users, error)
	GetUsersByDateAndAgeRange(c *gin.Context, startDate, endDate *time.Time, minAge, maxAge *int) ([]domain.Users, int64, error)
}

type UserController struct {
	database Database
	logger   *zap.Logger
}

func NewUserController(logger *zap.Logger, database Database) *UserController {
	return &UserController{
		database: database,
		logger:   logger,
	}
}

func (uc *UserController) AddUser(c *gin.Context) {
	var request model.UserRequestAdd

	if err := c.ShouldBindJSON(&request); err != nil {
		uc.logger.Error("invalid input format")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := request.MapToDomain()

	userID, err := uc.database.Add(c, &user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": userID})
	uc.logger.Info(fmt.Sprintf("user's id:%d", userID))
}

func (uc *UserController) GetUserByID(c *gin.Context) {
	idParam := c.Query("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		uc.logger.Error("invalid user id")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	user, err := uc.database.GetUserByID(c, id)
	if err != nil {
		uc.logger.Error("user not found")
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (uc *UserController) GetUsersByDateAndAgeRange(c *gin.Context) {
	var startDate, endDate *time.Time
	var minAge, maxAge *int

	if start := c.Query("start_date"); start != "" {
		parsedStartDate, err := time.Parse(time.RFC3339, start)
		if err != nil {
			uc.logger.Error("invalid start_date format")
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start_date format"})
			return
		}
		startDate = &parsedStartDate
	}

	if end := c.Query("end_date"); end != "" {
		parsedEndDate, err := time.Parse(time.RFC3339, end)
		if err != nil {
			uc.logger.Error("invalid end_date format")
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end_date format"})
			return
		}
		endDate = &parsedEndDate
	}

	if min := c.Query("min_age"); min != "" {
		parsedMinAge, err := strconv.Atoi(min)
		if err != nil {
			uc.logger.Error("invalid min_age format")
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid min_age format"})
			return
		}
		minAge = &parsedMinAge
	}

	if max := c.Query("max_age"); max != "" {
		parsedMaxAge, err := strconv.Atoi(max)
		if err != nil {
			uc.logger.Error("invalid max_age format")
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid max_age format"})
			return
		}
		maxAge = &parsedMaxAge
	}

	users, count, err := uc.database.GetUsersByDateAndAgeRange(c, startDate, endDate, minAge, maxAge)
	if err != nil {
		uc.logger.Error("failed to fetch users")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch users"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"count": count, "users": users})
}
