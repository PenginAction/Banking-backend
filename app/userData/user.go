package userdata

import (
	"bank/pkg/database"
	"log"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	UserId   int    `json:"user_id"`
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

	statement1 := `INSERT INTO User (Name, Email, Password) VALUES (?, ?, ?)`
	res, err := database.Db.Exec(statement1, name, email, hashPass)
	if err != nil {
		log.Printf("データベースへのユーザー挿入に失敗しました: %s", err.Error())
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Printf("挿入されたデータのIDの取得に失敗しました: %s", err.Error())
		return nil, err
	}

	firstBalance := 0.0
	statement2 := `INSERT INTO Account (UserId, Balance) VALUES (?, ?)`
	_, err = database.Db.Exec(statement2, id, firstBalance)
	if err != nil{
		log.Printf("口座開設できませんでした: %s", err.Error())
		return nil, err
	}

	return &User{
		UserId:   int(id),
		Name:     name,
		Email:    email,
	}, nil

}


func GetAccountByEmail(email string) (*User, error){
	var user User
	statement := `SELECT UserId, Name, Email, Password FROM User WHERE email = ?`
	err := database.Db.QueryRow(statement, email).Scan(&user.UserId, &user.Name, &user.Email, &user.Password)
	if err != nil{
		log.Printf("データベースからメールアドレスの取得に失敗しました: %s", err.Error())
		return nil, err
	}
	return &user, nil
}

func CompareHashAndPassword(userPassword string, password string) error{
	err := bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(password))
	if err != nil{
		log.Printf("認証に失敗しました：パスワードが一致しません．")
		return err
	}
	return nil
}

func GetAccountById(id int) (*User, error){
	var user User
	statement := `SELECT UserId, Name, Email, Password FROM User WHERE UserId = ?`
	err := database.Db.QueryRow(statement, id).Scan(&user.UserId, &user.Name, &user.Email, &user.Password)
	if err != nil {
		log.Printf("データベースからユーザーIDの取得に失敗しました: %s", err.Error())
		return nil, err
	}
	return &user, nil
}