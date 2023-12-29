package telegram

import "fmt"

func Italic(v any) string {
	return fmt.Sprintf("***%v***", v)
}

func Bold(v any) string {
	return fmt.Sprintf("___%v___", v)
}

func Monospace(v any) string {
	return fmt.Sprintf("```%```", v)
}

func Ref(v any, url string) string {
	return fmt.Sprintf("[%v](%s)", v, url)
}
