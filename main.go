package main

import (
	"log"
	"net/http"
	"tusk/config"
	"tusk/controllers"
	"tusk/models"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	//Database
	db := config.DatabaseConnection()
	db.AutoMigrate(&models.User{}, &models.Task{})
	config.CreatedOwnerAccount(db)

	//Controllers
	userController := controllers.UserController{DB: db}

	//Routes
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, "Welcome to tusk API")
	})

	router.POST("/users/login", userController.Login)
	router.POST("/users", userController.CreateAccount)
	router.DELETE("/users/:id", userController.Delete)

	router.Static("/attachment", "./attachment")
	router.Run("192.168.158.28:8080")
}
