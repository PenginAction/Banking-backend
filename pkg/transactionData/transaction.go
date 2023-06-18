package transactiondata

import (
	"bank/pkg/database"
	"log"
)

type Transaction struct {
	TransactionID     int     `json:"transaction_id"`
	TransactionAmount float64 `json:"transaction_amount"`
	TransactionType   string  `json:"transaction_type"`
	TransactionDate   string  `json:"transaction_date"`
	AccountID         int     `json:"amount_id"`
}

func CreateTransaction(accountId int, amount float64, transactionType string) error {
	statement := `INSERT INTO Transaction (TransactionAmount, TransactionType, AccountID) VALUES (?, ?, ?)`
	_, err := database.Db.Exec(statement, amount, transactionType, accountId)
	if err != nil {
		log.Printf("新しい取引の作成に失敗しました: %s", err.Error())
		return err
	}
	return nil
}

func GetTransactionByAccountId(id int) ([]Transaction, error) {
	var transactions []Transaction
	statement := `SELECT * FROM Transaction WHERE AccountID = ?`
	rows, err := database.Db.Query(statement, id)
	if err != nil {
		log.Printf("アカントIDに対応する取引の取得に失敗しました:%s", err.Error())
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var transaction Transaction
		err := rows.Scan(&transaction.TransactionID, &transaction.TransactionAmount, &transaction.TransactionType, &transaction.TransactionDate, &transaction.AccountID)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return transactions, nil
}
