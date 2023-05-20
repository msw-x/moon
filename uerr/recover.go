package uerr

import "fmt"

func Recover(onError func(string)) {
	if r := recover(); r != nil {
		if onError != nil {
			onError(fmt.Sprint(r))
		}
	}
}
