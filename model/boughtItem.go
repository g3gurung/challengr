package model

//BoughtItem is a model/schema for a bought_item table
type BoughtItem struct {
	ID       int64  `json:"id" sql:"id"`
	Name     string `json:"name" sql:"name"`
	Currency string `json:"currency" sql:"currency"`
	Amount   string `json:"amount" sql:"amount"`

	Payload map[string]interface{} `json:"-"`
}

//ParseNotAllowedJSON unmarshalls JSON payload to struct payload and fields. Plus parses the JSON payload.
func (b *BoughtItem) ParseNotAllowedJSON() []string {
	errSlice := []string{}

	delete(b.Payload, "name")
	delete(b.Payload, "email")
	delete(b.Payload, "currency")
	delete(b.Payload, "ammount")

	for key := range b.Payload {
		errSlice = append(errSlice, key)
	}

	return errSlice
}
