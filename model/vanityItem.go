package model

//VanityItem struct is a model/schema for vanity_item table
type VanityItem struct {
	ID       int64  `json:"id" sql:"id"`
	Name     string `json:"name" sql:"name"`
	Amount   string `json:"amount" sql:"amount"`
	Currency string `json:"currency" sql:"currency"`
	Coins    string `json:"coins" sql:"coins"`
}

//Get func fetches the vanity items from the db
func (v *VanityItem) Get(whereClause string, args ...interface{}) ([]*VanityItem, error) {
	vanityItemList := []*VanityItem{}
	return vanityItemList, nil
}
