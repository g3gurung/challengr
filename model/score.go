package model

import (
	"encoding/json"
	"errors"
	"log"
	"time"
)

//Score struct is a model/schema for a score table
type Score struct {
	ID             int64      `json:"id" sql:"id"`
	UserID         int64      `json:"user_id,omitempty" sql:"user_id"`
	Exp            int        `json:"exp" sql:"exp"`
	Coins          int64      `json:"coins" sql:"coins"`
	LikesRemaining int        `json:"like_remaining" sql:"likes_remaining"`    //will be descresed every time when liking posts
	LikesUpdatedAt *time.Time `json:"likes_updated_at" sql:"likes_updated_at"` //this will be set when the likes_remaining becomes zero
	CreatedAt      time.Time  `json:"created_at" sql:"created_at"`
	UpdatedAt      *time.Time `json:"updated_at" sql:"updated_at"`

	LevelID int64 `json:"-" sql:"level_id"`

	User        *User         `json:"user" sql:"-"`
	TotalPost   int32         `json:"total_post" sql:"-"`
	Level       *Level        `json:"level" sql:"-"`
	BoughtItems []*BoughtItem `json:"bought_items" sql:"-"` //this will be fetched via user_id
}

//Create func inserts a new score for a new user
func (s *Score) Create() (int, error) {

	stmt, err := db.Prepare("INSERT INTO scores(user_id, exp, coins, likes_remaining, level_id, created_at) VALUES($1,$2,$3,$4,$5,$6);")
	if err != nil {
		log.Printf("Create score: create prepare statement error: %v", err)
		return 0, err
	}

	res, err := stmt.Exec(s.UserID, s.Exp, s.Coins, s.LikesRemaining, s.LevelID, time.Now())
	if err != nil {
		log.Printf("Create score: exec statement error: %v", err)
		return 0, err
	}

	s.ID, err = res.LastInsertId()
	if err != nil {
		log.Printf("last insert id error: %v", err)
		return 0, err
	}

	log.Printf("score successfully created with id %v", s.ID)

	affected, err := res.RowsAffected()
	if err != nil {
		log.Printf("rows effected error: %v", err)
		return 0, err
	}
	if affected == 0 {
		log.Printf("rows effected -> %v", affected)
		return 500, errors.New("Server error")
	}

	return 0, nil
}

//UpdateExp func updates the experience points in db
func (s *Score) UpdateExp(amount int) (int, error) {
	count, err := s.Count("WHERE id=$1 AND user_id=$2", s.ID, s.UserID)
	if err != nil {
		log.Printf("Score count error: %v", err)
		return 500, errors.New("Server error")
	}

	if count == 0 {
		return 404, errors.New("Score not found")
	}

	if count != 1 {
		return 409, errors.New("Conflict in records detected")
	}

	stmt, err := db.Prepare("UPDATE scores SET exp=$1, updated_at=$2 WHERE id=$3 AND user_id=$4;")
	if err != nil {
		log.Printf("create prepare statement error: %v", err)
		return 500, err
	}

	now := time.Now()
	s.UpdatedAt = &now

	res, err := stmt.Exec(amount, s.UpdatedAt, s.ID, s.UserID)
	if err != nil {
		log.Printf("exec statement error: %v", err)
		return 500, err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		log.Printf("rows effected error: %v", err)
		return 500, err
	}
	if affected == 0 {
		log.Printf("rows effected -> %v", affected)
		return 500, errors.New("Server error")
	}
	return 0, nil
}

