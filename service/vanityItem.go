package service

import (
	"log"

	"net/http"

	"github.com/challengr/model"
	"github.com/gin-gonic/gin"
)

//GetVanityItem handler func fetches the list of
func GetVanityItem(c *gin.Context) {
	vanityItemList, err := (&model.VanityItem{}).Get("")
	if err != nil {
		log.Printf("db fetching vanityItmeList error: %v", err)
		c.JSON(http.StatusInternalServerError, &model.ErrResp{Error: "Server error"})
		return
	}

	c.JSON(http.StatusOK, vanityItemList)
}
