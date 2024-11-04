package authen

import (
	"fmt"
	"net/http"
	"os"
	"sms2pro/internal/models"
	"time"

	service "sms2pro/internal/database"
	middleware "sms2pro/internal/middleware"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var jwtKey = []byte(os.Getenv("SecretKey"))

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func Authenticate(c *gin.Context) {

	username := c.Request.Header.Get("Username")
	password := c.Request.Header.Get("Password")

	db := service.New().GetDB()

	var existingUser models.User

	if err := db.Where("username = ?", username).First(&existingUser).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	if !middleware.CheckPasswordHash(password, existingUser.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token, err := GenerateToken(existingUser.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func GenerateToken(username string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour) // Token expiration time (e.g., 24 hours)

	claims := &Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateToken(c *gin.Context) {
	tokenString := c.Request.Header.Get("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token required"})
		c.Abort()
		return
	}

	// Remove "Bearer " prefix if it exists
	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:]
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unautorize"})
		c.Abort()
		return
	}
	c.Next()
}
