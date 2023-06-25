package acountdata

import (
	"bank/pkg/database"
	transactiondata "bank/pkg/transactionData"
	"log"
)

type Account struct {
	AccountID int     `json:"acount_id"`
	Balance   float64 `json:"balance"`
	UserID    int     `json:"user_id"`
}

func GetAccountByUserId(id int) (*Account, error) {
	var account Account
	statement := `SELECT AccountID, Balance, UserID FROM Account WHERE UserID = ?`
	err := database.Db.QueryRow(statement, id).Scan(&account.AccountID, &account.Balance, &account.UserID)
	if err != nil {
		log.Printf("ユーザーIDに対応するアカウントの取得に失敗しました:%s", err.Error())
		return nil, err
	}
	return &account, nil
}

func GetAccountByAccountId(accountId int) (*Account, error) {
	var account Account
	statement := `SELECT AccountID, Balance, UserID FROM Account WHERE AccountID = ?`
	err := database.Db.QueryRow(statement, accountId).Scan(&account.AccountID, &account.Balance, &account.UserID)
	if err != nil {
		log.Printf("口座番号の取得に失敗しました:%s", err.Error())
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

func TransferFunds(senderAccountId int, recipientAccountId int, amount float64) error {
  
	_, err := database.Db.Exec("UPDATE Account SET Balance = Balance - ? WHERE AccountID = ?", amount, senderAccountId)
	if err != nil {
		log.Printf("送金者のアカウントからの引き出しに失敗しました: %s", err.Error())
		return err
	}

	err = transactiondata.CreateTransaction(senderAccountId, recipientAccountId,  amount, "振込")
	if err != nil {
		log.Printf("送金者の取引ログの作成に失敗しました: %s", err.Error())
		return err
	}
  
	_, err = database.Db.Exec("UPDATE Account SET Balance = Balance + ? WHERE AccountID = ?", amount, recipientAccountId)
	if err != nil {
		log.Printf("受取人のアカウントへの入金に失敗しました: %s", err.Error())
		return err
	}

	err = transactiondata.CreateTransaction(recipientAccountId, recipientAccountId, amount, "振込")
	if err != nil {
		log.Printf("受取人の取引ログの作成に失敗しました: %s", err.Error())
		return err
	}

	return nil
}