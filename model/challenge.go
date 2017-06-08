package model

import "time"

//Challenge struct is a model/schema for a challenge table
type Challenge struct {
	ID          int64      `json:"id" sql:"id"`
	Name        *string    `json:"name" sql:"name"`
	Description *string    `json:"description" sql:"description"`
	Location    *geoCoords `json:"geo_coords" sql:"-"`
	Status      string     `json:"status" sql:"status"`
	CreatedAt   *time.Time `json:"created_at" sql:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at" sql:"updated_at"`

	Payload map[string]interface{} `json:"-"`
}

//ParseNotAllowedJSON unmarshalls JSON payload to struct payload and fields. Plus parses the JSON payload.
func (c *Challenge) ParseNotAllowedJSON() []string {
	errSlice := []string{}

	delete(c.Payload, "name")
	delete(c.Payload, "description")
	delete(c.Payload, "geo_coords")

	for key := range c.Payload {
		errSlice = append(errSlice, key)
	}

	return errSlice
}
