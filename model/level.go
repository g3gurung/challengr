package model

//Level struct is a model/schema for a level table
type Level struct {
	ID   int64  `json:"id" sql:"id"`
	Name string `json:"name" sql:"name"`
}
