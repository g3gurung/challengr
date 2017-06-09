package model

import "time"

//Score struct is a model/schema for a score table
type Score struct {
	ID        int64      `json:"id" sql:"id"`
	UserID    int64      `json:"user_id" sql:"user_id"`
	Exp       int        `json:"exp" sql:"exp"`
	Coins     int64      `json:"coins" sql:"coins"`
	CreatedAt *time.Time `json:"created_at" sql:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" sql:"updated_at"`

	LevelID int64 `json:"-" sql:"level_id"`

	TotalPost   int32         `json:"total_post" sql:"-"`
	Level       *Level        `json:"level" sql:"-"`
	BoughtItems []*BoughtItem `json:"bought_items" sql:"-"` //this will be fetched via user_id
}

//AddExp func adds exp based on the amount provided. e.g. 0-99
func (s *Score) AddExp(amount int) error {
	return nil
}

//SubtractExp func subtracts exp based on the amount provided. e.g. 0-99
func (s *Score) SubtractExp(amount int) error {
	return nil
}

//AddCoins func adds coins based on the amount provided.
func (s *Score) AddCoins(amount int) error {
	return nil
}

//SubtractCoins func subtracts coins based on the amount provided.
func (s *Score) SubtractCoins(amount int) error {
	return nil
}

//GoLevel func upgrades or degrades to specific level
func (s *Score) GoLevel(levelID int64) error {
	return nil
}

//Get func fetches the scores of the users
func (s *Score) Get(whereClause string, args ...interface{}) ([]*Score, error) {
	scoreList := []*Score{}
	return scoreList, nil
}
