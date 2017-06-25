package model

import "log"

//GroupChallengeRequest struct is a schema or model for a group_challenge_requests table
type GroupChallengeRequest struct {
	ID            string  `json:"id" sql:"id"`
	FromID        int64   `json:"from_id,omitempty" sql:"from_id"`
	From          *User   `json:"from" sql:"from"`
	ToIDs         []int64 `json:"to_ids,omitempty" sql:"to_ids"`
	ToUsers       []*User `json:"to_users" sql:"to_users"`
	AcceptedIDs   []int64 `json:"accepted_ids,omitempty" sql:"accepted_ids"`
	AcceptedUsers []*User `json:"accepted_users" sql:"accepted_users"`
	Message       string  `json:"message" sql:"message"`

	Payload map[string]interface{} `json:"-"`
}

//ParseNotAllowedJSON unmarshalls JSON payload to struct payload and fields. Plus parses the JSON payload.
func (g *GroupChallengeRequest) ParseNotAllowedJSON() []string {
	errSlice := []string{}

	delete(g.Payload, "from_id")
	delete(g.Payload, "to_ids")
	delete(g.Payload, "message")

	for key := range g.Payload {
		errSlice = append(errSlice, key)
	}

	return errSlice
}

//PostValidate func validates incoming post payload fields
func (g *GroupChallengeRequest) PostValidate() []string {
	errSlice := []string{}

	if g.FromID <= 0 {
		errSlice = append(errSlice, "from_id")
	}

	if len(g.ToIDs) == 0 {
		errSlice = append(errSlice, "to_id")
	}

	return errSlice
}

//Count func counts the users from db
func (g *GroupChallengeRequest) Count(whereClause string, args ...interface{}) (int64, error) {
	var count int64

	if err := db.QueryRow("SELECT COUNT(id) FROM group_challenge_requests "+whereClause+";", args...).Scan(&count); err != nil {
		log.Printf("Count levels: sql error %v", err)
		return count, err
	}

	return count, nil
}
