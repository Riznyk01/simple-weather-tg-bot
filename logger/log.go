package logger

import (
	"log"
	"os"
)

func ForError(err error) {
	var (
		outFileFatal, _ = os.OpenFile("logs/fatal.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0755)
		LogFileFatal    = log.New(outFileFatal, "", log.Ldate|log.Ltime)
	)

	if err != nil {
		LogFileFatal.Fatalln(err)
	}
}
func ForErrorPrint(err error) {
	var (
		outFilePrint, _ = os.OpenFile("logs/info.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0755)
		LogFilePrint    = log.New(outFilePrint, "", log.Ldate|log.Ltime)
	)

	if err != nil {
		LogFilePrint.Println(err)
	}
}
