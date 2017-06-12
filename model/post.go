package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"
)

//Post struct is a model/schema for post table
type Post struct {
	ID          int64      `json:"id" sql:"id"`
	UserID      int64      `json:"user_id" sql:"user_id"`
	ChallengeID int64      `json:"challenge_id" sql:"challenge_id"`
	LikesNeeded int        `json:"likes_needed" sql:"likes_needed"`
	FileURL     string     `json:"file_url" sql:"file_url"`
	ContentType string     `json:"content_type" sql:"content_type"`
	ContentSize int64      `json:"content_size" sql:"content_size"`
	CreatedAt   *time.Time `json:"created_at" sql:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at" sql:"updated_at"`

	Flags []*Flag `json:"flags" sql:"-"`
	Likes []*Like `json:"likes" sql:"-"`

	Payload map[string]interface{} `json:"-"`
}

//ParseNotAllowedJSON unmarshalls JSON payload to struct payload and fields. Plus parses the JSON payload.
func (p *Post) ParseNotAllowedJSON() []string {
	errSlice := []string{}

	delete(p.Payload, "file_url")
	delete(p.Payload, "content_type")
	delete(p.Payload, "content_size")

	for key := range p.Payload {
		errSlice = append(errSlice, key)
	}

	return errSlice
}

//PostValidate func validates the incoming allowed post fields
func (p *Post) PostValidate() []string {
	errSlice := []string{}

	if p.FileURL == "" {
		errSlice = append(errSlice, "file_url")
	}

	if p.ContentSize == 0 {
		errSlice = append(errSlice, "content_size")
	}

	if p.ContentType == "" {
		errSlice = append(errSlice, "content_type")
	}
	return errSlice
}

//Create func inserts new post in db
func (p *Post) Create() error {
	now := time.Now()
	p.CreatedAt = &now

	stmt, err := db.Prepare("INSERT INTO posts(user_id, likes_needed, challenge_id, file_url, content_type, content_size, created_at) VALUES($1,$2,$3,$4,$5,$6);")
	if err != nil {
		log.Printf("create prepare statement error: %v", err)
		return err
	}

	res, err := stmt.Exec(p.UserID, p.LikesNeeded, p.ChallengeID, p.FileURL, p.ContentType, p.ContentSize, p.CreatedAt)
	if err != nil {
		log.Printf("exec statement error: %v", err)
		return err
	}

	p.ID, err = res.LastInsertId()
	if err != nil {
		log.Printf("last insert id error: %v", err)
		return err
	}

	log.Printf("onesignal successfully created with id %v", p.ID)

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

//Count func counts the total posts in db
func (p *Post) Count(whereClause string, args ...interface{}) (int64, error) {
	var count int64

	if err := db.QueryRow("SELECT COUNT(id) FROM posts "+whereClause+";", args...).Scan(&count); err != nil {
		log.Printf("Count onesignal: sql error %v", err)
		return count, err
	}

	return count, nil
}

//Get func counts the total posts in db
func (p *Post) Get(whereClause string, args ...interface{}) ([]*Post, error) {
	postList := []*Post{}
	rows, err := db.Query("SELECT id, likes_needed, file_url, content_type, content_size, created_at, updated_at, (SELECT array_to_json(array_agg(likes)) FROM likes WHERE post_id=posts.id) as likes, (SELECT array_to_json(array_agg(flags)) FROM flags WHERE post_id=posts.id) as flags FROM posts "+whereClause+" ORDER BY created_at DESC;", args...)
	if err != nil {
		log.Printf("Get users: sql error %v", err)
		return nil, err
	}
	for rows.Next() {
		post := Post{}
		likesStr := ""
		flagsStr := ""
		if err = rows.Scan(&p.ID, &p.LikesNeeded, &p.FileURL, &p.ContentType, &p.ContentSize, &p.CreatedAt, &p.UpdatedAt, &likesStr, &flagsStr); err != nil {
			log.Printf("scanning row to struct error: %v", err)
			return nil, err
		}

		if err = json.Unmarshal([]byte(likesStr), &p.Likes); err != nil {
			log.Printf("Unmarshaling of likes subquery error: %v", err)
			return nil, err
		}

		if err = json.Unmarshal([]byte(flagsStr), &p.Flags); err != nil {
			log.Printf("Unmarshaling of flags subquery error: %v", err)
			return nil, err
		}

		postList = append(postList, &post)
	}
	return postList, nil
}

//Delete func deletes the post record of the user. Delete meaning it doesnt purge it. Just hides it.
func (p *Post) Delete() error {
	count, err := p.Count("WHERE id=$1", p.ID)
	if err != nil {
		log.Printf("Post delete: error on fetching Post record count: %v", err)
		return err
	}

	if count == 0 {
		err = fmt.Errorf("Post account not found-> id %v, total found %v", p.ID, count)
		log.Printf("%v", err)
		return err
	} else if count == 1 {
		stmt, err := db.Prepare("UPDATE posts SET deleted_at=$1 WHERE id=$2;")
		if err != nil {
			log.Printf("create prepare statement error: %v", err)
			return err
		}

		res, err := stmt.Exec(time.Now(), p.ID)
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
		err = fmt.Errorf("multiple posts found-> id %v, total found %v", p.ID, count)
		log.Printf("%v", err)
		return err
	}

	return nil
}
