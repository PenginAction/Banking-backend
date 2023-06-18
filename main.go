package main

import (
	"bank/app/webserver"
	"bank/config"
	"bank/utils"
	
)

func main() {
	utils.Logging(config.Config.BankLogfile)
	webserver.Start()
}

