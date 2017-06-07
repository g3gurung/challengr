package model

//Post struct is a model/schema for post table
type Post struct {
	ID          int64  `json:"id" sql:"id"`
	FileURL     string `json:"file_url" binding:"required" sql:"file_url"`
	ContentType string `json:"content_type" binding:"required" sql:"content_type"`
	ContentSize int64  `json:"content_size" binding:"required" sql:"content_size"`

	Hearts []*Heart `json:"hearts" sql:"-"`

	Payload map[string]interface{} `json:"-"`
}

//ParseNotAllowedJSON unmarshalls JSON payload to struct payload and fields. Plus parses the JSON payload.
func (p *Post) ParseNotAllowedJSON() []string {
	errSlice := []string{}

	delete(p.Payload, "id")
	delete(p.Payload, "file_url")
	delete(p.Payload, "content_type")
	delete(p.Payload, "content_size")

	for key := range p.Payload {
		errSlice = append(errSlice, key)
	}

	return errSlice
}
