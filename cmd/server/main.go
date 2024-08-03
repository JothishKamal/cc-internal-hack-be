package main

import (
	"fmt"
	"log"
	"net/http"

	"trip-planner-be/internal/config"
	"trip-planner-be/internal/handlers"
)

func main() {
	config.LoadEnv()
	config.SetupGoogleOAuth()

	http.HandleFunc("/", handlers.HandleHome)
	http.HandleFunc("/login", handlers.HandleLogin)
	http.HandleFunc("/callback", handlers.HandleCallback)

	fmt.Println("Started server on :3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
