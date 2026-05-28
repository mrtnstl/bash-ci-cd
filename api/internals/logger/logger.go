package logger

import (
	"log"
	"os"

	"example.com/api/internals/utils"
)

func Log(path string, data... string){
	env := utils.GetEnvString("GO_ENV")
	var logRow string

	for _, value := range data {
		logRow += ";" + value
	}
	logRow+="\n"

	if env != "production" {
		log.Println(logRow)
	}
	
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("error opening log file: %v", err)
	}
	defer f.Close()

	builtinLogger := log.New(f, "bash-ci-cd;", log.LstdFlags)
	builtinLogger.Print(logRow)
}