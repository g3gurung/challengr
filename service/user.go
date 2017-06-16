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
		userlist []*model.User
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
		userList, err := (&model.User{}).Get("WHERE last_id>$1 AND deleted_at IS NULL AND " + geomQ + " ORDER BY users.level_id DESC LIMIT 100")
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
		for i, v := range IDs {
			index = index + 1
			args[index] = v
			idBuilder = append(idBuilder, "users.id == $"+strconv.Itoa(index))
		}
		for i, v := range fbIDs {
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
		userList, err := (&model.User{}).Get("WHERE deleted_at IS NULL AND "+q+" ORDER BY users.level_id DESC", args)
	}

	if err != nil {
		log.Printf("User fetching error %v", err)
		c.JSON(http.StatusInternalServerError, &model.ErrResp{Error: "Server error"})
		return
	}

	c.JSON(http.StatusOK, &userlist)
}

//PutUser func is a handler for updateing a user info.
func PutUser(c *gin.Context) {

}
