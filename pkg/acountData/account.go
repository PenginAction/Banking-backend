package acountdata

import (
	"bank/pkg/database"
	"log"
)

type Account struct {
	AccountID int     `json:"acount_id"`
	Balance   float64 `json:"balance"`
	UserID    int     `json:"user_id"`
}

func GetAccountByUserId(id int) (*Account, error) {
	var account Account
	statement := `SELECT * FROM Account WHERE UserID = ?`
	err := database.Db.QueryRow(statement, id).Scan(&account.AccountID, &account.Balance, &account.UserID)
	if err != nil {
		log.Printf("ユーザーIDに対応するアカウントの取得に失敗しました:%s", err.Error())
		return nil, err
	}
	return &account, nil
}

func UpdateBalance(accountId int, newBalance float64) error {
	statement := `UPDATE Account SET Balance = ? WHERE AccountID = ?`
	_, err := database.Db.Exec(statement, newBalance, accountId)
	if err != nil {
		log.Printf("口座の残高の更新に失敗しました: %s", err.Error())
		return err
	}
	return nil
}
