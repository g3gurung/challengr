package model

//BoughtItem is a model/schema for a bought_item table
type BoughtItem struct {
	ID       int64  `json:"id" sql:"id"`
	Name     string `json:"name" sql:"name"`
	Currency string `json:"currency" sql:"currency"`
	Amount   string `json:"amount" sql:"amount"`
}
