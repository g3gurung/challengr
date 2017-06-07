package model

//level struct is a model/schema for a level table
type level struct {
	ID   int    `json:"id" sql:"id"`
	Name string `json:"name" sql:"name"`
}

//Score struct is a model/schema for a score table
type Score struct {
	ID  int64 `json:"id" sql:"id"`
	Exp int   `json:"exp" sql:"exp"`

	LevelID int `json:"-" sql:"level_id"`

	TotalPost   int32         `json:"total_post" sql:"-"`
	Level       *level        `json:"level" sql:"-"`
	BoughtItems []*BoughtItem `json:"bought_items" sql:"-"`
}
