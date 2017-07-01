package service

import (
	"log"
	"net/http"
	"strconv"

	"github.com/challengr/model"
	"github.com/gin-gonic/gin"
)

//PostPost func is a handler for creating a new post
func PostPost(c *gin.Context) {
	paramChallengeID := c.Param("challenge_id")
	challengeID, err := strconv.ParseInt(paramChallengeID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, &model.ErrResp{Error: "Invalid path params", Fields: &[]string{"challenge_id"}})
		return
	}

	userID, ok := c.MustGet("user_id").(int64)
	if !ok {
		log.Println("token parsing not ok")
		c.JSON(http.StatusBadRequest, &model.ErrResp{Error: "Invalid token"})
		return
	}

	var post model.Post
	if err := c.BindJSON(&post); err != nil {
		log.Printf("post struct JSON bind error: %v", err)
		c.JSON(http.StatusBadRequest, &model.ErrResp{Error: err.Error()})
		return
	}

	if err := c.BindJSON(&post.Payload); err != nil {
		log.Printf("post Payload JSON bind error: %v", err)
		c.JSON(http.StatusBadRequest, &model.ErrResp{Error: err.Error()})
		return
	}

	if errSlice := post.ParseNotAllowedJSON(); len(errSlice) > 0 {
		log.Printf("post not allowed fields detected: %v", errSlice)
		c.JSON(http.StatusBadRequest, &model.ErrResp{Error: "Some fields are not allowed", Fields: &errSlice})
		return
	}

	if errSlice := post.PostValidate(); len(errSlice) > 0 {
		log.Printf("post validate err: %v", errSlice)
		c.JSON(http.StatusBadRequest, &model.ErrResp{Error: "Invalid fields detected", Fields: &errSlice})
		return
	}

	post.UserID = userID
	post.ChallengeID = challengeID

	if err := post.Create(); err != nil {
		c.JSON(http.StatusInternalServerError, &model.ErrResp{Error: "Server error"})
		return
	}

	c.JSON(http.StatusOK, &post)
}

//GetPost func is a handler for fetching posts
func GetPost(c *gin.Context) {
	paramChallengeID := c.Param("challenge_id")
	challengeID, err := strconv.ParseInt(paramChallengeID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, &model.ErrResp{Error: "Invalid path params", Fields: &[]string{"challenge_id"}})
		return
	}

	queryLastID := c.Query("last_id")
	var lastID int64
	lastID = 0
	if queryLastID != "" {
		lastID, err = strconv.ParseInt(queryLastID, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, &model.ErrResp{Error: "Invalid path params", Fields: &[]string{"last_id"}})
			return
		}
	}

	postList, err := (&model.Post{}).Get("WHERE challenge_id=$1 AND deleted_at IS NULL AND id > $2 ORDER BY created_at DESC LIMIT 30", challengeID, lastID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &model.ErrResp{Error: "Server error"})
		return
	}

	c.JSON(http.StatusOK, &postList)
}

//DeletePost func is a handler for deleting a post
func DeletePost(c *gin.Context) {
	paramChallengeID := c.Param("challenge_id")
	challengeID, err := strconv.ParseInt(paramChallengeID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, &model.ErrResp{Error: "Invalid path params", Fields: &[]string{"challenge_id"}})
		return
	}

	paramPostID := c.Param("challenge_id")
	postID, err := strconv.ParseInt(paramPostID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, &model.ErrResp{Error: "Invalid path params", Fields: &[]string{"post_id"}})
		return
	}

	userID, ok := c.MustGet("user_id").(int64)
	if !ok {
		log.Println("token parsing not ok")
		c.JSON(http.StatusBadRequest, &model.ErrResp{Error: "Invalid token"})
		return
	}

	userRole, ok := c.MustGet("user_role").(string)
	if !ok {
		log.Println("token parsing not ok")
		c.JSON(http.StatusBadRequest, &model.ErrResp{Error: "Invalid token"})
		return
	}

	switch userRole {
	case constAdminRole:
		status, err := (&model.Post{ID: postID, ChallengeID: challengeID}).Delete()
		if err != nil {
			c.JSON(status, &model.ErrResp{Error: err.Error()})
			return
		}
	case constUserRole:
		status, err := (&model.Post{ID: postID, ChallengeID: challengeID, UserID: userID}).Delete()
		if err != nil {
			c.JSON(status, &model.ErrResp{Error: err.Error()})
			return
		}
	default:
		log.Printf("token parsing invalid role -> %v", userRole)
		c.JSON(http.StatusForbidden, &model.ErrResp{Error: "Invalid token"})
		return
	}

	c.JSON(http.StatusOK, &model.SuccessResp{Message: "Post successfully deleted", Status: http.StatusOK})
}

