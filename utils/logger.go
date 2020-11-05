package utils

import "log"

func HandleFatal(message string, err error) {
	if err != nil {
		log.Fatalf("[FATAL]: %s: %v", message, err)
	}
}
