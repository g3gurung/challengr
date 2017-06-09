package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

//Challenge struct is a model/schema for a challenge table
type Challenge struct {
	ID          int64      `json:"id" sql:"id"`
	UserID      int64      `json:"user_id" sql:"user_id"`
	Name        string     `json:"name" sql:"name"`
	Description *string    `json:"description" sql:"description"`
	Status      string     `json:"status" sql:"status"`
	Weight      *int       `json:"weight" sql:"weight"`
	CreatedAt   *time.Time `json:"created_at" sql:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at" sql:"updated_at"`

	TotalPost int64     `json:"total_post" sql:"-"`
	Location  *geometry `json:"geo_coords" sql:"-"`

	Payload map[string]interface{} `json:"-"`
}

//ParseNotAllowedJSON unmarshalls JSON payload to struct payload and fields. Plus parses the JSON payload.
func (c *Challenge) ParseNotAllowedJSON() []string {
	errSlice := []string{}

	delete(c.Payload, "name")
	delete(c.Payload, "description")
	delete(c.Payload, "geo_coords")

	for key := range c.Payload {
		errSlice = append(errSlice, key)
	}

	return errSlice
}

//PostValidate func validates incoming post payload fields
func (c *Challenge) PostValidate() []string {
	errSlice := []string{}

	if c.Name == "" {
		errSlice = append(errSlice, "name")
	}

	return errSlice
}

//Create func inserts a new challenge in the db
func (c *Challenge) Create() error {
	now := time.Now()
	c.CreatedAt = &now

	geomStr, err := json.Marshal(c.Location)
	if err != nil {
		log.Printf("Bad location value err: %v\n", err)
		return err
	}

	geometryValue := "ST_GeomFromGeoJSON('" + string(geomStr) + "')"

	stmt, err := db.Prepare("INSERT INTO cahllenge(user_id, name, description, status, weight, geometry, created_at) VALUES($1,$2,$3,$4,$5" + geometryValue + ",$6);")
	if err != nil {
		log.Printf("create prepare statement error: %v", err)
		return err
	}

	res, err := stmt.Exec(c.UserID, c.Name, c.Description, c.Status, c.Weight, c.CreatedAt)
	if err != nil {
		log.Printf("exec statement error: %v", err)
		return err
	}

	c.ID, err = res.LastInsertId()
	if err != nil {
		log.Printf("last insert id error: %v", err)
		return err
	}

	log.Printf("onesignal successfully created with id %v", c.ID)

	affected, err := res.RowsAffected()
	if err != nil {
		log.Printf("rows effected error: %v", err)
		return err
	}
	if affected == 0 {
		log.Printf("rows effected -> %v", affected)
		return errors.New("Server error")
	}

	return nil
}

//Update func updates a challenge in the db
func (c *Challenge) Update() error {
	sets := []string{}
	values := make(map[int]interface{})
	index := 0

	if c.Description != nil {
		values[index] = *c.Description
		index = index + 1
		sets = append(sets, "description=$"+strconv.Itoa(index))
	}

	if c.Status != "" {
		values[index] = c.Status
		index = index + 1
		sets = append(sets, "status=$"+strconv.Itoa(index))
	}

	if c.Weight != nil {
		values[index] = *c.Weight
		index = index + 1
		sets = append(sets, "weight=$"+strconv.Itoa(index))
	}

	if c.UpdatedAt != nil {
		values[index] = c.UpdatedAt
		index = index + 1
		sets = append(sets, "updated_at=$"+strconv.Itoa(index))
	}

	stmt, err := db.Prepare("UPDATE challenges SET " + strings.Join(sets, ", ") + " WHERE id=" + fmt.Sprintf("%v", c.ID) + ";")
	if err != nil {
		log.Printf("UPDATE challegne prepare statement error: %v", err)
		return err
	}

	argsValues := make([]interface{}, len(values))
	for k, v := range values {
		argsValues[k] = v
	}

	res, err := stmt.Exec(argsValues...)
	if err != nil {
		log.Printf("exec statement error: %v", err)
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		log.Printf("rows effected error: %v", err)
		return err
	}
	if affected == 0 {
		log.Printf("rows effected -> %v", affected)
		return errors.New("Server error")
	}

	return nil
}

//Get func fetches the challenges from the db based on the query
func (c *Challenge) Get(whereClause string, args ...interface{}) ([]*Challenge, error) {
	challengeList := []*Challenge{}

	rows, err := db.Query("SELECT id, name, description, ST_AsGeoJSON(geometry) AS location, (SELECT COUNT(id) FROM posts WHERE posts.challenge_id=challenges.id) AS total_post, created_at, updated_at FROM flags WHERE post_id=posts.id) as flags FROM posts "+whereClause+" ORDER BY created_at DESC;", args...)
	if err != nil {
		log.Printf("Get users: sql error %v", err)
		return nil, err
	}
	for rows.Next() {
		challenge := Challenge{}
		geomStr := ""
		if err = rows.Scan(&c.ID, &c.Name, &c.Description, &geomStr, &c.TotalPost, &c.CreatedAt, &c.UpdatedAt); err != nil {
			log.Printf("scanning row to struct error: %v", err)
			return nil, err
		}

		if err = json.Unmarshal([]byte(geomStr), &c.Location); err != nil {
			log.Printf("Unmarshaling of location subquery error: %v", err)
			return nil, err
		}

		challengeList = append(challengeList, &challenge)
	}

	return challengeList, nil
}

//Count func counts the total challenges in db
func (c *Challenge) Count(whereClause string, args ...interface{}) (int64, error) {
	var count int64

	if err := db.QueryRow("SELECT COUNT(id) FROM challenges "+whereClause+";", args...).Scan(&count); err != nil {
		log.Printf("Count onesignal: sql error %v", err)
		return count, err
	}

	return count, nil
}

//Delete func deletes the post record of the user. Delete meaning it doesnt purge it. Just hides it.
func (c *Challenge) Delete() error {
	count, err := c.Count("WHERE id=$1", c.ID)
	if err != nil {
		log.Printf("Post delete: error on fetching Post record count: %v", err)
		return err
	}

	if count == 0 {
		err = fmt.Errorf("Post account not found-> id %v, total found %v", c.ID, count)
		log.Printf("%v", err)
		return err
	} else if count == 1 {
		stmt, err := db.Prepare("UPDATE challenges SET deleted_at=$1 WHERE id=$2;")
		if err != nil {
			log.Printf("create prepare statement error: %v", err)
			return err
		}

		res, err := stmt.Exec(time.Now(), c.ID)
		if err != nil {
			log.Printf("exec statement error: %v", err)
			return err
		}

		affected, err := res.RowsAffected()
		if err != nil {
			log.Printf("rows effected error: %v", err)
			return err
		}
		if affected == 0 {
			log.Printf("rows effected -> %v", affected)
			return errors.New("Server error")
		}
	} else {
		err = fmt.Errorf("multiple posts found-> id %v, total found %v", c.ID, count)
		log.Printf("%v", err)
		return err
	}

	return nil
}
