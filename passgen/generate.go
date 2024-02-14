package passgen

import (
	"math/rand"
	"time"

	"golang.org/x/exp/slices"
)

const letters string = "abcdfghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const numbers string = "0123456789"
const symbols string = "!@#$%&*+_-="

func init() {
	rand.Seed(time.Now().UnixNano())
}

func Valid(password string, withLetters bool, withNumbers bool, withSymbols bool) bool {
	has := func(chars string, enable bool) bool {
		if !enable {
			return true
		}
		for _, r := range []rune(chars) {
			if slices.Contains([]rune(password), r) {
				return true
			}
		}
		return false
	}
	return has(letters, withLetters) && has(numbers, withNumbers) && has(symbols, withSymbols)
}

func Generate(length int, withNumbers bool, withSymbols bool) string {
	chars := letters
	if withNumbers {
		chars = chars + numbers
	}
	if withSymbols {
		chars = chars + symbols
	}
	var password string
	for {
		password = generate(length, chars)
		if Valid(password, true, withNumbers, withSymbols) {
			break
		}
	}
	return password
}

func GenerateHard(length int) string {
	return Generate(length, true, true)
}

func GenerateMild(length int) string {
	return Generate(length, true, false)
}

func generate(length int, chars string) string {
	password := ""
	for i := 0; i < length; i++ {
		password += string([]rune(chars)[rand.Intn(len(chars))])
	}
	return password
}
