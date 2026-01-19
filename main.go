package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/kartikgujjar22/tokenvault/internal/api"
	"github.com/kartikgujjar22/tokenvault/internal/database"
)

func main() {
	//Initialize the V2 Database (Encryption + Migrations)
	fmt.Println("TokenVault V2: Initializing Secure Database...")
	if err := database.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Setup the Router
	r := gin.Default()

	// Define Routes
	r.POST("/store", api.HandleStoreToken)
	r.GET("/fetch/:project", api.HandleFetchToken)

	// Start Server
	fmt.Println("TokenVault is running on http://localhost:9999")
	if err := r.Run(":9999"); err != nil {
		log.Fatal(err)
	}
}
