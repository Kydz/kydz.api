package utils

import (
	"log"
	"math/rand"
	"strconv"
	"time"
)

const runes = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
const (
	width = 6
	mask = 1 << width - 1
	segments = 64 / width
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

func GenerateRandomString(n int) string {
	rand.Seed(time.Now().UnixNano())
	source := rand.Int63()
	re := make([]byte, n)

	for i, remaining := 0, segments; i < n ; {
		if remaining <= 0 {
			source, remaining = rand.Int63(), segments
		}
		if idx := int(source & mask); idx < len(runes) - 1 {
			re[i] = runes[idx]
			i++
		}
		source >>= width
		remaining--
	}

	return string(re)
}
