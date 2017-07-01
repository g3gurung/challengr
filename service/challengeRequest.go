package service

import (
	"log"
	"net/http"

	"strconv"

	"github.com/challengr/model"
	"github.com/gin-gonic/gin"
)

//GetChallengeRequest handler func fetches the challenge request to or from the user depending on the query parameters
func GetChallengeRequest(c *gin.Context) {
	userID, ok := c.MustGet("user_id").(int64)
	if !ok {
		log.Println("invalid token, user_id error")
		c.JSON(http.StatusBadRequest, &model.ErrResp{Error: "Invalid token."})
		return
	}
	userRole, ok := c.MustGet("role").(string)
	if !ok {
		log.Println("invalid token, role error")
		c.JSON(http.StatusBadRequest, &model.ErrResp{Error: "Invalid token"})
		return
	}

	paramUserID, err := strconv.ParseInt(c.Param("user_id"), 10, 64)
	if err != nil {
		log.Printf("Invalid path param user_id: %v", err)
		c.JSON(http.StatusBadRequest, &model.ErrResp{Error: "Invalid path param user_id"})
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

	queryType := c.Query("type")
	if queryType != "sent" && queryType != "recieved" {
		c.JSON(http.StatusBadRequest, &model.ErrResp{Error: "Invalid query string", Fields: &[]string{"type"}})
		return
	}

	if userID != paramUserID && userRole != constAdminRole {
		c.JSON(http.StatusMethodNotAllowed, &model.ErrResp{Error: "Not allowed"})
		return
	}

	challengeRequestList := []*model.ChallengeRequest{}
	if queryType == "sent" {
		challengeRequestList, err = (&model.ChallengeRequest{}).Get("WHERE from_id=$1 AND last_id>$2 ORDER BY DESC created_at LIMIT 20", paramUserID, lastID)
	} else {
		challengeRequestList, err = (&model.ChallengeRequest{}).Get("WHERE to_id=$1 AND last_id>$2 ORDER BY DESC created_at LIMIT 20", paramUserID, lastID)
	}
	if err != nil {
		log.Printf("Challenge request fetching err: %v", err)
		c.JSON(http.StatusInternalServerError, &model.ErrResp{Error: "Server error"})
		return
	}

	c.JSON(http.StatusOK, &challengeRequestList)
}

//PostChallengeRequest handler func sends challenge to the user who is in friendlist in fb
func PostChallengeRequest(c *gin.Context) {
	userID, ok := c.MustGet("user_id").(int64)
	if !ok {
		log.Println("invalid token, user_id error")
		c.JSON(http.StatusBadRequest, &model.ErrResp{Error: "Invalid token."})
		return
	}
	userRole, ok := c.MustGet("role").(string)
	if !ok {
		log.Println("invalid token, role error")
		c.JSON(http.StatusBadRequest, &model.ErrResp{Error: "Invalid token"})
		return
	}

	paramUserID, err := strconv.ParseInt(c.Param("user_id"), 10, 64)
	if err != nil {
		log.Printf("Invalid path param user_id: %v", err)
		c.JSON(http.StatusBadRequest, &model.ErrResp{Error: "Invalid path param user_id"})
		return
	}

	if userID != paramUserID && userRole != constAdminRole {
		c.JSON(http.StatusMethodNotAllowed, &model.ErrResp{Error: "Not allowed"})
		return
	}

	var challengeRequest model.ChallengeRequest
	if err := c.BindJSON(&challengeRequest); err != nil {
		log.Printf("challenge struct JSON bind error: %v", err)
		c.JSON(http.StatusBadRequest, &model.ErrResp{Error: err.Error()})
		return
	}

	if err := c.BindJSON(&challengeRequest.Payload); err != nil {
		log.Printf("challenge Payload JSON bind error: %v", err)
		c.JSON(http.StatusBadRequest, &model.ErrResp{Error: err.Error()})
		return
	}

	if errSlice := challengeRequest.ParseNotAllowedJSON(); len(errSlice) > 0 {
		log.Printf("challenge not allowed fields detected: %v", errSlice)
		c.JSON(http.StatusBadRequest, &model.ErrResp{Error: "Some fields are not allowed", Fields: &errSlice})
		return
	}

	if errSlice := challengeRequest.PostValidate(); len(errSlice) > 0 {
		log.Printf("challenge not allowed fields detected: %v", errSlice)
		c.JSON(http.StatusBadRequest, &model.ErrResp{Error: "Some fields are invalid", Fields: &errSlice})
		return
	}

	count, err := (&model.User{}).Count("WHERE id=$1", challengeRequest.ToID)
	if err != nil {
		log.Printf("User count error: %v", err)
		c.JSON(http.StatusInternalServerError, &model.ErrResp{Error: "Server error"})
		return
	}

	if count != 1 {
		c.JSON(http.StatusNotFound, &model.ErrResp{Error: "to_id, user not found"})
		return
	}

	count, err = (&model.Challenge{}).Count("WHERE id=$1", challengeRequest.ToID)
	if err != nil {
		log.Printf("Challenge count error: %v", err)
		c.JSON(http.StatusInternalServerError, &model.ErrResp{Error: "Server error"})
		return
	}

	if count != 1 {
		c.JSON(http.StatusNotFound, &model.ErrResp{Error: "challeng_id, challenge not found"})
		return
	}

	challengeRequest.FromID = paramUserID
	challengeRequest.Status = "open"
	if err = challengeRequest.Create(); err != nil {
		log.Printf("Challenge request create error: %v", err)
		c.JSON(http.StatusInternalServerError, &model.ErrResp{Error: "Server error"})
		return
	}

	c.JSON(http.StatusOK, &challengeRequest)
}

//PutChallengeRequest handler func updates the challenge request, used basically for updating the status
func PutChallengeRequest(c *gin.Context) {
	userID, ok := c.MustGet("user_id").(int64)
	if !ok {
		log.Println("invalid token, user_id error")
		c.JSON(http.StatusBadRequest, &model.ErrResp{Error: "Invalid token."})
		return
	}
	userRole, ok := c.MustGet("role").(string)
	if !ok {
		log.Println("invalid token, role error")
		c.JSON(http.StatusBadRequest, &model.ErrResp{Error: "Invalid token"})
		return
	}

	paramUserID, err := strconv.ParseInt(c.Param("user_id"), 10, 64)
	if err != nil {
		log.Printf("Invalid path param user_id: %v", err)
		c.JSON(http.StatusBadRequest, &model.ErrResp{Error: "Invalid path param user_id"})
		return
	}

	challengeStatus := c.Query("status")
	if challengeStatus != "rejected" && challengeStatus != "accepted" {
		c.JSON(http.StatusBadRequest, &model.ErrResp{Error: "Invalid query string", Fields: &[]string{"status"}})
		return
	}

	if userID != paramUserID && userRole != constAdminRole {
		c.JSON(http.StatusMethodNotAllowed, &model.ErrResp{Error: "Not allowed"})
		return
	}

	paramChallengeRequestID, err := strconv.ParseInt(c.Param("challenge_request_id"), 10, 64)
	if err != nil {
		log.Printf("Invalid path param user_id: %v", err)
		c.JSON(http.StatusBadRequest, &model.ErrResp{Error: "Invalid path param challenge_request_id"})
		return
	}

	challengeRequest := model.ChallengeRequest{ID: paramChallengeRequestID, ToID: paramUserID}
	status, err := challengeRequest.UpdateStatus(challengeStatus)
	if err != nil {
		log.Printf("Challenge request update status error: %v", err)
		c.JSON(status, &model.ErrResp{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, &model.SuccessResp{Status: http.StatusOK, Message: "challenge request successfully udated", Response: map[string]interface{}{"challenge_id": challengeRequest.ChallengeID}})
}

//DeleteChallengeRequest handler func updates the challenge request, used basically for updating the status
func DeleteChallengeRequest(c *gin.Context) {
	userID, ok := c.MustGet("user_id").(int64)
	if !ok {
		log.Println("invalid token, user_id error")
		c.JSON(http.StatusBadRequest, &model.ErrResp{Error: "Invalid token."})
		return
	}
	userRole, ok := c.MustGet("role").(string)
	if !ok {
		log.Println("invalid token, role error")
		c.JSON(http.StatusBadRequest, &model.ErrResp{Error: "Invalid token"})
		return
	}

	paramUserID, err := strconv.ParseInt(c.Param("user_id"), 10, 64)
	if err != nil {
		log.Printf("Invalid path param user_id: %v", err)
		c.JSON(http.StatusBadRequest, &model.ErrResp{Error: "Invalid path param user_id"})
		return
	}

	if userID != paramUserID && userRole != constAdminRole {
		c.JSON(http.StatusMethodNotAllowed, &model.ErrResp{Error: "Not allowed"})
		return
	}

	paramChallengeRequestID, err := strconv.ParseInt(c.Param("challenge_request_id"), 10, 64)
	if err != nil {
		log.Printf("Invalid path param user_id: %v", err)
		c.JSON(http.StatusBadRequest, &model.ErrResp{Error: "Invalid path param challenge_request_id"})
		return
	}

	challengeRequest := model.ChallengeRequest{ID: paramChallengeRequestID, FromID: paramUserID}
	status, err := challengeRequest.Delete()
	if err != nil {
		log.Printf("Challenge request update status error: %v", err)
		c.JSON(status, &model.ErrResp{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, &model.SuccessResp{Status: http.StatusOK, Message: "challenge request successfully deleted"})
}
