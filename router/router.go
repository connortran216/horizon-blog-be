package router

import (
	"go-crud/initializers"
	"go-crud/views"

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	swaggerFiles "github.com/swaggo/files"
)

// SetupRouter creates and configures the Gin router
func SetupRouter() *gin.Engine {
	// Initialize dependencies
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()

	router := gin.Default()

	postViews := views.NewPostViews()
	postViews.RegisterRoutes(router)

	userViews := views.NewUserViews()
	userViews.RegisterRoutes(router)

	authViews := views.NewAuthViews()
	authViews.RegisterRoutes(router)

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"service": "go-crud-api",
		})
	})

	// Swagger endpoint
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}
