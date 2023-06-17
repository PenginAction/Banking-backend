package transactiondata

import (
	"bank/pkg/database"
	"log"
)

type Transaction struct{
	TransactionID int `json:"transaction_id"`
	TransactionAmount float64 `json:"transaction_amount"`
	TransactionType string `json:"transaction_type"`
	TransactionDate string `json:"transaction_date"`
	AccountID string `json:"amount_id"`
}

func CreateTransaction(accountId int, amount float64, transactionType string) error{
	statement := `INSERT INTO Transaction (TransactionAmount, TransactionType, AccountID) VALUES (?, ?, ?)`
	_, err := database.Db.Exec(statement, amount, transactionType, accountId)
	if err != nil{
		log.Printf("新しい取引の作成に失敗しました: %s", err.Error())
		return err
	}
	return nil
}