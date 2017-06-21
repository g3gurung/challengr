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
		port = "8080"
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	log.SetFlags(log.Ldate | log.Ltime | log.LUTC | log.Llongfile)

	router.GET("/s3Sign", service.PreSignS3)

	router.PUT("/onesignal", service.UpdateOneSignal)

	router.POST("/login", service.LogIn)
	router.POST("/logout", service.LogOut)

	router.GET("/user", service.GetUser)
	router.PUT("/user/:user_id/weight", service.UpdateUserWeight)
	router.PUT("/user/:user_id/level", service.UpdateUserLevel)

	router.PUT("/user/:user_id/score/:score_id/add_coins", service.AddCoins)
	router.PUT("/user/:user_id/score/:score_id/add_exp", service.AddExp)
	router.PUT("/user/:user_id/score/:score_id/add_likes", service.AddLikes)

	router.GET("/challenge", service.GetChellenge)
	router.POST("/challenge", service.PostChallenge)
	router.PUT("/challenge/:challenge_id", service.PutChallenge)
	router.DELETE("/challenge/:challenge_id", service.DeleteChallenge)

	router.GET("/challenge/:challenge_id/post", service.GetPost)
	router.POST("/challenge/:challenge_id/post", service.PostPost)
	router.PUT("/challenge/:challenge_id/post/:post_id/like", service.LikePost)
	router.PUT("/challenge/:challenge_id/post/:post_id/flag", service.FlagPost)
	router.PUT("/challenge/:challenge_id/post/:post_id/unflag", service.UnFlagPost)
	router.DELETE("/challenge/:challenge_id/post/:post_id", service.DeletePost)

	router.Run(":" + port)
}
