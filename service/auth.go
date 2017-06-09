package service

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"net/http"

	"github.com/challengr/model"
	"github.com/gin-gonic/gin"
)

const fbGraphURI = "https://graph.facebook.com/me?fields=id,email&access_token="

//LogIn func handler logs in a user based on a facebook token and email. Also updates or sets a onesignal detail.
func LogIn(c *gin.Context) {
	var logIn model.LogIn
	if err := c.BindJSON(&logIn); err != nil {
		log.Printf("LogIn struct JSON bind error: %v", err)
		c.JSON(http.StatusBadRequest, &model.ErrResp{Error: "Invalid login payload"})
		return
	}

	if err := c.BindJSON(&logIn.Payload); err != nil {
		log.Printf("LogIn Payload JSON bind error: %v", err)
		c.JSON(http.StatusBadRequest, &model.ErrResp{Error: "Invalid login payload"})
		return
	}

	if errSlice := logIn.ParseNotAllowedJSON(); len(errSlice) > 0 {
		log.Printf("LogIn not allowed fields detected: %v", errSlice)
		c.JSON(http.StatusBadRequest, &model.ErrResp{Error: "Some fields are not allowed", Fields: &errSlice})
		return
	}

	if errSlice := logIn.PostValidate(); len(errSlice) > 0 {
		log.Printf("LogIn post validate err: %v", errSlice)
		c.JSON(http.StatusBadRequest, &model.ErrResp{Error: "Invalid fields detected", Fields: &errSlice})
		return
	}

	url := fbGraphURI + logIn.FacebookToken
	res, err := http.Get(url)
	if err != nil {
		log.Printf("facebook req err: %q", err)
		c.JSON(http.StatusInternalServerError, &model.ErrResp{Error: "Facebook server error"})
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("facebook resp io read err: %q", err)
		c.JSON(http.StatusInternalServerError, &model.ErrResp{Error: "Server error"})
		return
	}
	var facebookResp model.FacebookResp
	if err = json.Unmarshal(body, &facebookResp); err != nil {
		log.Printf("Facebook body byte unmarshalling err: %q", err)
		c.JSON(http.StatusInternalServerError, &model.ErrResp{Error: "Server error"})
		return
	}

	if facebookResp.Email != logIn.Email {
		log.Printf("LogIn facebook email: %v != given login emai: %v", facebookResp.Email, logIn.Email)
		c.JSON(http.StatusBadRequest, &model.ErrResp{Error: "Invalid fields detected", Fields: &[]string{"email"}})
		return
	}

	user := model.User{Email: &facebookResp.Email, FacebookUserID: &facebookResp.ID}

	count, err := user.Count("WHERE email=$1 AND facebook_user_id=$2", user.Email, facebookResp.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, &model.ErrResp{Error: "Server error"})
		return
	}

	if count != 0 {
		user.Token = user.CreateTokenString()
		c.JSON(http.StatusOK, &user)
		return
	}

	if err = user.Create(); err != nil {
		c.JSON(http.StatusBadRequest, &model.ErrResp{Error: "Server error"})
		return
	}

	user.Token = user.CreateTokenString()
	c.JSON(http.StatusOK, &user)
}

//LogOut func handler logs out a user based on an user_id and imei
func LogOut(c *gin.Context) {
	var logOut model.LogOut
	if err := c.BindJSON(&logOut); err != nil {
		log.Printf("logOut struct JSON bind error: %v", err)
		c.JSON(http.StatusBadRequest, &model.ErrResp{Error: "Invalid logOut payload"})
		return
	}

	if err := c.BindJSON(&logOut.Payload); err != nil {
		log.Printf("logOut Payload JSON bind error: %v", err)
		c.JSON(http.StatusBadRequest, &model.ErrResp{Error: "Invalid logOut payload"})
		return
	}

	if errSlice := logOut.ParseNotAllowedJSON(); len(errSlice) > 0 {
		log.Printf("logOut not allowed fields detected: %v", errSlice)
		c.JSON(http.StatusBadRequest, &model.ErrResp{Error: "Some fields are not allowed", Fields: &errSlice})
		return
	}

	if errSlice := logOut.Validate(); len(errSlice) > 0 {
		log.Printf("logOut post validate err: %v", errSlice)
		c.JSON(http.StatusBadRequest, &model.ErrResp{Error: "Invalid fields detected", Fields: &errSlice})
		return
	}

	logOut.UserID = c.MustGet("user_id").(int64)

	oneSignal := model.OneSignal{UserID: logOut.UserID, Imei: logOut.Imei}

	if err := oneSignal.Delete(); err != nil {
		log.Printf("Delete onesignal error: %v", err)
		c.JSON(http.StatusInternalServerError, &model.ErrResp{Error: "Server error"})
		return
	}

	c.JSON(http.StatusOK, &model.SuccessResp{Status: http.StatusOK, Message: "Successfully logged out"})
}
