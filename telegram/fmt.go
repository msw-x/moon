package telegram

import (
	"fmt"

	botapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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
	return botapi.EscapeText(botapi.ModeMarkdownV2, fmt.Sprint(v))
}
