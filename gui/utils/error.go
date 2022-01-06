package utils

import "github.com/Pauloo27/logger"

func HandleError(err error) {
	if err != nil {
		logger.Fatal(err)
	}
}
