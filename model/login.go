package model

import "github.com/challengr/lib"

//LogIn is a struct used for logging
type LogIn struct {
	Email         string `json:"email"`
	FacebookToken string `json:"facebook_token"`
	Imei          string `json:"imei"`

	Payload map[string]interface{} `json:"-"`
}

//ParseNotAllowedJSON unmarshalls JSON payload to struct payload and fields. Plus parses the JSON payload.
func (l *LogIn) ParseNotAllowedJSON() []string {
	errSlice := []string{}

	delete(l.Payload, "email")
	delete(l.Payload, "imei")
	delete(l.Payload, "facebook_token")

	for key := range l.Payload {
		errSlice = append(errSlice, key)
	}

	return errSlice
}

//PostValidate func validates a post payload data
func (l *LogIn) PostValidate() []string {
	errSlice := []string{}

	if !lib.ValidateEmail(l.Email) {
		errSlice = append(errSlice, "email")
	}

	if l.FacebookToken == "" {
		errSlice = append(errSlice, "facebook_token")
	}

	if l.Imei == "" {
		errSlice = append(errSlice, "imei")
	}

	return errSlice
}
