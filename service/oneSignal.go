package service

import (
	"log"
	"net/http"

	"github.com/challengr/model"
	"github.com/gin-gonic/gin"
)

//UpdateOneSignal func is a handler for updating onesignal account info of an user
func UpdateOneSignal(c *gin.Context) {
	var oneSignal model.OneSignal
	if err := c.BindJSON(&oneSignal); err != nil {
		log.Printf("oneSignal struct JSON bind error: %v", err)
		c.JSON(http.StatusBadRequest, &model.ErrResp{Error: err.Error()})
		return
	}

	if err := c.BindJSON(&oneSignal.Payload); err != nil {
		log.Printf("oneSignal Payload JSON bind error: %v", err)
		c.JSON(http.StatusBadRequest, &model.ErrResp{Error: err.Error()})
		return
	}

	if errSlice := oneSignal.ParseNotAllowedJSON(); len(errSlice) > 0 {
		log.Printf("LogIn not allowed fields detected: %v", errSlice)
		c.JSON(http.StatusBadRequest, &model.ErrResp{Error: "Some fields are not allowed", Fields: &errSlice})
		return
	}

	if errSlice := oneSignal.Validate(); len(errSlice) > 0 {
		log.Printf("LogIn post validate err: %v", errSlice)
		c.JSON(http.StatusBadRequest, &model.ErrResp{Error: "Invalid fields detected", Fields: &errSlice})
		return
	}

	if err := oneSignal.Upsert(); err != nil {
		log.Printf("oneSignal upsert error: %v", err)
		c.JSON(http.StatusInternalServerError, &model.ErrResp{Error: "Server errror"})
		return
	}

	c.JSON(http.StatusOK, &model.SuccessResp{Status: http.StatusOK, Message: "Success"})
}
