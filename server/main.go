package main

import (
	"fmt"
	"log"
	"os"
	"server/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	r := gin.Default()
	err := godotenv.Load()
	routes.SetupRouter(r)
	r.Use(cors.Default())
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	var PORT = os.Getenv("PORT")
	r.Run(":" + PORT)

	fmt.Println("Hello World")
}
