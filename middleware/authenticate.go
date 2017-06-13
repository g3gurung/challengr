package middleware

import (
	"log"
	"time"

	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

//JWTSecret is user for encrypting and decrypting jwt
const JWTSecret = "challengr app secret"

//respondWithError func responds with error
func respondWithError(code int, message string, c *gin.Context) {
	resp := map[string]string{"error": message}

	c.JSON(code, resp)
	c.AbortWithStatus(code)
}

//Authenticate func middleware authenticates incoming request
func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()
		tokenstring := c.Query("token")
		if tokenstring == "" {
			tokenstring = c.Request.Header.Get("access-token")
		}

		if tokenstring == "" {
			respondWithError(http.StatusForbidden, "token required", c)
			return
		}

		user := JWTUser{}
		token, err := jwt.ParseWithClaims(tokenstring, &user, func(token *jwt.Token) (interface{}, error) {
			return []byte(JWTSecret), nil
		})

		if err != nil {
			log.Printf("token parse err: %v", err)
			respondWithError(http.StatusForbidden, "invalid  token", c)
			return
		}

		if !token.Valid {
			log.Println("token.Valid false")
			respondWithError(http.StatusForbidden, "invalid  token", c)
			return
		}

		// Set example variable
		c.Set("user_id", user.ID)
		c.Set("facebook_user_id", user.FacebookUserID)
		c.Set("weight", user.Weight)
		c.Set("role", user.Role)

		// before request

		c.Next()

		// after request
		latency := time.Since(t)
		log.Print(latency)

		// access the status we are sending
		status := c.Writer.Status()
		log.Println(status)
	}
}
