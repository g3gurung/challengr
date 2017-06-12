package model

import "log"

//Level struct is a model/schema for a level table
type Level struct {
	ID   int64  `json:"id" sql:"id"`
	Name string `json:"name" sql:"name"`
}

//Count func counts the users from db
func (l *Level) Count(whereClause string, args ...interface{}) (int64, error) {
	var count int64

	if err := db.QueryRow("SELECT COUNT(id) FROM levels "+whereClause+";", args...).Scan(&count); err != nil {
		log.Printf("Count levels: sql error %v", err)
		return count, err
	}

	return count, nil
}
