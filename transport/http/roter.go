package http

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"users/config"
)

type Router interface {
	Start()
	RegisterRoutes()
}

type RouterImpl struct {
	userController *UserController
	logger         *zap.Logger
	config         *config.Config
	server         *gin.Engine
}

func NewRouter(userController *UserController, logger *zap.Logger, cfg *config.Config) *RouterImpl {
	return &RouterImpl{
		userController: userController,
		logger:         logger,
		config:         cfg,
	}
}

func (r *RouterImpl) RegisterRoutes() {
	r.logger.Info("Registering routes")

	router := gin.Default()

	router.POST("/users", func(c *gin.Context) {
		r.userController.AddUser(c)
	})

	router.GET("/users", func(c *gin.Context) {
		r.userController.GetUserByID(c)
	})

	router.GET("/rangeUsersDateAge", func(c *gin.Context) {
		r.userController.GetUsersByDateAndAgeRange(c)
	})

	r.server = router

}

func (r *RouterImpl) Start() error {
	port := r.config.LocalURL
	return r.server.Run(port)
}
