package service

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/challengr/model"
	"github.com/gin-gonic/gin"
)

//AddCoins func handler adds coins to user score
func AddCoins(c *gin.Context) {
	userID, ok := c.MustGet("user_id").(int64)
	if !ok {
		log.Printf("invalid userid in token, userid: %v", c.MustGet("user_id"))
		c.JSON(http.StatusForbidden, &model.ErrResp{Error: "Invalid token"})
		return
	}

	paramScoreID, err := strconv.ParseInt(c.Param("score_id"), 10, 64)
	if err != nil {
		log.Printf("path parm user_id err: %v", err)
		c.JSON(http.StatusBadRequest, &model.ErrResp{Error: "Invalid path param", Fields: &[]string{"user_id"}})
		return
	}

	m := make(map[string]interface{})
	if err := c.BindJSON(&m); err != nil {
		log.Printf("add coins struct JSON bind error: %v", err)
		c.JSON(http.StatusBadRequest, &model.ErrResp{Error: "Invalid payload"})
		return
	}

	errSlice := []string{}
	for key := range m {
		if key != "amount" {
			errSlice = append(errSlice, key)
		}
	}

	if len(errSlice) > 0 {
		c.JSON(http.StatusBadRequest, &model.ErrResp{Error: "Some fields are not allowed", Fields: &errSlice})
		return
	}

	jsonString, err := json.Marshal(m)
	if err != nil {
		log.Printf("add coins struct JSON bind error: %v", err)
		c.JSON(http.StatusBadRequest, &model.ErrResp{Error: "Invalid payload"})
		return
	}
	amount := model.Amount{}
	json.Unmarshal(jsonString, &amount)

	if amount.Amount == 0 {
		log.Printf("add coins struct JSON bind error: %v", err)
		c.JSON(http.StatusBadRequest, &model.ErrResp{Error: "Missing or invalid payload field", Fields: &[]string{"amount"}})
		return
	}

	var status int
	if status, err := (&model.Score{ID: paramScoreID, UserID: userID}).AddCoins(amount.Amount); err != nil {
		log.Printf("add coinsdb error: %v", err)
		c.JSON(status, &model.ErrResp{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, &model.SuccessResp{Message: "Coins successfully added", Status: http.StatusOK})
}

//AddExp func handler adds experiences to user score
func AddExp(c *gin.Context) {
	userID, ok := c.MustGet("user_id").(int64)
	if !ok {
		log.Printf("invalid userid in token, userid: %v", c.MustGet("user_id"))
		c.JSON(http.StatusForbidden, &model.ErrResp{Error: "Invalid token"})
		return
	}

	paramScoreID, err := strconv.ParseInt(c.Param("score_id"), 10, 64)
	if err != nil {
		log.Printf("path parm user_id err: %v", err)
		c.JSON(http.StatusBadRequest, &model.ErrResp{Error: "Invalid path param", Fields: &[]string{"user_id"}})
		return
	}

	m := make(map[string]interface{})
	if err := c.BindJSON(&m); err != nil {
		log.Printf("add exp struct JSON bind error: %v", err)
		c.JSON(http.StatusBadRequest, &model.ErrResp{Error: "Invalid payload"})
		return
	}

	errSlice := []string{}
	for key := range m {
		if key != "amount" {
			errSlice = append(errSlice, key)
		}
	}

	if len(errSlice) > 0 {
		c.JSON(http.StatusBadRequest, &model.ErrResp{Error: "Some fields are not allowed", Fields: &errSlice})
		return
	}

	jsonString, err := json.Marshal(m)
	if err != nil {
		log.Printf("add exp struct JSON bind error: %v", err)
		c.JSON(http.StatusBadRequest, &model.ErrResp{Error: "Invalid payload"})
		return
	}
	amount := model.Amount{}
	json.Unmarshal(jsonString, &amount)

	if amount.Amount == 0 {
		log.Printf("add exp struct JSON bind error: %v", err)
		c.JSON(http.StatusBadRequest, &model.ErrResp{Error: "Missing or invalid payload field", Fields: &[]string{"amount"}})
		return
	}

	var status int
	if status, err := (&model.Score{ID: paramScoreID, UserID: userID}).AddExp(amount.Amount); err != nil {
		log.Printf("add exp db error: %v", err)
		c.JSON(status, &model.ErrResp{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, &model.SuccessResp{Message: "Exp successfully added", Status: http.StatusOK})
}

//AddLikes func handler adds likes to user score
func AddLikes(c *gin.Context) {
	userID, ok := c.MustGet("user_id").(int64)
	if !ok {
		log.Printf("invalid userid in token, userid: %v", c.MustGet("user_id"))
		c.JSON(http.StatusForbidden, &model.ErrResp{Error: "Invalid token"})
		return
	}

	paramScoreID, err := strconv.ParseInt(c.Param("score_id"), 10, 64)
	if err != nil {
		log.Printf("path parm user_id err: %v", err)
		c.JSON(http.StatusBadRequest, &model.ErrResp{Error: "Invalid path param", Fields: &[]string{"user_id"}})
		return
	}

	m := make(map[string]interface{})
	if err := c.BindJSON(&m); err != nil {
		log.Printf("add exp struct JSON bind error: %v", err)
		c.JSON(http.StatusBadRequest, &model.ErrResp{Error: "Invalid payload"})
		return
	}

	errSlice := []string{}
	for key := range m {
		if key != "amount" {
			errSlice = append(errSlice, key)
		}
	}

	if len(errSlice) > 0 {
		c.JSON(http.StatusBadRequest, &model.ErrResp{Error: "Some fields are not allowed", Fields: &errSlice})
		return
	}

	jsonString, err := json.Marshal(m)
	if err != nil {
		log.Printf("add likes struct JSON bind error: %v", err)
		c.JSON(http.StatusBadRequest, &model.ErrResp{Error: "Invalid payload"})
		return
	}
	amount := model.Amount{}
	json.Unmarshal(jsonString, &amount)

	if amount.Amount == 0 {
		log.Printf("add likes struct JSON bind error: %v", err)
		c.JSON(http.StatusBadRequest, &model.ErrResp{Error: "Missing or invalid payload field", Fields: &[]string{"amount"}})
		return
	}

	var status int
	if status, err := (&model.Score{ID: paramScoreID, UserID: userID}).AddLikes(amount.Amount); err != nil {
		log.Printf("add likes db error: %v", err)
		c.JSON(status, &model.ErrResp{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, &model.SuccessResp{Message: "Likes successfully added", Status: http.StatusOK})
}

/*
//ResetLikes func handler resets the likes
func ResetLikes(c *gin.Context) {

}
*/
