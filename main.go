package main

import (
	"bank/pkg/webserver"
	"bank/config"
	"bank/utils"
	
)

func main() {
	utils.Logging(config.Config.BankLogfile)
	webserver.Start()
}