//UpdateCoins func updates coins on db
func (s *Score) UpdateCoins(amount int64) (int, error) {
	count, err := s.Count("WHERE id=$1 AND user_id=$2", s.ID, s.UserID)
	if err != nil {
		log.Printf("Score count error: %v", err)
		return 500, errors.New("Server error")
	}

	if count == 0 {
		return 404, errors.New("Score not found")
	}

	if count != 1 {
		return 409, errors.New("Conflict in records detected")
	}

	stmt, err := db.Prepare("UPDATE scores SET coins=$1, updated_at=$2  WHERE id=$3 AND user_id=$4;")
	if err != nil {
		log.Printf("create prepare statement error: %v", err)
		return 500, err
	}

	now := time.Now()
	s.UpdatedAt = &now

	res, err := stmt.Exec(amount, s.UpdatedAt, s.ID, s.UserID)
	if err != nil {
		log.Printf("exec statement error: %v", err)
		return 500, err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		log.Printf("rows effected error: %v", err)
		return 500, err
	}
	if affected == 0 {
		log.Printf("rows effected -> %v", affected)
		return 500, errors.New("Server error")
	}
	return 0, err
}

//ResetLikes func updates likes on db
func (s *Score) ResetLikes(amount int) (int, error) {
	count, err := s.Count("WHERE id=$1 AND user_id=$2 AND likes_remaining=0 AND likes_updated_at < NOW() - INTERVAL '1 hour'", s.ID, s.UserID)
	if err != nil {
		log.Printf("Score count error: %v", err)
		return 500, errors.New("Server error")
	}

	if count == 0 {
		return 405, errors.New("Not allowed")
	}

	if count != 1 {
		return 409, errors.New("Conflict in records detected")
	}

	stmt, err := db.Prepare("UPDATE scores SET likes_remaining=$1 WHERE id=$2 AND user_id=$3;")
	if err != nil {
		log.Printf("create prepare statement error: %v", err)
		return 500, err
	}

	res, err := stmt.Exec(amount, s.ID, s.UserID)
	if err != nil {
		log.Printf("exec statement error: %v", err)
		return 500, err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		log.Printf("rows effected error: %v", err)
		return 500, err
	}
	if affected == 0 {
		log.Printf("rows effected -> %v", affected)
		return 500, errors.New("Server error")
	}
	return 0, err
}

//AddLikes func updates likes on db
func (s *Score) AddLikes(amount int) (int, error) {
	count, err := s.Count("WHERE id=$1 AND user_id=$2", s.ID, s.UserID)
	if err != nil {
		log.Printf("Score count error: %v", err)
		return 500, errors.New("Server error")
	}

	if count == 0 {
		return 405, errors.New("Not allowed")
	}

	if count != 1 {
		return 409, errors.New("Conflict in records detected")
	}

	stmt, err := db.Prepare("UPDATE scores SET likes_remaining=likes_remaining + $1 WHERE id=$2 AND user_id=$3;")
	if err != nil {
		log.Printf("create prepare statement error: %v", err)
		return 500, err
	}

	res, err := stmt.Exec(amount, s.ID, s.UserID)
	if err != nil {
		log.Printf("exec statement error: %v", err)
		return 500, err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		log.Printf("rows effected error: %v", err)
		return 500, err
	}
	if affected == 0 {
		log.Printf("rows effected -> %v", affected)
		return 500, errors.New("Server error")
	}
	return 0, err
}

//DecreaseLikes func subtracts one like from likes_remaining in db
func (s *Score) DecreaseLikes(amount int) (int, error) {
	count, err := s.Count("WHERE id=$1 AND user_id=$2 AND (likes_remaining - $3) >= 0", s.ID, s.UserID, amount)
	if err != nil {
		log.Printf("Score count error: %v", err)
		return 500, errors.New("Server error")
	}

	if count == 0 {
		return 405, errors.New("Not allowed")
	}

	if count != 1 {
		return 409, errors.New("Conflict in records detected")
	}

	//update test set mynum = case when (0 < (mynum - 5)) then (mynum - 5) else (0) end;
	stmt, err := db.Prepare("UPDATE scores SET likes_remaining = likes_remaining-$1 WHERE id=$2 AND user_id=$3;")
	if err != nil {
		log.Printf("create prepare statement error: %v", err)
		return 500, err
	}

	res, err := stmt.Exec(amount, s.ID, s.UserID)
	if err != nil {
		log.Printf("exec statement error: %v", err)
		return 500, err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		log.Printf("rows effected error: %v", err)
		return 500, err
	}
	if affected == 0 {
		log.Printf("rows effected -> %v", affected)
		return 500, errors.New("Server error")
	}
	return 0, err
}

