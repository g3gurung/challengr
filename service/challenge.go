package service

import (
	"net/http"
	"strconv"

	"strings"

	"log"

	"github.com/challengr/model"
	"github.com/gin-gonic/gin"
)

//GetChellenge func handler fetches challenges
func GetChellenge(c *gin.Context) {
	whereClause := []string{}
	preparedValues := make(map[int]interface{})

	queryUserID := c.Query("user_id")
	if queryUserID != "" {
		userID, err := strconv.ParseInt(queryUserID, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, &model.ErrResp{Error: "Invalid query strings", Fields: &[]string{"user_id"}})
			return
		}
		preparedValues[len(whereClause)] = userID
		whereClause = append(whereClause, "user_id=$"+strconv.Itoa(len(whereClause)+1))
	}

	queryLastID := c.Query("last_id")
	if queryLastID != "" {
		lastID, err := strconv.ParseInt(queryLastID, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, &model.ErrResp{Error: "Invalid query strings", Fields: &[]string{"last_id"}})
			return
		}
		preparedValues[len(whereClause)] = lastID
		whereClause = append(whereClause, "last_id>$"+strconv.Itoa(len(whereClause)+1))
	} else {
		preparedValues[len(whereClause)] = 0
		whereClause = append(whereClause, "last_id>$"+strconv.Itoa(len(whereClause)+1))
	}

	whereQueryStr := ""
	status := "active"
	queryType := c.Query("type")
	switch queryType {
	case "hot":
		orderByStr := "(((SELECT COUNT(posts.id) FROM posts WHERE posts.challenge_id=challenges.id AND (SELECT COUNT(likes.id) FROM likes WHERE likes.post_id=posts.id)) / (SELECT COUNT(posts.id) FROM posts WHERE posts.challenge_id=id)) * (SELECT COUNT(posts.id) FROM posts WHERE posts.challenge_id=id) / 100)"
		orderByStr = "(" + orderByStr + " + " + orderByStr + " * challenges.weight) DESC"
		whereQueryStr = "WHERE " + strings.Join(whereClause, " AND ") + " AND status='" + status + "' ORDER BY " + orderByStr
	case "fresh":
		orderByStr := "challenges.created_at DESC, (SELECT COUNT(posts.id) FROM posts WHERE posts.challenge_id=challenges.id) DESC"
		whereQueryStr = "WHERE " + strings.Join(whereClause, " AND ") + " AND status='" + status + "' ORDER BY " + orderByStr

	default:
		c.JSON(http.StatusBadRequest, &model.ErrResp{Error: "Invalid query strings", Fields: &[]string{"type"}})
		return
	}

	args := make([]interface{}, len(preparedValues))
	for i, val := range preparedValues {
		args[i] = val
	}

	challengeList, err := (&model.Challenge{}).Get(whereQueryStr, args)
	if err != nil {
		log.Printf("Fetch challenge error: %v", err)
		c.JSON(http.StatusInternalServerError, &model.ErrResp{Error: "Server error"})
		return
	}

	c.JSON(http.StatusOK, &challengeList)
}

//PostChallenge func handler creates a new challenge
func PostChallenge(c *gin.Context) {

}

//PutChallenge func handler updates a challenge. PS: It cant update 'name'.
func PutChallenge(c *gin.Context) {

}

//DeleteChallenge func handler deletes a challenge, if there is no post made.
func DeleteChallenge(c *gin.Context) {

}

//DeActivateChallenge func handler de-activates challenges which are not being used for a while.
func DeActivateChallenge(c *gin.Context) {

}

//ActivateChallenge func handler de-activates challenges which are not being used for a while.
func ActivateChallenge(c *gin.Context) {

}
