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
		log.Printf("Failed to read file: %v", err)
		os.Exit(1)
	}

	Config = Configs{
		DbUser:      cfg.Section("database").Key("dbuser").String(),
		DbPass:      cfg.Section("database").Key("dbpass").String(),
		DbLocal:     cfg.Section("database").Key("dblocal").String(),
		DbPort:      cfg.Section("database").Key("dbport").String(),
		DbName:      cfg.Section("database").Key("dbname").String(),
		BankLogfile: cfg.Section("bank").Key("logfile").String(),
		Port:        cfg.Section("web").Key("port").MustInt(),
	}
}