//UpdateLevel func upgrades or degrades to specific level
func (s *Score) UpdateLevel(levelID int64) (int, error) {
	count, err := s.Count("WHERE id=$1 AND user_id=$2 AND $3 IN (SELECT id FROM levels)", s.ID, s.UserID, levelID)
	if err != nil {
		log.Printf("Score count error: %v", err)
		return 500, errors.New("Server error")
	}

	if count == 0 {
		return 404, errors.New("Score/level not found")
	}

	if count != 1 {
		return 409, errors.New("Conflict in records detected")
	}

	stmt, err := db.Prepare("UPDATE scores SET level_id=$1, updated_at=$2  WHERE id=$3 AND user_id=$4;")
	if err != nil {
		log.Printf("create prepare statement error: %v", err)
		return 500, err
	}

	now := time.Now()
	s.UpdatedAt = &now

	res, err := stmt.Exec(levelID, s.UpdatedAt, s.ID, s.UserID)
	if err != nil {
		log.Printf("exec statement error: %v", err)
		return 500, err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		log.Printf("rows effected error: %v", err)
		return 500, err
	}
	if affected == 0 {
		log.Printf("rows effected -> %v", affected)
		return 500, errors.New("Server error")
	}
	return 0, err
}

//Count func counts the users from db
func (s *Score) Count(whereClause string, args ...interface{}) (int64, error) {
	var count int64

	if err := db.QueryRow("SELECT COUNT(id) FROM scores "+whereClause+";", args...).Scan(&count); err != nil {
		log.Printf("Count scores: sql error %v", err)
		return count, err
	}

	return count, nil
}

//Get func fetches the scores of the users
func (s *Score) Get(whereClause string, args ...interface{}) ([]*Score, int, error) {
	scoreList := []*Score{}
	rows, err := db.Query("SELECT id, (SELECT row_to_json(row) FROM (SELECT id, name, role FROM users WHERE user.id=scores.user_id) row) as user, exp, coins, (SELECT row_to_json(levels) FROM levels WHERE levels.id=scores.level_id) as level, (SELECT COUNT(id) FROM posts WHERE posts.user_id=scores.user_id) AS total_post, (SELECT array_to_json(array_agg(bought_items)) FROM bought_items WHERE bought_items.user_id=scores.user_id) as bought_items, created_at, updated_at FROM flags WHERE post_id=posts.id) as flags FROM scores "+whereClause+";", args...)
	if err != nil {
		log.Printf("Get scores: sql error %v", err)
		return nil, 500, err
	}
	for rows.Next() {
		score := Score{}
		levelStr := ""
		boughtItemsStr := ""
		userStr := ""
		if err = rows.Scan(&s.ID, &userStr, &s.Exp, &levelStr, &s.TotalPost, &boughtItemsStr, &s.CreatedAt, &s.UpdatedAt); err != nil {
			log.Printf("scanning row to struct error: %v", err)
			return nil, 500, err
		}

		if err = json.Unmarshal([]byte(userStr), &s.User); err != nil {
			log.Printf("Unmarshaling of user subquery error: %v", err)
			return nil, 500, err
		}

		if err = json.Unmarshal([]byte(levelStr), &s.Level); err != nil {
			log.Printf("Unmarshaling of level subquery error: %v", err)
			return nil, 500, err
		}

		if err = json.Unmarshal([]byte(boughtItemsStr), &s.BoughtItems); err != nil {
			log.Printf("Unmarshaling of bought items subquery error: %v", err)
			return nil, 500, err
		}

		scoreList = append(scoreList, &score)
	}
	return scoreList, 500, nil
}
