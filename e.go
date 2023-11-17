package main

import "log"

func HandleError(msg string, err error) {
	errorMessage := err.Error()
	log.Println(msg, errorMessage)
}
func HandleErrorMessage(msg string, err error) string {
	errorMessage := err.Error()
	log.Println(msg, errorMessage)
	return errorMessage
}
