//go:generate swag init

package main

import (
	_ "go-crud/docs" // This will be generated
	"go-crud/router"
)

// @title Go CRUD API
// @version 1.0
// @description A simple CRUD API for posts built with Go and Gin
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /
// @schemes http https

func main() {
	r := router.SetupRouter()
	r.Run() // listen and serve on 0.0.0.0:8080
}
