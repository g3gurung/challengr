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

	router.PUT("/onesignal", service.UpdateOneSignal)

	router.POST("/login", service.LogIn)
	router.POST("/logout", service.LogOut)

	router.GET("/user", service.GetUser)

	router.GET("/challenge", service.GetChellenge)
	router.POST("/challenge", service.PostChallenge)
	router.PUT("/challenge/:challenge_id", service.PutChallenge)
	router.DELETE("/challenge/:challenge_id", service.DeleteChallenge)

	router.GET("/challenge/:challenge_id/post/:post", service.GetPost)

	router.Run(":" + port)
}
