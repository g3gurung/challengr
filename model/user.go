package model

import "time"

//imei struct is a model/schema for imei table
type imei struct {
	ID string `json:"id" sql:"id"`
}

//User struct is a model/schema for user table
type User struct {
	ID            int64      `json:"id" sql:"id"`
	Name          *string    `json:"name" sql:"name"`
	Email         *string    `json:"email" sql:"email"`
	Role          string     `json:"-" sql:"role"`
	Gender        *string    `json:"gender" sql:"gender"`
	DOB           *string    `json:"date_of_birth" sql:"date_of_birth"`
	FacebookToken *string    `json:"facebook_token" sql:"facebook_token"`
	CreatedAt     *time.Time `json:"created_at" sql:"created_at"`
	UpdatedAt     *time.Time `json:"updated_at" sql:"updated_at"`

	Payload map[string]interface{} `json:"-"`
}

//ParseNotAllowedJSON unmarshalls JSON payload to struct payload and fields. Plus parses the JSON payload.
func (u *User) ParseNotAllowedJSON() []string {
	errSlice := []string{}

	delete(u.Payload, "id")
	delete(u.Payload, "name")
	delete(u.Payload, "email")
	delete(u.Payload, "facebook_token")

	for key := range u.Payload {
		errSlice = append(errSlice, key)
	}

	return errSlice
}
