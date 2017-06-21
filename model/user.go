package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/challengr/middleware"
	jwt "github.com/dgrijalva/jwt-go"
)

//User struct is a model/schema for user table
type User struct {
	ID             int64      `json:"id" sql:"id"`
	Name           string     `json:"name" sql:"name"`
	Email          string     `json:"email" sql:"email"`
	FacebookUserID string     `json:"facebook_user_id" sql:"facebook_user_id"`
	Role           string     `json:"role" sql:"role"`
	Gender         string     `json:"gender" sql:"gender"`
	DOB            string     `json:"date_of_birth" sql:"date_of_birth"`
	Weight         *float32   `json:"weight" sql:"weight"`
	CreatedAt      *time.Time `json:"created_at" sql:"created_at"`
	UpdatedAt      *time.Time `json:"updated_at" sql:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at" sql:"deleted_at"`

	Token    string    `json:"token,omitempty" sql:"-"`
	Location *geometry `json:"geo_coords" sql:"-"`
	LevelID  int64     `json:"-" sql:"level_id"`

	//User        *User         `json:"user" sql:"-"`
	TotalPost   int64         `json:"total_post" sql:"-"`
	Level       *Level        `json:"level" sql:"-"`
	BoughtItems []*BoughtItem `json:"bought_items" sql:"-"` //this will be fetched via user_id
	Score       *Score        `json:"score" sql:"score"`

	Payload map[string]interface{} `json:"-"`
}

//CreateTokenString func creates a new jwt token
func (u *User) CreateTokenString() string {
	// Embed User information to `token`
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), &middleware.JWTUser{
		ID:             u.ID,
		FacebookUserID: u.FacebookUserID,
		Weight:         *u.Weight,
		Role:           u.Role,
	})
	// token -> string. Only server knows this secret (foobar).
	tokenstring, err := token.SignedString([]byte(JWTSecret))
	if err != nil {
		log.Fatalln(err)
	}
	return tokenstring
}

//ParseNotAllowedJSON unmarshalls JSON payload to struct payload and fields. Plus parses the JSON payload.
func (u *User) ParseNotAllowedJSON() []string {
	errSlice := []string{}

	delete(u.Payload, "name")
	delete(u.Payload, "role")
	delete(u.Payload, "gender")
	delete(u.Payload, "date_of_birth")

	for key := range u.Payload {
		errSlice = append(errSlice, key)
	}

	return errSlice
}

