package model

import (
	"errors"
	"fmt"
	"log"
	"time"
)

//OneSignal struct is a model/schema for one_signal table
type OneSignal struct {
	ID        int64      `json:"id" sql:"id"`
	UserID    int64      `json:"user_id" sql:"user_id"`
	Imei      string     `json:"imei" sql:"imei"`
	PlayerID  string     `json:"player_id" sql:"player_id"`
	CreatedAt *time.Time `json:"created_at" sql:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" sql:"updated_at"`

	Payload map[string]interface{} `json:"-"`
}

//ParseNotAllowedJSON unmarshalls JSON payload to struct payload and fields. Plus parses the JSON payload.
func (o *OneSignal) ParseNotAllowedJSON() []string {
	errSlice := []string{}

	delete(o.Payload, "user_id")
	delete(o.Payload, "imei")
	delete(o.Payload, "player_id")

	for key := range o.Payload {
		errSlice = append(errSlice, key)
	}

	return errSlice
}

//Validate func validates a post payload data
func (o *OneSignal) Validate() []string {
	errSlice := []string{}

	if o.UserID < 1 {
		errSlice = append(errSlice, "user_id")
	}

	if o.Imei == "" {
		errSlice = append(errSlice, "imei")
	}

	if o.PlayerID == "" {
		errSlice = append(errSlice, "player_id")
	}

	return errSlice
}

//Get func fetches the onesignal records of a user from db
func (o *OneSignal) Get(whereClause string, args ...interface{}) ([]*OneSignal, error) {
	oneSignalList := []*OneSignal{}
	rows, err := db.Query("SELECT id, user_id, imei, player_id, created_at, updated_at FROM users "+whereClause+" ORDER BY created_at DESC;", args...)
	if err != nil {
		log.Printf("Get users: sql error %v", err)
		return nil, err
	}
	for rows.Next() {
		oneSignal := OneSignal{}
		if err = rows.Scan(&o.ID, &o.UserID, &o.Imei, &o.PlayerID, &o.CreatedAt, &o.UpdatedAt); err != nil {
			log.Printf("scanning row to struct error: %v", err)
			return nil, err
		}

		oneSignalList = append(oneSignalList, &oneSignal)
	}
	return oneSignalList, nil
}

//Count func counts the users from db
func (o *OneSignal) Count(whereClause string, args ...interface{}) (int64, error) {
	var count int64

	if err := db.QueryRow("SELECT COUNT(id) FROM onesignal "+whereClause+";", args...).Scan(&count); err != nil {
		log.Printf("Count onesignal: sql error %v", err)
		return count, err
	}

	return count, nil
}

//Upsert func inserts or updates the onesignal info in the db of the user
func (o *OneSignal) Upsert() error {
	count, err := o.Count("WHERE user_id=$1 AND imei=$2", o.UserID, o.Imei)
	if err != nil {
		log.Printf("Onesignal upsert: error on fetching onesignal record counts: %v", err)
		return err
	}

	now := time.Now()

	if count == 0 {
		o.CreatedAt = &now
		stmt, err := db.Prepare("INSERT INTO users(user_id, imei, player_id, created_at) VALUES($1,$2,$3,$4);")
		if err != nil {
			log.Printf("create prepare statement error: %v", err)
			return err
		}

		res, err := stmt.Exec(o.UserID, o.Imei, o.PlayerID, o.CreatedAt)
		if err != nil {
			log.Printf("exec statement error: %v", err)
			return err
		}

		o.ID, err = res.LastInsertId()
		if err != nil {
			log.Printf("last insert id error: %v", err)
			return err
		}

		log.Printf("onesignal successfully created with id %v", o.ID)

		affected, err := res.RowsAffected()
		if err != nil {
			log.Printf("rows effected error: %v", err)
			return err
		}
		if affected == 0 {
			log.Printf("rows effected -> %v", affected)
			return errors.New("Server error")
		}
	} else if count == 1 {
		stmt, err := db.Prepare("INSERT INTO users(player_id, updated_at) VALUES($1,$2);")
		if err != nil {
			log.Printf("create prepare statement error: %v", err)
			return err
		}

		o.UpdatedAt = &now

		res, err := stmt.Exec(o.PlayerID, o.UpdatedAt)
		if err != nil {
			log.Printf("exec statement error: %v", err)
			return err
		}

		affected, err := res.RowsAffected()
		if err != nil {
			log.Printf("rows effected error: %v", err)
			return err
		}
		if affected == 0 {
			log.Printf("rows effected -> %v", affected)
			return errors.New("Server error")
		}
	} else {
		err = fmt.Errorf("Multiple oneSignal record found for -> user_id %v, imei %v, total found %v", o.UserID, o.Imei, count)
		log.Printf("%v", err)
		return err
	}

	return nil
}

//Delete func deletes the onesignal record of the user
func (o *OneSignal) Delete() error {
	count, err := o.Count("WHERE user_id=$1 AND imei=$2", o.UserID, o.Imei)
	if err != nil {
		log.Printf("Onesignal delete: error on fetching onesignal record count: %v", err)
		return err
	}

	if count == 0 {
		err = fmt.Errorf("Onesignal account not found-> user_id %v, imei %v, total found %v", o.UserID, o.Imei, count)
		log.Printf("%v", err)
		return err
	} else if count == 1 {
		stmt, err := db.Prepare("DELETE FROM onesignal WHERE user_id=$1 AND imei=$2;")
		if err != nil {
			log.Printf("create prepare statement error: %v", err)
			return err
		}

		res, err := stmt.Exec(o.UserID, o.Imei)
		if err != nil {
			log.Printf("exec statement error: %v", err)
			return err
		}

		affected, err := res.RowsAffected()
		if err != nil {
			log.Printf("rows effected error: %v", err)
			return err
		}
		if affected == 0 {
			log.Printf("rows effected -> %v", affected)
			return errors.New("Server error")
		}
	} else {
		err = fmt.Errorf("Onesignal account multiple record found-> user_id %v, imei %v, total found %v", o.UserID, o.Imei, count)
		log.Printf("%v", err)
		return err
	}

	return nil
}
