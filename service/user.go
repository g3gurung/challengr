package service

import (
	"net/http"
	"strconv"

	"strings"

	"regexp"

	"log"

	"github.com/challengr/model"
	"github.com/gin-gonic/gin"
)

//GetUser func is a handler for fetching user list. The level of detials depends on the role of the user.
func GetUser(c *gin.Context) {
	var (
		err      error
		userList []*model.User
	)
	queryType := c.Query("type")
	switch queryType {
	case "ranking":
		queryLastID := c.Query("last_id")
		var lastID int64
		if queryLastID != "" {
			lastID, err = strconv.ParseInt(queryLastID, 10, 64)
			if err != nil {
				c.JSON(http.StatusBadRequest, &model.ErrResp{Error: "Invalid query string", Fields: &[]string{"last_id"}})
				return
			}
		} else {
			lastID = 0
		}

		queryRadius := strings.TrimSpace(c.Query("radius"))
		radius, err := strconv.Atoi(queryRadius)
		if err != nil {
			c.JSON(http.StatusBadRequest, &model.ErrResp{Error: "Invalid query string", Fields: &[]string{"radius"}})
			return
		}
		if radius <= 0 {
			c.JSON(http.StatusBadRequest, &model.ErrResp{Error: "Invalid query string", Fields: &[]string{"radius"}})
			return
		}

		queryLong := strings.TrimSpace(c.Query("longitude"))
		if r := regexp.MustCompile("^[-+]?(180(\\.0+)?|((1[0-7]\\d)|([1-9]?\\d))(\\.\\d+)?)$"); !r.MatchString(queryLong) {
			c.JSON(http.StatusBadRequest, &model.ErrResp{Error: "Invalid query string", Fields: &[]string{"longitude"}})
			return
		}

		queryLat := strings.TrimSpace(c.Query("latitude"))
		if r := regexp.MustCompile("^[-+]?([1-8]?\\d(\\.\\d+)?|90(\\.0+)?)$"); !r.MatchString(queryLat) {
			c.JSON(http.StatusBadRequest, &model.ErrResp{Error: "Invalid query string", Fields: &[]string{"latitude"}})
			return
		}
		geomQ := string("ST_Distance_Sphere(geometry, ST_MakePoint(" + queryLong + "," + queryLat + ")) <= " + queryRadius)
		userList, err = (&model.User{}).Get("WHERE last_id>$1 AND deleted_at IS NULL AND "+geomQ+" ORDER BY users.level_id DESC LIMIT 100", lastID)
	default:
		IDs := []int64{}
		queryIDs := strings.Split(c.Query("ids"), ",")
		for i, v := range queryIDs {
			queryIDs[i] = strings.TrimSpace(v)
			if queryIDs[i] != "" {
				ID, err := strconv.ParseInt(queryIDs[i], 10, 64)
				if err != nil {
					c.JSON(http.StatusBadRequest, &model.ErrResp{Error: "Invalid query string", Fields: &[]string{"ids"}})
					return
				}
				IDs = append(IDs, ID)
			} else if len(queryIDs) == 1 {
				queryIDs = []string{}
			} else {
				c.JSON(http.StatusBadRequest, &model.ErrResp{Error: "Invalid query string", Fields: &[]string{"ids"}})
				return
			}
		}

		fbIDs := strings.Split(c.Query("fb_ids"), ",")
		for i, v := range queryIDs {
			queryIDs[i] = strings.TrimSpace(v)
			if queryIDs[i] == "" {
				if len(queryIDs) == 1 {
					fbIDs = []string{}
				} else {
					c.JSON(http.StatusBadRequest, &model.ErrResp{Error: "Invalid query string", Fields: &[]string{"fb_ids"}})
					return
				}
			}
		}

		if len(IDs) == 0 && len(fbIDs) == 0 {
			c.JSON(http.StatusBadRequest, &model.ErrResp{Error: "Invalid query string"})
			return
		}

		args := make([]interface{}, len(IDs)+len(fbIDs))
		idBuilder := []string{}
		fbIDQBuilder := []string{}
		index := 0
		for _, v := range IDs {
			index = index + 1
			args[index] = v
			idBuilder = append(idBuilder, "users.id == $"+strconv.Itoa(index))
		}
		for _, v := range fbIDs {
			index = index + 1
			args[index] = v
			fbIDQBuilder = append(fbIDQBuilder, "users.facebook_user_id = $"+strconv.Itoa(index))
		}
		idQ := strings.Join(idBuilder, " OR ")
		fbIDQ := strings.Join(fbIDQBuilder, " OR ")

		q := ""
		if idQ != "" {
			q = "(" + idQ + ")"
		}
		if fbIDQ != "" {
			if q != "" {
				q = q + " AND " + "(" + fbIDQ + ")"
			} else {
				q = "(" + fbIDQ + ")"
			}
		}
		userList, err = (&model.User{}).Get("WHERE deleted_at IS NULL AND "+q+" ORDER BY users.level_id DESC", args)
	}

	if err != nil {
		log.Printf("User fetching error %v", err)
		c.JSON(http.StatusInternalServerError, &model.ErrResp{Error: "Server error"})
		return
	}

	c.JSON(http.StatusOK, &userList)
}

