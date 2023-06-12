package utils

import (
	"io"
	"log"
	"os"
)

func Logging(LogFile string) {
	logfile, err := os.OpenFile(LogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("file=LogFile err=%s", err.Error())
	}

	multiLogOutputs := io.MultiWriter(os.Stdout, logfile)

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.SetOutput(multiLogOutputs)
}
