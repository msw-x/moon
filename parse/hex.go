package parse

import (
	"encoding/hex"

	"github.com/msw-x/moon/uerr"
)

func Hex(s string) ([]byte, error) {
	return hex.DecodeString(s)
}

func HexStrict(s string) []byte {
	b, err := Hex(s)
	uerr.Strict(err, "parse hex")
	return b
}