//FlagPost func is a handler for flagging a post
func FlagPost(c *gin.Context) {
	paramChallengeID := c.Param("challenge_id")
	challengeID, err := strconv.ParseInt(paramChallengeID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, &model.ErrResp{Error: "Invalid path params", Fields: &[]string{"challenge_id"}})
		return
	}

	paramPostID := c.Param("challenge_id")
	postID, err := strconv.ParseInt(paramPostID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, &model.ErrResp{Error: "Invalid path params", Fields: &[]string{"post_id"}})
		return
	}

	userID, ok := c.MustGet("user_id").(int64)
	if !ok {
		log.Println("token parsing not ok")
		c.JSON(http.StatusBadRequest, &model.ErrResp{Error: "Invalid token"})
		return
	}

	status, err := (&model.Post{ID: postID, ChallengeID: challengeID}).Flag(userID)
	if err != nil {
		c.JSON(status, &model.ErrResp{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, &model.SuccessResp{Message: "Post successfully flagged", Status: http.StatusOK})
}

//UnFlagPost func is a handler or unflagging a post
func UnFlagPost(c *gin.Context) {
	paramChallengeID := c.Param("challenge_id")
	challengeID, err := strconv.ParseInt(paramChallengeID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, &model.ErrResp{Error: "Invalid path params", Fields: &[]string{"challenge_id"}})
		return
	}

	paramPostID := c.Param("challenge_id")
	postID, err := strconv.ParseInt(paramPostID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, &model.ErrResp{Error: "Invalid path params", Fields: &[]string{"post_id"}})
		return
	}

	userID, ok := c.MustGet("user_id").(int64)
	if !ok {
		log.Println("token parsing not ok")
		c.JSON(http.StatusBadRequest, &model.ErrResp{Error: "Invalid token"})
		return
	}

	status, err := (&model.Post{ID: postID, ChallengeID: challengeID}).UnFlag(userID)
	if err != nil {
		c.JSON(status, &model.ErrResp{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, &model.SuccessResp{Message: "Post successfully unflagged", Status: http.StatusOK})
}

//LikePost func is a handler for licking a post
func LikePost(c *gin.Context) {
	paramChallengeID := c.Param("challenge_id")
	challengeID, err := strconv.ParseInt(paramChallengeID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, &model.ErrResp{Error: "Invalid path params", Fields: &[]string{"challenge_id"}})
		return
	}

	paramPostID := c.Param("challenge_id")
	postID, err := strconv.ParseInt(paramPostID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, &model.ErrResp{Error: "Invalid path params", Fields: &[]string{"post_id"}})
		return
	}

	userID, ok := c.MustGet("user_id").(int64)
	if !ok {
		log.Println("token parsing not ok")
		c.JSON(http.StatusBadRequest, &model.ErrResp{Error: "Invalid token"})
		return
	}

	status, err := (&model.Post{ID: postID, ChallengeID: challengeID}).Like(userID)
	if err != nil {
		c.JSON(status, &model.ErrResp{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, &model.SuccessResp{Message: "Post successfully liked", Status: http.StatusOK})
}
