package main

import (
	"log"
	"trip-planner-be/internal/config"
	"trip-planner-be/internal/handlers"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	config.SetupGoogleOAuth()

	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/", handlers.HandleHome)
	e.GET("/login", handlers.HandleLogin)
	e.GET("/auth/google/callback", handlers.HandleCallback)

	// Start server
	log.Fatal(e.Start(":3000"))
}
