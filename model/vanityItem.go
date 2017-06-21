package model

import "log"

//VanityItem struct is a model/schema for vanity_item table
type VanityItem struct {
	ID        int64  `json:"id" sql:"id"`
	Name      string `json:"name" sql:"name"`
	Amount    string `json:"amount" sql:"amount"`
	Currency  string `json:"currency" sql:"currency"`
	Coins     string `json:"coins" sql:"coins"`
	CreatedAt string `json:"created_at" sql:"created_at"`
}

//Get func fetches the vanity items from the db
func (v *VanityItem) Get(whereClause string, args ...interface{}) ([]*VanityItem, error) {
	vanityItemList := []*VanityItem{}
	rows, err := db.Query("SELECT id, name, amount, currency, coins, created_at FROM vanity_items "+whereClause+" ORDER BY created_at DESC;", args...)
	if err != nil {
		log.Printf("Get users: sql error %v", err)
		return nil, err
	}
	for rows.Next() {
		vanityItem := VanityItem{}
		if err = rows.Scan(&vanityItem.ID, &vanityItem.Name, &vanityItem.Amount, &vanityItem.Currency, &vanityItem.Coins, &vanityItem.CreatedAt); err != nil {
			log.Printf("scanning row to struct error: %v", err)
			return nil, err
		}

		vanityItemList = append(vanityItemList, &vanityItem)
	}
	return vanityItemList, nil
}

//Count func counts the total vanity_items in the database
func (v *VanityItem) Count(whereClause string, args ...interface{}) (int64, error) {
	var count int64

	if err := db.QueryRow("SELECT COUNT(id) FROM vanity_items "+whereClause+";", args...).Scan(&count); err != nil {
		log.Printf("Count levels: sql error %v", err)
		return count, err
	}

	return count, nil
}
