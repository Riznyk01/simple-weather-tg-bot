package logger

import (
	"log"
	"os"
)

var (
	outFileFatal, _ = os.OpenFile("logs/fatal.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0755)
	LogFileFatal    = log.New(outFileFatal, "", 0)

	outFilePrint, _ = os.OpenFile("logs/info.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0755)
	LogFilePrint    = log.New(outFilePrint, "", 0)
)

func ForError(err error) {
	if err != nil {
		LogFileFatal.Fatalln(err)
	}
}
func ForErrorPrint(err error) {
	if err != nil {
		LogFilePrint.Println(err)
	}
}
