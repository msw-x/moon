package parse

import (
	"encoding/hex"
	"github.com/msw-x/moon"
)

func Hex(s string) []byte {
	b, err := hex.DecodeString(s)
	moon.Check(err, "ufmt hex")
	return b
}
