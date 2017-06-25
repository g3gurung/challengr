package model

import (
	"encoding/json"
	"errors"
	"log"
	"time"
)

//BoughtItem is a model/schema for a bought_item table
type BoughtItem struct {
	ID           int64       `json:"id" sql:"id"`
	VanityItemID *int64      `json:"vanity_item_id,omitempty" sql:"vanity_item_id"`
	UserID       int64       `json:"-" sql:"user_id"`
	VanityItem   *VanityItem `json:"vanity_item" sql:"-"`
	LevelID      *int64      `json:"level_id,omitempty" sql:"level_id"`
	Level        *Level      `json:"level" sql:"-"`
	Amount       *string     `json:"amount" sql:"amount"`
	Currency     *string     `json:"currency" sql:"currency"`
	CreatedAt    time.Time   `json:"created_at" sql:"created_at"`

	Payload map[string]interface{} `json:"-"`
}

//ParseNotAllowedJSON unmarshalls JSON payload to struct payload and fields. Plus parses the JSON payload.
func (b *BoughtItem) ParseNotAllowedJSON() []string {
	errSlice := []string{}

	delete(b.Payload, "amount")
	delete(b.Payload, "currency")
	delete(b.Payload, "level_id")
	delete(b.Payload, "vanity_item_id")

	for key := range b.Payload {
		errSlice = append(errSlice, key)
	}

	return errSlice
}

//PostValidate func validates incoming post payload fields
func (b *BoughtItem) PostValidate() []string {
	errSlice := []string{}

	if *b.Amount == "" {
		errSlice = append(errSlice, "amount")
	}

	if *b.Currency == "" {
		errSlice = append(errSlice, "currency")
	}

	if *b.LevelID <= 0 && *b.VanityItemID <= 0 {
		errSlice = append(errSlice, "vanity_item_id/level_id")
	}

	return errSlice
}

//Get func fetches the vanity items from the db
func (b *BoughtItem) Get(whereClause string, args ...interface{}) ([]*BoughtItem, error) {
	boughtItemList := []*BoughtItem{}
	rows, err := db.Query(`SELECT id, (SELECT row_to_json(vanity_items) FROM vanity_items WHERE vanity_items.id=bought_items.vanity_item_id) as vanity_items, 
	(SELECT row_to_json(levels) FROM levels WHERE levels.id=bought_item.level_id) as level, amount, currency, created_at FROM bought_items `+whereClause+" ORDER BY created_at DESC;", args...)
	if err != nil {
		log.Printf("Get users: sql error %v", err)
		return nil, err
	}
	for rows.Next() {
		boughtItem := BoughtItem{}
		vanityItemsStr := ""
		levelStr := ""
		if err = rows.Scan(&boughtItem.ID, &vanityItemsStr, &levelStr, &boughtItem.Amount, &boughtItem.Currency, &boughtItem.CreatedAt); err != nil {
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

//Create func adds an item to the table
func (b *BoughtItem) Create() error {
	now := time.Now()
	b.CreatedAt = now

	stmt, err := db.Prepare("INSERT INTO bought_items (user_id, vanity_item_id, level_id, amount, currency, created_at) VALUES($1,$2,$3,$4,$5,$6);")
	if err != nil {
		log.Printf("create prepare statement error: %v", err)
		return err
	}

	res, err := stmt.Exec(b.UserID, b.VanityItemID, b.LevelID, b.Amount, b.Currency, b.CreatedAt)
	if err != nil {
		log.Printf("exec statement error: %v", err)
		return err
	}

	b.ID, err = res.LastInsertId()
	if err != nil {
		log.Printf("last insert id error: %v", err)
		return err
	}

	log.Printf("boughtitem successfully created with id %v", b.ID)

	affected, err := res.RowsAffected()
	if err != nil {
		log.Printf("rows effected error: %v", err)
		return err
	}
	if affected == 0 {
		log.Printf("rows effected -> %v", affected)
		return errors.New("Server error")
	}

	return nil
}
