package passgen

import (
	"math/rand"
	"time"
)

const voc string = "abcdfghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const numbers string = "0123456789"
const symbols string = "!@#$%&*+_-="

func init() {
	rand.Seed(time.Now().UnixNano())
}

func generate(length int, chars string) string {
	password := ""
	for i := 0; i < length; i++ {
		password += string([]rune(chars)[rand.Intn(len(chars))])
	}
	return password
}

func Generate(length int, hasNumbers bool, hasSymbols bool) string {
	chars := voc
	if hasNumbers {
		chars = chars + numbers
	}
	if hasSymbols {
		chars = chars + symbols
	}
	return generate(length, chars)
}
