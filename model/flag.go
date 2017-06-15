package model

import "time"

//Flag struct is a model/schema for a flags table
type Flag struct {
	ID     string `json:"id" sql:"id"`
	UserID int64  `json:"-" sql:"user_id"`
	PostID int64  `json:"post_id" sql:"post_id"`

	User      *User     `json:"user" sql:"-"`
	CreatedAt time.Time `json:"created_at" sql:"created_at"`
}

//Delete func deletes a flag
func (f *Flag) Delete() error {
	return nil
}

//Create func creates a new flag
func (f *Flag) Create() error {
	return nil
}
