package model

import "time"

//Like struct is model/Schema for like table
type Like struct {
	ID        int64     `json:"id" sql:"id"`
	UserID    string    `json:"user_id" sql:"user_id"`
	CreatedAt time.Time `json:"created_at" sql:"created_at"`
}

//Delete func deletes a like
func (l *Like) Delete() error {
	return nil
}

//Create func creates a new like
func (l *Like) Create() error {
	return nil
}
