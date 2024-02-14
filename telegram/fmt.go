package telegram

import (
	"fmt"

	botapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Bold(v any) string {
	return fmt.Sprintf("***%v***", v)
}

func Italic(v any) string {
	return fmt.Sprintf("___%v___", v)
}

func Monospace(v any) string {
	return fmt.Sprintf("```%v```", v)
}

func Strikethrough(v any) string {
	return fmt.Sprintf("~~~%v~~~", v)
}

func Spoiler(v any) string {
	return fmt.Sprintf("||%v||", v)
}

func Ref(v any, url string) string {
	return fmt.Sprintf("[%v](%s)", v, url)
}

func EscapeMarkdown(text string) string {
	return botapi.EscapeText(botapi.ModeMarkdownV2, text)
}
