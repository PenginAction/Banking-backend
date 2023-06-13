package userdata

import (
	"bank/pkg/database"
	"log"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	UserId   int    `json:"userId"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func CreateAccount(name, email, password string) (*User, error) {
	hashPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("パスワードのハッシュ化に失敗しました: %s", err.Error())
		return nil, err
	}

	statement := `INSERT INTO User (Name, Email, Password) VALUES (?, ?, ?)`
	res, err := database.Db.Exec(statement, name, email, hashPass)
	if err != nil {
		log.Printf("データベースへのユーザー挿入に失敗しました: %s", err.Error())
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Printf("挿入されたデータのIDの取得に失敗しました: %s", err.Error())
		return nil, err
	}

	return &User{
		UserId:   int(id),
		Name:     name,
		Email:    email,
		Password: string(hashPass),
	}, nil

}

