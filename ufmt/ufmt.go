package ufmt

import (
	"encoding/hex"
	"fmt"
)

func Print(v ...any) {
	fmt.Println(Join(v...))
}

func Printf(format string, v ...any) {
	fmt.Println(fmt.Sprintf(format, v...))
}

func Bool(v bool, yes, no string) string {
	if v {
		return yes
	}
	return no
}

func YesNo(v bool) string {
	return Bool(v, "yes", "no")
}

func OnOff(v bool) string {
	return Bool(v, "on", "off")
}

func EnableDisable(v bool) string {
	return Bool(v, "enable", "disable")
}

func OnlineOffline(v bool) string {
	return Bool(v, "online", "offline")
}

func SuccessFailure(v bool) string {
	return Bool(v, "success", "failure")
}

func UpDown(v bool) string {
	return Bool(v, "up", "down")
}

func CheckBox(v bool) string {
	return Bool(v, "[x]", "[ ]")
}

func Hex(buf []byte) string {
	return hex.EncodeToString(buf)
}
