package model

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

//imei struct is a model/schema for imei table
type imei struct {
	ID string `json:"id" sql:"id"`
}

//User struct is a model/schema for user table
type User struct {
	ID             int64      `json:"id" sql:"id"`
	Name           *string    `json:"name" sql:"name"`
	Email          *string    `json:"email" sql:"email"`
	FacebookUserID *string    `json:"facebook_user_id" sql:"facebook_user_id"`
	Role           string     `json:"-" sql:"role"`
	Gender         *string    `json:"gender" sql:"gender"`
	DOB            *string    `json:"date_of_birth" sql:"date_of_birth"`
	Token          string     `json:"token,omitempty" sql:"-"`
	CreatedAt      *time.Time `json:"created_at" sql:"created_at"`
	UpdatedAt      *time.Time `json:"updated_at" sql:"updated_at"`

	Payload map[string]interface{} `json:"-"`
}

//CreateTokenString func creates a new jwt token
func (u *User) CreateTokenString() string {
	// Embed User information to `token`
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), &JWTUser{
		ID:             u.ID,
		FacebookUserID: *u.FacebookUserID,
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
	delete(u.Payload, "email")
	delete(u.Payload, "facebook_user_id")
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

	stmt, err := db.Prepare("INSERT INTO users(name, email, facebook_user_id, role, gender, date_of_birth, created_at, updated_at) VALUES($1,$2,$3,$4,$5,$6,$7,$8);")
	if err != nil {
		log.Printf("Save ticket: create prepare statement error: %v", err)
		return err
	}

	res, err := stmt.Exec(u.Name, u.Email, u.FacebookUserID, u.Role, u.Gender, u.DOB, u.CreatedAt, u.UpdatedAt)
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
	rows, err := db.Query("SELECT id, name, email, facebook_user_id, role, created_at, updated_at FROM users "+whereClause+" ORDER BY created_at DESC;", args...)
	if err != nil {
		log.Printf("Get users: sql error %v", err)
		return nil, err
	}
	for rows.Next() {
		user := User{}
		if err = rows.Scan(&u.ID, &u.Name, &u.Email, &u.FacebookUserID, &u.Role, &u.CreatedAt, &u.UpdatedAt); err != nil {
			log.Printf("scanning row to struct error: %v", err)
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

//Update func updates a user in db
func (u *User) Update() error {
	sets := []string{}
	values := make(map[int]interface{})
	index := 0

	if *u.Name != "" {
		values[index] = *u.Name
		index = index + 1
		sets = append(sets, "name=$"+strconv.Itoa(index))

	}

	if u.UpdatedAt != nil {
		values[index] = *u.UpdatedAt
		index = index + 1
		sets = append(sets, "email=$"+strconv.Itoa(index))
	}

	stmt, err := db.Prepare("UPDATE hub SET " + strings.Join(sets, ", ") + " WHERE id=" + fmt.Sprintf("%v", u.ID) + ";")
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
