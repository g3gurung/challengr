package main

import (
	"log"
	"os"

	"github.com/challengr/service"
	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	router := gin.New()
	router.Use(gin.Logger())

	log.SetFlags(log.Ldate | log.Ltime | log.LUTC | log.Llongfile)

	router.GET("/s3Sign", service.PreSignS3)

	router.Run(":" + port)
}