//UpdateUserWeight func is a handler for updateing a user info.
func UpdateUserWeight(c *gin.Context) {
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

	var user model.User
	if err := c.BindJSON(&user); err != nil {
		log.Printf("user struct JSON bind error: %v", err)
		c.JSON(http.StatusBadRequest, &model.ErrResp{Error: err.Error()})
		return
	}

	if err := c.BindJSON(&user.Payload); err != nil {
		log.Printf("user Payload JSON bind error: %v", err)
		c.JSON(http.StatusBadRequest, &model.ErrResp{Error: err.Error()})
		return
	}

	if errSlice := user.ParseNotAllowedJSON(); len(errSlice) > 0 {
		log.Printf("user not allowed fields detected: %v", errSlice)
		c.JSON(http.StatusBadRequest, &model.ErrResp{Error: "Some fields are not allowed", Fields: &errSlice})
		return
	}

	if user.Weight == nil {
		c.JSON(http.StatusBadRequest, &model.ErrResp{Error: "Invalid payload", Fields: &[]string{"weight"}})
		return
	}

	if err := (&model.User{ID: paramUserID, Weight: user.Weight}).Update(); err != nil {
		log.Printf("Error user weight update: %v", err)
		c.JSON(http.StatusInternalServerError, &model.ErrResp{Error: "Server error"})
		return
	}

	c.JSON(http.StatusOK, &model.SuccessResp{Message: "User weight succesfully updated", Status: http.StatusOK})
}

//UpdateUserLevel func handler updates the user level
func UpdateUserLevel(c *gin.Context) {
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

	var user model.User
	if err := c.BindJSON(&user); err != nil {
		log.Printf("user struct JSON bind error: %v", err)
		c.JSON(http.StatusBadRequest, &model.ErrResp{Error: err.Error()})
		return
	}

	if err := c.BindJSON(&user.Payload); err != nil {
		log.Printf("user Payload JSON bind error: %v", err)
		c.JSON(http.StatusBadRequest, &model.ErrResp{Error: err.Error()})
		return
	}

	if errSlice := user.ParseNotAllowedJSON(); len(errSlice) > 0 {
		log.Printf("user not allowed fields detected: %v", errSlice)
		c.JSON(http.StatusBadRequest, &model.ErrResp{Error: "Some fields are not allowed", Fields: &errSlice})
		return
	}

	if user.LevelID == 0 {
		c.JSON(http.StatusBadRequest, &model.ErrResp{Error: "Invalid payload", Fields: &[]string{"level_id"}})
		return
	}

	if err := (&model.User{ID: paramUserID, LevelID: user.LevelID}).Update(); err != nil {
		log.Printf("Error user weight update: %v", err)
		c.JSON(http.StatusInternalServerError, &model.ErrResp{Error: "Server error"})
		return
	}

	c.JSON(http.StatusOK, &model.SuccessResp{Message: "User level succesfully updated", Status: http.StatusOK})
}
