package server

import (
	"net/http"

	auth "sms2pro/internal/authen"
	services "sms2pro/internal/service"

	"github.com/gin-gonic/gin"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()

	// r.POST("/register", controllers.Register)
	// r.POST("/login", controllers.Login)

	// auth := r.Group("/")
	// auth.Use(controllers.AuthMiddleware)

	r.GET("/", s.HelloWorldHandler)
	r.GET("/health", s.healthHandler)
	r.POST("/register", services.Register)
	r.POST("/login", auth.Authenticate)
	// auth := r.Group("/")
	// // r.GET("/api/set", service.SavePokemon)
	// // r.GET("/api/get", service.FetchPokemon)
	// auth.GET("/pokemon/:name", services.GetPokemon)
	// auth.GET("/pokemon/:name/ability", services.GetPokemonAbility)

	authGroup := r.Group("/pokemon")
	authGroup.Use(auth.ValidateToken) // Apply the token validation middleware
	{
		authGroup.GET("/:name", services.GetPokemon)
		authGroup.GET("/:name/ability", services.GetPokemonAbility)
	}

	return r
}

func (s *Server) HelloWorldHandler(c *gin.Context) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	c.JSON(http.StatusOK, resp)
}

func (s *Server) healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, s.db.Health())
}
