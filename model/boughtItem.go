package model

import (
	"encoding/json"
	"log"
	"time"
)

//BoughtItem is a model/schema for a bought_item table
type BoughtItem struct {
	ID           int64       `json:"id" sql:"id"`
	VanityItemID int64       `json:"-" sql:"vanity_item_id"`
	UserID       int64       `json:"-" sql:"user_id"`
	VanityItem   *VanityItem `json:"vanity_item" sql:"-"`
	LevelID      int64       `json:"-" sql:"level_id"`
	Level        *Level      `json:"level" sql:"-"`
	CreatedAt    time.Time   `json:"created_at" sql:"created_at"`
}

//Get func fetches the vanity items from the db
func (b *BoughtItem) Get(whereClause string, args ...interface{}) ([]*BoughtItem, error) {
	boughtItemList := []*BoughtItem{}
	rows, err := db.Query(`SELECT id, (SELECT row_to_json(vanity_items) FROM vanity_items WHERE vanity_items.id=bought_items.vanity_item_id) as vanity_items, 
	(SELECT row_to_json(levels) FROM levels WHERE levels.id=bought_item.level_id) as level,  created_at FROM bought_items `+whereClause+" ORDER BY created_at DESC;", args...)
	if err != nil {
		log.Printf("Get users: sql error %v", err)
		return nil, err
	}
	for rows.Next() {
		boughtItem := BoughtItem{}
		vanityItemsStr := ""
		levelStr := ""
		if err = rows.Scan(&boughtItem.ID, &vanityItemsStr, &levelStr, &boughtItem.CreatedAt); err != nil {
			log.Printf("scanning row to struct error: %v", err)
			return nil, err
		}

		if err = json.Unmarshal([]byte(vanityItemsStr), &boughtItem.VanityItem); err != nil {
			log.Printf("Unmarshaling of subquery error: %v", err)
			return nil, err
		}

		if err = json.Unmarshal([]byte(vanityItemsStr), &boughtItem.Level); err != nil {
			log.Printf("Unmarshaling of subquery error: %v", err)
			return nil, err
		}

		boughtItemList = append(boughtItemList, &boughtItem)
	}
	return boughtItemList, nil
}

//PurchaseVanityItem func adds a vanity item to the table
func (b *BoughtItem) PurchaseVanityItem() (int, error) {
	return 0, nil
}

//PurchaseLevel func adds a level to the table
func (b *BoughtItem) PurchaseLevel() (int, error) {
	return 0, nil
}
