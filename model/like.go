package model

import "time"

//Like struct is model/Schema for like table
type Like struct {
	ID        int64     `json:"id" sql:"id"`
	UserID    string    `json:"user_id" sql:"user_id"`
	CreatedAt time.Time `json:"created_at" sql:"created_at"`
}
