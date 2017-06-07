package model

import "time"

//Heart struct is model/Schema for heart table
type Heart struct {
	ID        int64     `json:"id" sql:"id"`
	UserID    string    `json:"user_id" sql:"user_id"`
	CreatedAt time.Time `json:"created_at" sql:"created_at"`
}
