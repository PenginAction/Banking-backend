package config

import (
	"log"
	"os"

	"github.com/go-ini/ini"
)

type Configs struct {
	DbUser      string
	DbPass      string
	DbLocal     string
	DbPort      string
	DbName      string
	BankLogfile string
	Port        int
}

var Config Configs

func init() {
	cfg, err := ini.Load("config.ini")
	if err != nil {
		log.Printf("ファイルの読み込みに失敗しました: %v", err)
		os.Exit(1)
	}

	Config = Configs{
		DbUser:      os.Getenv("DB_USER"),
		DbPass:      os.Getenv("DB_PASSWORD"),
		DbLocal:     os.Getenv("DB_LOCAL"),
		DbPort:      os.Getenv("DB_PORT"),
		DbName:      os.Getenv("DB_NAME"),
		BankLogfile: cfg.Section("bank").Key("logfile").String(),
		Port:        cfg.Section("web").Key("port").MustInt(),
	}
}
