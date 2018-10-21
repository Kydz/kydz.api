package utils

import (
	"log"
	"strconv"
)

func StringToInteger(value string) int {
	converted, err := strconv.Atoi(value)
	if err != nil {
		log.Fatal(err)
	}
	return converted
}

func IntegerToString(value int) string {
	converted := strconv.Itoa(value)
	return converted
}
