package app

import (
	"medods/internal/adapters"
	"medods/internal/adapters/database"
	"medods/internal/adapters/handlers"
	"medods/internal/service"

	"github.com/gin-gonic/gin"
)

func Initialize() *gin.Engine {
	router := gin.Default()
	configs := adapters.Configuration()
	mongoDB := database.NewDatabase(configs)
	storage := database.NewCollection(mongoDB, configs.NameCollection)
	authService := service.NewService(storage)
	handler := handlers.NewHandlers(authService)
	handler.Register(router)
	return router
}
