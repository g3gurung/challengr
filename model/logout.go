package model

//LogOut is a struct which is used for loggin out users from the system based on the database
type LogOut struct {
	UserID int64  `json:"user_id"`
	Imei   string `json:"imei"`

	Payload map[string]interface{} `json:"-"`
}

//ParseNotAllowedJSON unmarshalls JSON payload to struct payload and fields. Plus parses the JSON payload.
func (l *LogOut) ParseNotAllowedJSON() []string {
	errSlice := []string{}

	delete(l.Payload, "imei")

	for key := range l.Payload {
		errSlice = append(errSlice, key)
	}

	return errSlice
}

//Validate func validates a post payload data
func (l *LogOut) Validate() []string {
	errSlice := []string{}

	if l.Imei == "" {
		errSlice = append(errSlice, "imei")
	}

	return errSlice
}
