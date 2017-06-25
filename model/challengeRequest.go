package model

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"
)

//ChallengeRequest struct is a schema or a model for a challenge_requests table
type ChallengeRequest struct {
	ID          int64      `json:"id" sql:"id"`
	FromID      int64      `json:"from_id,omitempty" sql:"from_id"`
	From        *User      `json:"from" sql:"-"`
	ToID        int64      `json:"to_id,omitempty" sql:"to_id"`
	To          *User      `json:"to" sql:"-"`
	ChallengeID int64      `json:"challenge_id,omitempty" sql:"challenge_id"`
	Challenge   *Challenge `json:"challenge" sql:"-"`
	Message     string     `json:"message" sql:"message"`
	Status      string     `json:"status" sql:"status"` //rejected, accepted and completed
	CreatedAt   time.Time  `json:"created_at" sql:"created_at"`

	Payload map[string]interface{} `json:"-"`
}

//ParseNotAllowedJSON unmarshalls JSON payload to struct payload and fields. Plus parses the JSON payload.
func (c *ChallengeRequest) ParseNotAllowedJSON() []string {
	errSlice := []string{}

	delete(c.Payload, "from_id")
	delete(c.Payload, "to_id")
	delete(c.Payload, "message")

	for key := range c.Payload {
		errSlice = append(errSlice, key)
	}

	return errSlice
}

//PostValidate func validates incoming post payload fields
func (c *ChallengeRequest) PostValidate() []string {
	errSlice := []string{}

	if c.FromID <= 0 {
		errSlice = append(errSlice, "from_id")
	}

	if c.ToID <= 0 {
		errSlice = append(errSlice, "to_id")
	}

	return errSlice
}

//Count func counts the users from db
func (c *ChallengeRequest) Count(whereClause string, args ...interface{}) (int64, error) {
	var count int64

	if err := db.QueryRow("SELECT COUNT(id) FROM challenge_requests "+whereClause+";", args...).Scan(&count); err != nil {
		log.Printf("Count levels: sql error %v", err)
		return count, err
	}

	return count, nil
}

//Get func fetches the vanity items from the db
func (c *ChallengeRequest) Get(whereClause string, args ...interface{}) ([]*ChallengeRequest, error) {
	ChallengeRequestList := []*ChallengeRequest{}
	rows, err := db.Query(`SELECT id, (SELECT row_to_json(users) FROM users WHERE users.id=challenge_requests.to_id) as to, 
	(SELECT row_to_json(users) FROM users WHERE users.id=challenge_requests.from_id) as from, 
	(SELECT row_to_json(challenges) FROM challenges WHERE challenges.id=challenge_requests.from_challenge_id) as challenge, message, created_at FROM bought_items `+whereClause+" ORDER BY created_at DESC;", args...)
	if err != nil {
		log.Printf("Get users: sql error %v", err)
		return nil, err
	}
	for rows.Next() {
		ChallengeRequest := ChallengeRequest{}
		toStr := ""
		fromStr := ""
		challengeStr := ""
		if err = rows.Scan(&ChallengeRequest.ID, &toStr, &fromStr, &challengeStr, &ChallengeRequest.Message, &ChallengeRequest.CreatedAt); err != nil {
			log.Printf("scanning row to struct error: %v", err)
			return nil, err
		}

		if err = json.Unmarshal([]byte(toStr), &ChallengeRequest.To); err != nil {
			log.Printf("Unmarshaling of subquery error: %v", err)
			return nil, err
		}

		if err = json.Unmarshal([]byte(fromStr), &ChallengeRequest.From); err != nil {
			log.Printf("Unmarshaling of subquery error: %v", err)
			return nil, err
		}

		if err = json.Unmarshal([]byte(challengeStr), &ChallengeRequest.Challenge); err != nil {
			log.Printf("Unmarshaling of subquery error: %v", err)
			return nil, err
		}

		ChallengeRequestList = append(ChallengeRequestList, &ChallengeRequest)
	}
	return ChallengeRequestList, nil
}

//Create func adds an item to the table
func (c *ChallengeRequest) Create() error {
	now := time.Now()
	c.CreatedAt = now

	stmt, err := db.Prepare("INSERT INTO bought_items (to_id, from_id, challenge_id, message, created_at) VALUES($1,$2,$3,$4,$5);")
	if err != nil {
		log.Printf("create prepare statement error: %v", err)
		return err
	}

	res, err := stmt.Exec(c.ToID, c.FromID, c.ChallengeID, c.Message, c.CreatedAt)
	if err != nil {
		log.Printf("exec statement error: %v", err)
		return err
	}

	c.ID, err = res.LastInsertId()
	if err != nil {
		log.Printf("last insert id error: %v", err)
		return err
	}

	log.Printf("ChallengeRequest successfully created with id %v", c.ID)

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

//UpdateStatus func updates the status
func (c *ChallengeRequest) UpdateStatus(status string) (int, error) {
	count, err := c.Count("WHERE id=$1 AND to_id=$2", c.ID, c.ToID)
	if err != nil {
		log.Printf("count challenge request err: %v", err)
		return http.StatusInternalServerError, errors.New("Sever error")
	}

	if count != 1 {
		log.Printf("challenge rewuest count: %v, must be 1", count)
		return http.StatusBadRequest, errors.New("Invalid challenge_id and to_id")
	}

	stmt, err := db.Prepare("UPDATE challenge_requests SET (status) VALUES($1);")
	if err != nil {
		log.Printf("create prepare statement error: %v", err)
		return http.StatusInternalServerError, errors.New("Sever error")
	}

	res, err := stmt.Exec(status)
	if err != nil {
		log.Printf("exec statement error: %v", err)
		return http.StatusInternalServerError, errors.New("Sever error")
	}

	c.ID, err = res.LastInsertId()
	if err != nil {
		log.Printf("last insert id error: %v", err)
		return http.StatusInternalServerError, errors.New("Sever error")
	}

	log.Printf("ChallengeRequest successfully created with id %v", c.ID)

	affected, err := res.RowsAffected()
	if err != nil {
		log.Printf("rows effected error: %v", err)
		return http.StatusInternalServerError, errors.New("Sever error")
	}
	if affected == 0 {
		log.Printf("rows effected -> %v", affected)
		return http.StatusNotFound, errors.New("Not found")
	}

	return 0, nil
}
