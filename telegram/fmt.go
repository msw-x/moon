package telegram

import (
	"fmt"
	"strings"
)

func Spoiler(v any) string {
	return fmt.Sprintf("||%v||", v)
}

func Ref(v any, url string) string {
	return fmt.Sprintf("[%v](%s)", v, url)
}

func UserRef(v any, user string) string {
	return Ref(v, "https://t.me/"+user)
}

func UserIdRef(v any, id int64) string {
	return Ref(v, "tg://user?id="+fmt.Sprint(id))
}

func EscapeMarkdownV2(v any) string {
	// Fixed EscapeMarkdown for \ symbol #44
	// https://github.com/go-telegram/bot/pull/44
	// return botapi.EscapeText(botapi.ModeMarkdownV2, fmt.Sprint(v))
	// EscapeMarkdown escapes special symbols for Telegram MarkdownV2 syntax
	// https://github.com/go-telegram/bot/blob/main/common.go
	// https://core.telegram.org/bots/api
	// In all other places characters '_', '*', '[', ']', '(', ')', '~', '`', '>', '#', '+', '-', '=', '|', '{', '}', '.', '!' must be escaped with the preceding character '\'.
	var s = fmt.Sprint(v)
	var shouldBeEscaped = "_*[]()~`>#+-=|{}.!\\"
	var result []rune
	for _, r := range s {
		if strings.ContainsRune(shouldBeEscaped, r) {
			result = append(result, '\\')
		}
		result = append(result, r)
	}
	return string(result)
}
