package router

import (
	"go-crud/initializers"
	"go-crud/middleware"
	"go-crud/views"
	"time"

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

	// Add global middlewares
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.LoggerMiddleware())
	router.Use(middleware.ErrorLoggerMiddleware())

	// Rate limiting for public endpoints (100ms between requests, burst of 5)
	publicRateLimiter := middleware.RateLimitMiddleware(100*time.Millisecond, 5)

	// Public endpoints with rate limiting
	public := router.Group("/")
	public.Use(publicRateLimiter)
	{
		public.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status":  "healthy",
				"service": "go-crud-api",
			})
		})

		// Swagger endpoint (public but no rate limiting for docs)
		public.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	postViews := views.NewPostViews()
	postViews.RegisterRoutes(router)

	versionViews := views.NewPostVersionViews()
	versionViews.RegisterRoutes(router)

	userViews := views.NewUserViews()
	userViews.RegisterRoutes(router)

	authViews := views.NewAuthViews()
	authViews.RegisterRoutes(router)

	tagViews := views.NewTagViews()
	tagViews.RegisterRoutes(router)

	return router
}
