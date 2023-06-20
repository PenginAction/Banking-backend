package database

import (
	"bank/config"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var Db *sql.DB

func init(){
	cmd := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", 
		config.Config.DbUser,
		config.Config.DbPass,
		config.Config.DbLocal,
		config.Config.DbPort,
		config.Config.DbName,
	)

	var  err error
	Db, err = sql.Open("mysql",cmd)
	if err != nil{
		log.Fatalf("データベース接続に失敗しました: %v", err)
	}
}