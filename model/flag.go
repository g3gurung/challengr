package model

//Flag struct is a model/schema for a flags table
type Flag struct {
	ID     string `json:"id" sql:"id"`
	UserID string `json:"-" sql:"user_id"`

	User *User `json:"user" sql:"-"`
}

//Delete func deletes a flag
func (f *Flag) Delete() error {
	return nil
}

//Create func creates a new flag
func (f *Flag) Create() error {
	return nil
}
