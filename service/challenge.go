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
		whereQueryStr = "WHERE " + strings.Join(whereClause, " AND ") + " AND status='" + status + "' AND deleted_at IS NULL ORDER BY " + orderByStr + " LIMIT 20;"
	case "fresh":
		orderByStr := "challenges.created_at DESC, (SELECT COUNT(posts.id) FROM posts WHERE posts.challenge_id=challenges.id) DESC"
		whereQueryStr = "WHERE " + strings.Join(whereClause, " AND ") + " AND status='" + status + "' AND deleted_at IS NULL ORDER BY " + orderByStr + " LIMIT 20;"

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
	weight, ok := c.MustGet("weight").(float32)
	if !ok {
		log.Println("invalid token, wight error")
		c.JSON(http.StatusBadRequest, &model.ErrResp{Error: "Invalid token"})
		return
	}

	var challenge model.Challenge
	if err := c.BindJSON(&challenge); err != nil {
		log.Printf("challenge struct JSON bind error: %v", err)
		c.JSON(http.StatusBadRequest, &model.ErrResp{Error: err.Error()})
		return
	}

	if err := c.BindJSON(&challenge.Payload); err != nil {
		log.Printf("challenge Payload JSON bind error: %v", err)
		c.JSON(http.StatusBadRequest, &model.ErrResp{Error: err.Error()})
		return
	}

	if errSlice := challenge.ParseNotAllowedJSON(); len(errSlice) > 0 {
		log.Printf("challenge not allowed fields detected: %v", errSlice)
		c.JSON(http.StatusBadRequest, &model.ErrResp{Error: "Some fields are not allowed", Fields: &errSlice})
		return
	}

	if userRole != constAdminRole {
		usersList, err := (&model.User{}).Get("WHERE id=$1", userID)
		if err != nil {
			log.Printf("User fetch error: %v", err)
			c.JSON(http.StatusInternalServerError, &model.ErrResp{Error: "Server error"})
			return
		}
		if usersList[0].LevelID < 5 {
			log.Printf("User levelID: %v, must be 5 or bigger", usersList[0].LevelID)
			c.JSON(http.StatusMethodNotAllowed, &model.ErrResp{Error: "Not allowed. Need level 5 or more."})
			return
		}
	}

	if errSlice := challenge.PostValidate(); len(errSlice) > 0 {
		log.Printf("challenge post validate err: %v", errSlice)
		c.JSON(http.StatusBadRequest, &model.ErrResp{Error: "Invalid fields detected", Fields: &errSlice})
		return
	}

	challenge.UserID = userID
	challenge.Weight = &weight

	if err := challenge.Create(); err != nil {
		log.Printf("challenge create err: %v", err)
		c.JSON(http.StatusInternalServerError, &model.ErrResp{Error: "Server error"})
		return
	}

	c.JSON(http.StatusOK, &challenge)
}