//Create func inserts a new user in db
func (u *User) Create() error {
	now := time.Now()
	u.CreatedAt = &now
	u.Role = roleUser

	stmt, err := db.Prepare("INSERT INTO users(name, email, facebook_user_id, role, gender, date_of_birth, created_at) VALUES($1,$2,$3,$4,$5,$6,$7);")
	if err != nil {
		log.Printf("create user prepare statement error: %v", err)
		return err
	}

	res, err := stmt.Exec(u.Name, u.Email, u.FacebookUserID, u.Role, u.Gender, u.DOB, u.CreatedAt)
	if err != nil {
		log.Printf("Create user: exec statement error: %v", err)
		return err
	}

	u.ID, err = res.LastInsertId()
	if err != nil {
		log.Printf("last insert id error: %v", err)
		return err
	}

	log.Printf("user successfully created with id %v", u.ID)

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

//Get func fetches the users from db
func (u *User) Get(whereClause string, args ...interface{}) ([]*User, error) {
	userList := []*User{}
	rows, err := db.Query(`SELECT id, name, email, facebook_user_id, role, (SELECT COUNT(posts.id) FROM posts WHERE posts.user_id=users.user_id) AS total_post, (SELECT row_to_json(levels) FROM levels WHERE levels.id=users.level_id) as level, 
	(SELECT array_to_json(array_agg(bought_items)) FROM bought_items WHERE bought_items.user_id=users.user_id) as bought_items,
	(SELECT id, exp, coins, likes_remaining, created_at FROM scores WHERE scores.user_id=users.id) as score, created_at, updated_at, deleted_at FROM users `+whereClause, args...)
	if err != nil {
		log.Printf("Get users: sql error %v", err)
		return nil, err
	}
	for rows.Next() {
		user := User{}
		levelStr := ""
		boughtItemsStr := ""
		scoreStr := ""
		if err = rows.Scan(&u.ID, &u.Name, &u.Email, &u.FacebookUserID, &u.Role, &u.TotalPost, &levelStr, &boughtItemsStr, &scoreStr, &u.CreatedAt, &u.UpdatedAt, &u.DeletedAt); err != nil {
			log.Printf("scanning row to struct error: %v", err)
			return nil, err
		}

		if err = json.Unmarshal([]byte(scoreStr), &u.Score); err != nil {
			log.Printf("Unmarshaling of score subquery error: %v", err)
			return nil, err
		}

		if err = json.Unmarshal([]byte(levelStr), &u.Level); err != nil {
			log.Printf("Unmarshaling of level subquery error: %v", err)
			return nil, err
		}

		if err = json.Unmarshal([]byte(boughtItemsStr), &u.BoughtItems); err != nil {
			log.Printf("Unmarshaling of bought items subquery error: %v", err)
			return nil, err
		}
		userList = append(userList, &user)
	}
	return userList, nil
}

//Count func counts the users from db
func (u *User) Count(whereClause string, args ...interface{}) (int64, error) {
	var count int64

	if err := db.QueryRow("SELECT COUNT(id) FROM users "+whereClause+";", args...).Scan(&count); err != nil {
		log.Printf("Count users: sql error %v", err)
		return count, err
	}

	return count, nil
}

//UpdateFBFields func updates the name of the user
func (u *User) UpdateFBFields(updateQuery string, args ...interface{}) error {
	stmt, err := db.Prepare(updateQuery)
	if err != nil {
		log.Printf("UPDATE user prepare statement error: %v", err)
		return err
	}

	res, err := stmt.Exec(args...)
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

	return err
}

//Update func updates a user in db
func (u *User) Update() error {
	sets := []string{}
	values := make(map[int]interface{})
	index := 0

	if u.Name != "" {
		values[index] = u.Name
		index = index + 1
		sets = append(sets, "name=$"+strconv.Itoa(index))
	}

	if u.Email != "" {
		values[index] = u.Email
		index = index + 1
		sets = append(sets, "email=$"+strconv.Itoa(index))
	}

	if u.FacebookUserID != "" {
		values[index] = u.FacebookUserID
		index = index + 1
		sets = append(sets, "facebook_user_id=$"+strconv.Itoa(index))
	}

	if u.LevelID != 0 {
		values[index] = u.LevelID
		index = index + 1
		sets = append(sets, "level_id=$"+strconv.Itoa(index))
	}

	if u.Role != "" {
		values[index] = u.Role
		index = index + 1
		sets = append(sets, "role=$"+strconv.Itoa(index))
	}

	if u.Gender != "" {
		values[index] = u.Gender
		index = index + 1
		sets = append(sets, "gender=$"+strconv.Itoa(index))
	}

	if u.DOB != "" {
		values[index] = u.DOB
		index = index + 1
		sets = append(sets, "date_of_birth=$"+strconv.Itoa(index))
	}

	if u.Weight != nil {
		values[index] = *u.Weight
		index = index + 1
		sets = append(sets, "weight=$"+strconv.Itoa(index))
	}

	if u.UpdatedAt != nil {
		values[index] = u.UpdatedAt
		index = index + 1
		sets = append(sets, "updated_at=$"+strconv.Itoa(index))
	}

	stmt, err := db.Prepare("UPDATE users SET " + strings.Join(sets, ", ") + " WHERE id=" + fmt.Sprintf("%v", u.ID) + " AND deleted IS NULL;")
	if err != nil {
		log.Printf("UPDATE user prepare statement error: %v", err)
		return err
	}

	argsValues := make([]interface{}, len(values))
	for k, v := range values {
		argsValues[k] = v
	}

	res, err := stmt.Exec(argsValues...)
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

	return nil
}

//Delete func deletes the user. Delete meaning it doesnt purge it. Just hides it.
func (u *User) Delete() error {
	count, err := u.Count("WHERE id=$1 AND deleted_at IS NULL", u.ID)
	if err != nil {
		log.Printf("User delete: error on fetching User record count: %v", err)
		return err
	}

	if count == 0 {
		err = fmt.Errorf("User account not found-> id %v, total found %v", u.ID, count)
		log.Printf("%v", err)
		return err
	} else if count == 1 {
		stmt, err := db.Prepare("UPDATE users SET deleted_at=$1 WHERE id=$2;")
		if err != nil {
			log.Printf("create prepare statement error: %v", err)
			return err
		}

		res, err := stmt.Exec(time.Now(), u.ID)
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
		err = fmt.Errorf("multiple Users found-> id %v, total found %v", u.ID, count)
		log.Printf("%v", err)
		return err
	}

	return nil
}
