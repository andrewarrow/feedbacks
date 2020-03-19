package models

import "fmt"
import "github.com/jmoiron/sqlx"
import "github.com/dgrijalva/jwt-go"
import "time"

type User struct {
	Id        int    `json:"id"`
	Email     string `json:"email"`
	Flavor    string `json:"flavor"`
	Phrase    string
	CreatedAt int64 `json:"created_at"`
}

const jwtSecret = "changeme-66c9dffa-a8b4-4f47-92e1-062298fcde79"

func (u *User) Encode() string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id": u.Id,
		"email": u.Email,
		"flavor": u.Flavor,
    "nbf": time.Now().Unix(),
})

  tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		fmt.Println(err)
	}
	return tokenString
}

func DecodeUser(s string) *User {
	var user User = User{}

	token, err := jwt.Parse(s, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(jwtSecret), nil
	})

	if err != nil {
		fmt.Println(err)
	return &user
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		user.Id = int(claims["id"].(float64))
		user.Email = claims["email"].(string)
		user.Flavor = claims["flavor"].(string)
	} else {
		fmt.Println(err)
	}

	return &user
}

func SelectUsers(db *sqlx.DB) ([]User, string) {
	users := []User{}
	sql := fmt.Sprintf("SELECT id, email, UNIX_TIMESTAMP(created_at) as createdat from users order by created_at desc")
	err := db.Select(&users, sql)
	s := ""
	if err != nil {
		s = err.Error()
	}

	return users, s
}
func InsertUser(db *sqlx.DB, u *User) string {
	_, err := db.NamedExec("INSERT INTO users (email, flavor, phrase) values (:email, :flavor, SHA1(:phrase))",
		map[string]interface{}{"email": u.Email, "flavor": u.Flavor,
			"phrase": u.Phrase})
	if err != nil {
		return err.Error()
	}
	return ""
}