//PutChallenge func handler updates a challenge. PS: It cant update 'name'.
func PutChallenge(c *gin.Context) {
	userID, ok := c.MustGet("user_id").(int64)
	if !ok {
		log.Println("invalid token, user_id error")
		c.JSON(http.StatusBadRequest, &model.ErrResp{Error: "Invalid token."})
		return
	}
	paramChallengeID := c.Param("challenge_id")
	challengeID, err := strconv.ParseInt(paramChallengeID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, &model.ErrResp{Error: "Invalid path params", Fields: &[]string{"challenge_id"}})
		return
	}

	challenge := model.Challenge{ID: challengeID}
	if err := c.BindJSON(&challenge); err != nil {
		log.Printf("challenge struct JSON bind error: %v", err)
		c.JSON(http.StatusBadRequest, &model.ErrResp{Error: err.Error()})
		return
	}

	if err := c.BindJSON(&challenge.Payload); err != nil {
		log.Printf("challenge Payload JSON bind error: %v", err)
		c.JSON(http.StatusBadRequest, &model.ErrResp{Error: err.Error()})
		return
	}

	if errSlice := challenge.ParseNotAllowedJSON(); len(errSlice) > 0 {
		log.Printf("challenge not allowed fields detected: %v", errSlice)
		c.JSON(http.StatusBadRequest, &model.ErrResp{Error: "Some fields are not allowed", Fields: &errSlice})
		return
	}

	notAllowedFields := []string{}
	if challenge.Status != "" {
		notAllowedFields = append(notAllowedFields, "status")
	}
	if challenge.Name != "" {
		notAllowedFields = append(notAllowedFields, "name")
	}
	if challenge.LikesNeededPerPost != 0 {
		notAllowedFields = append(notAllowedFields, "likes_needed_per_post")
	}
	if challenge.Weight != nil {
		notAllowedFields = append(notAllowedFields, "weight")
	}

	if len(notAllowedFields) > 0 {
		c.JSON(http.StatusBadRequest, &model.ErrResp{Error: "Some fields are not allowed", Fields: &notAllowedFields})
		return
	}

	if challenge.Description == nil && challenge.Location == nil {
		log.Printf("challenge Payload invalid json: %v", challenge.Payload)
		c.JSON(http.StatusBadRequest, &model.ErrResp{Error: "No valid payload detected"})
		return
	}

	challenge.UserID = userID

	if errStatus, err := challenge.Update(); err != nil {
		log.Printf("challenge update error: %v", err)
		c.JSON(errStatus, &model.ErrResp{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, &model.SuccessResp{Message: "Challenge successfuly updated", Status: http.StatusOK})
}

//DeleteChallenge func handler deletes a challenge, if there is no post made.
func DeleteChallenge(c *gin.Context) {
	userID, ok := c.MustGet("user_id").(int64)
	if !ok {
		log.Println("invalid token, user_id error")
		c.JSON(http.StatusBadRequest, &model.ErrResp{Error: "Invalid token."})
		return
	}
	role, ok := c.MustGet("role").(string)
	if !ok {
		log.Println("invalid token, role error")
		c.JSON(http.StatusBadRequest, &model.ErrResp{Error: "Invalid token."})
		return
	}

	paramChallengeID := c.Param("challenge_id")
	challengeID, err := strconv.ParseInt(paramChallengeID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, &model.ErrResp{Error: "Invalid path params", Fields: &[]string{"challenge_id"}})
		return
	}

	errStatus := 0
	switch role {
	case constAdminRole:
		errStatus, err = (&model.Challenge{ID: challengeID}).Delete()
	case constUserRole:
		errStatus, err = (&model.Challenge{ID: challengeID, UserID: userID}).AdminDelete()
	default:
		c.JSON(http.StatusForbidden, &model.ErrResp{Error: "Check token"})
		return
	}

	if err != nil {
		c.JSON(errStatus, err.Error())
		return
	}

	c.JSON(http.StatusOK, &model.SuccessResp{Message: "Challenge successfuly deleted", Status: http.StatusOK})
}

func activeDeactiveChallenge(val string, c *gin.Context) {
	userID, ok := c.MustGet("user_id").(int64)
	if !ok {
		log.Println("invalid token, user_id error")
		c.JSON(http.StatusBadRequest, &model.ErrResp{Error: "Invalid token."})
		return
	}
	role, ok := c.MustGet("role").(string)
	if !ok {
		log.Println("invalid token, role error")
		c.JSON(http.StatusBadRequest, &model.ErrResp{Error: "Invalid token."})
		return
	}
	paramChallengeID := c.Param("challenge_id")
	challengeID, err := strconv.ParseInt(paramChallengeID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, &model.ErrResp{Error: "Invalid path params", Fields: &[]string{"challenge_id"}})
		return
	}

	errStatus := 0
	switch role {
	case constAdminRole:
		errStatus, err = (&model.Challenge{ID: challengeID, Status: val}).AdminUpdate()
	case constUserRole:
		errStatus, err = (&model.Challenge{ID: challengeID, UserID: userID, Status: val}).Update()
	default:
		c.JSON(http.StatusForbidden, &model.ErrResp{Error: "Check token"})
		return
	}

	if err != nil {
		c.JSON(errStatus, err.Error())
		return
	}

	c.JSON(http.StatusOK, &model.SuccessResp{Message: "Challenge successfuly updated", Status: http.StatusOK})
}

//DeActivateChallenge func handler de-activates challenges which are not being used for a while.
func DeActivateChallenge(c *gin.Context) {
	activeDeactiveChallenge("inactive", c)
}

//ActivateChallenge func handler de-activates challenges which are not being used for a while.
func ActivateChallenge(c *gin.Context) {
	activeDeactiveChallenge("active", c)
}
