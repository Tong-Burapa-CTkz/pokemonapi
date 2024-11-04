package service

import (
	"fmt"
	"net/http"

	middleware "sms2pro/internal/middleware"
	models "sms2pro/internal/models"

	service "sms2pro/internal/database"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func Register(c *gin.Context) {

	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	user.ID = uuid.New()

	hashedPassword, err := middleware.HashPassword(user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	user.Password = hashedPassword

	fmt.Println(user.Password)

	// Get the DB instance from the service
	db := service.New().GetDB()

	// Insert the user into the database
	if err := db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}
