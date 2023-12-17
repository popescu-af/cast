package logger

import (
	"fmt"
	"log"
	"os"
	"time"
)

var logsDir = ""
var Log *log.Logger

func init() {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	logsDir = fmt.Sprintf("%s/logs", cwd)
	err = os.MkdirAll(logsDir, os.ModePerm)
	if err != nil {
		panic(err)
	}
	logFileLocation := fmt.Sprintf("%s/log-%s.log", logsDir, time.Now().Format("2006-01-02_15-04-05.000"))
	file, err := os.Create(logFileLocation)
	if err != nil {
		panic(err)
	}
	Log = log.New(file, "", log.Ldate|log.Ltime|log.Lshortfile)
}

func Directory() string {
	return logsDir
}
