package service

import (
	"log"
	"net/http"
	"strconv"

	"github.com/challengr/model"
	"github.com/gin-gonic/gin"
)

//GetBoughtItem handler func fetches all the bought items of an user
func GetBoughtItem(c *gin.Context) {
	userID, ok := c.MustGet("user_id").(int64)
	if !ok {
		log.Printf("invalid userid in token, userid: %v", c.MustGet("user_id"))
		c.JSON(http.StatusForbidden, &model.ErrResp{Error: "Invalid token"})
		return
	}

	userRole, ok := c.MustGet("role").(string)
	if !ok {
		log.Printf("invalid userole in token, userid: %v", c.MustGet("role"))
		c.JSON(http.StatusForbidden, &model.ErrResp{Error: "Invalid token"})
		return
	}

	paramUserID, err := strconv.ParseInt(c.Param("user_id"), 10, 64)
	if err != nil {
		log.Printf("path parm user_id err: %v", err)
		c.JSON(http.StatusBadRequest, &model.ErrResp{Error: "Invalid path param", Fields: &[]string{"user_id"}})
		return
	}

	if userID != paramUserID && userRole != constAdminRole {
		c.JSON(http.StatusMethodNotAllowed, &model.ErrResp{Error: "Not allowed please check token"})
		return
	}

	boughtItemList, err := (&model.BoughtItem{}).Get("user_id=$1", paramUserID)
	if err != nil {
		log.Printf("db fetching bought item error: %v", err)
		c.JSON(http.StatusInternalServerError, &model.ErrResp{Error: "Server error"})
		return
	}

	c.JSON(http.StatusOK, boughtItemList)
}

//Purchase handler func creates a new bought item record
func Purchase(c *gin.Context) {
	userID, ok := c.MustGet("user_id").(int64)
	if !ok {
		log.Printf("invalid userid in token, userid: %v", c.MustGet("user_id"))
		c.JSON(http.StatusForbidden, &model.ErrResp{Error: "Invalid token"})
		return
	}

	userRole, ok := c.MustGet("role").(string)
	if !ok {
		log.Printf("invalid userole in token, userid: %v", c.MustGet("role"))
		c.JSON(http.StatusForbidden, &model.ErrResp{Error: "Invalid token"})
		return
	}

	paramUserID, err := strconv.ParseInt(c.Param("user_id"), 10, 64)
	if err != nil {
		log.Printf("path parm user_id err: %v", err)
		c.JSON(http.StatusBadRequest, &model.ErrResp{Error: "Invalid path param", Fields: &[]string{"user_id"}})
		return
	}

	if userID != paramUserID && userRole != constAdminRole {
		c.JSON(http.StatusMethodNotAllowed, &model.ErrResp{Error: "Not allowed please check token"})
		return
	}

	var boughtItem model.BoughtItem
	if err := c.BindJSON(&boughtItem); err != nil {
		log.Printf("boughtItem struct JSON bind error: %v", err)
		c.JSON(http.StatusBadRequest, &model.ErrResp{Error: err.Error()})
		return
	}

	if err := c.BindJSON(&boughtItem.Payload); err != nil {
		log.Printf("boughtItem Payload JSON bind error: %v", err)
		c.JSON(http.StatusBadRequest, &model.ErrResp{Error: err.Error()})
		return
	}

	if errSlice := boughtItem.ParseNotAllowedJSON(); len(errSlice) > 0 {
		log.Printf("boughtItem not allowed fields detected: %v", errSlice)
		c.JSON(http.StatusBadRequest, &model.ErrResp{Error: "Some fields are not allowed", Fields: &errSlice})
		return
	}

	if errSlice := boughtItem.PostValidate(); len(errSlice) > 0 {
		log.Printf("boughtItem not allowed fields detected: %v", errSlice)
		c.JSON(http.StatusBadRequest, &model.ErrResp{Error: "Some fields are invalid", Fields: &errSlice})
		return
	}

	if boughtItem.LevelID != nil {
		count, err := (&model.Level{}).Count("WHERE level_id=$1", boughtItem.LevelID)
		if err != nil {
			log.Printf("Level count err: %v", err)
			c.JSON(http.StatusInternalServerError, &model.ErrResp{Error: "Server error"})
			return
		}

		if count != 1 {
			log.Printf("Level count: %v, must be 1", count)
			c.JSON(http.StatusBadRequest, &model.ErrResp{Error: "Invalid fields detected", Fields: &[]string{"level_id"}})
			return
		}
	}

	if boughtItem.VanityItemID != nil {
		count, err := (&model.Level{}).Count("WHERE level_id=$1", boughtItem.VanityItemID)
		if err != nil {
			log.Printf("vanity item count err: %v", err)
			c.JSON(http.StatusInternalServerError, &model.ErrResp{Error: "Server error"})
			return
		}

		if count != 1 {
			log.Printf("vanity item count: %v, must be 1", count)
			c.JSON(http.StatusBadRequest, &model.ErrResp{Error: "Invalid fields detected", Fields: &[]string{"vanity_item_id"}})
			return
		}
	}

	boughtItem.UserID = userID

	if err := boughtItem.Create(); err != nil {
		log.Printf("boughtItem insert error: %v", err)
		c.JSON(http.StatusInternalServerError, &model.ErrResp{Error: "Server error"})
		return
	}

	c.JSON(http.StatusOK, &model.SuccessResp{Message: "Item successfully purchased", Status: http.StatusOK})
}
