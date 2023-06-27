package config

import (
	"log"
	"os"
	"strconv"
)

type Configs struct {
	DbUser      string
	DbPass      string
	DbLocal     string
	DbPort      string
	DbName      string
	DbNameTest  string
	BankLogfile string
	Port        int
}

var Config Configs

func init() {
	port, err := strconv.Atoi(os.Getenv("WEBPORT"))
	if err != nil {
		log.Printf("ポートを整数に変換できませんでした: %v", err)
	}

	Config = Configs{
		DbUser:      os.Getenv("DB_USER"),
		DbPass:      os.Getenv("DB_PASSWORD"),
		DbLocal:     os.Getenv("DB_LOCAL"),
		DbPort:      os.Getenv("DB_PORT"),
		DbName:      os.Getenv("DB_NAME"),
		DbNameTest:  os.Getenv("DB_NAME_TEST"),
		BankLogfile: os.Getenv("LOGFILE"),
		Port:        port,
	}
}
