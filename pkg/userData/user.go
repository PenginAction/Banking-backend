package userdata

import (
	"bank/pkg/database"
	"errors"
	"log"
	"math/rand"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	UserId   int    `json:"user_id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Pin string `json:"pin"`
}

var r *rand.Rand

func init() {
	s := rand.NewSource(time.Now().UnixNano())
	r = rand.New(s)
}

func generateAccountID(length int) string {
	const charset = "1234567890"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[r.Intn(len(charset))]
	}
	return string(result)
}

func generateUniqueAccountID(length int) (string, error) {
	const maxAttempts = 5
	for attempt := 0; attempt < maxAttempts; attempt++ {
		accountID := generateAccountID(length)
		exists, err := checkAccountIDExists(accountID)
		if err != nil {
			return "", err
		}
		if !exists {
			return accountID, nil
		}
	}
	return "", errors.New("複数回の試行後でも一意のアカウントIDを生成できませんでした")
}

func checkAccountIDExists(accountID string) (bool, error) {
	statement := `SELECT COUNT(*) FROM Account WHERE AccountID = ?`
	var count int
	err := database.Db.QueryRow(statement, accountID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func CreateUser(name, email, password, pin string) (*User, error) {
	hashPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("パスワードのハッシュ化に失敗しました: %s", err.Error())
		return nil, err
	}

	hashPin, err := bcrypt.GenerateFromPassword([]byte(pin), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("暗証番号のハッシュ化に失敗しました: %s", err.Error())
		return nil, err
	}

	statement1 := `INSERT INTO User (Name, Email, Password, Pin) VALUES (?, ?, ?, ?)`
	res, err := database.Db.Exec(statement1, name, email, hashPass, hashPin)
	if err != nil {
		log.Printf("データベースへのユーザー挿入に失敗しました: %s", err.Error())
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Printf("挿入されたデータのIDの取得に失敗しました: %s", err.Error())
		return nil, err
	}

	accountID, err := generateUniqueAccountID(10)
	if err != nil {
		log.Printf("口座番号の作成ができません: %s", err.Error())
		return nil, err
	}

	firstBalance := 0.0
	statement2 := `INSERT INTO Account (AccountID, UserId, Balance) VALUES (?, ?, ?)`
	_, err = database.Db.Exec(statement2, accountID, id, firstBalance)
	if err != nil {
		log.Printf("口座開設できませんでした: %s", err.Error())
		return nil, err
	}

	return &User{
		UserId: int(id),
		Name:   name,
		Email:  email,
	}, nil

}

func GetAccountByEmail(email string) (*User, error) {
	var user User
	statement := `SELECT UserId, Name, Email, Password, Pin FROM User WHERE Email = ?`
	err := database.Db.QueryRow(statement, email).Scan(&user.UserId, &user.Name, &user.Email, &user.Password, &user.Pin)
	if err != nil {
		log.Printf("データベースからメールアドレスの取得に失敗しました: %s", err.Error())
		return nil, err
	}
	return &user, nil
}

func CompareHashAndPassword(userPassword, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(password))
	if err != nil {
		log.Printf("認証に失敗しました：パスワードが一致しません．")
		return err
	}
	return nil
}

func GetAccountById(id int) (*User, error) {
	var user User
	statement := `SELECT UserId, Name, Email, Password, Pin FROM User WHERE UserId = ?`
	err := database.Db.QueryRow(statement, id).Scan(&user.UserId, &user.Name, &user.Email, &user.Password, &user.Pin)
	if err != nil {
		log.Printf("データベースからユーザーIDの取得に失敗しました: %s", err.Error())
		return nil, err
	}
	return &user, nil
}

func CompareHashAndPin(userPin, pin string) error{
	err := bcrypt.CompareHashAndPassword([]byte(userPin), []byte(pin))
	if err != nil {
		log.Printf("認証に失敗しました: 暗証番号が一致しません")
		return err
	}
	return nil
}
