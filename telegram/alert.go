package telegram

import (
	"fmt"
	"strings"
	"time"

	botapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/msw-x/moon/parse"
	"github.com/msw-x/moon/ulog"
	"github.com/msw-x/moon/utime"
)

type AlertBot struct {
	log     *ulog.Log
	bot     *botapi.BotAPI
	token   string
	chatId  int64
	version string
	limiter *Limiter[string]
}

func NewAlertBot(token string, chatId string, version string) *AlertBot {
	o := &AlertBot{
		log:     ulog.New("alert-bot"),
		token:   token,
		version: version,
	}
	if chatId != "" {
		var err error
		o.chatId, err = parse.Int64(chatId)
		if err == nil && token != "" {
			o.bot, err = botapi.NewBotAPI(token)
			if err == nil {
				o.log.Info("authorized on account:", o.bot.Self.UserName)
				o.limiter = NewLimiter(o.send)
			}
		}
		if err != nil {
			o.log.Error(err)
		}
	}
	return o
}

func (o *AlertBot) WithPrePop(prepop func(func() (string, bool)) (string, bool)) *AlertBot {
	if o.limiter != nil {
		o.limiter.WithPrePop(prepop)
	}
	return o
}

func (o *AlertBot) Send(message string) {
	if o.limiter != nil {
		o.limiter.Push(message)
	}
}

func (o *AlertBot) Sendf(f string, v ...any) {
	o.Send(fmt.Sprintf(f, v...))
}

func (o *AlertBot) Startup() {
	o.send("🚀 ***Startup***\n`v" + o.version + "`")
}

func (o *AlertBot) Shutdown(ts time.Duration) {
	if o.limiter != nil {
		o.limiter.Close()
		n := o.limiter.Size()
		q := ""
		if n > 0 {
			q = fmt.Sprintf("\nmessage queue: ***%d***", n)
		}
		s := ""
		if ts > 0 {
			ts = utime.PrettyTruncate(ts)
			s = fmt.Sprintf(" %v", ts)
		}
		o.send("🏁 ***Shutdown***" + s + "\n`" + ulog.Stat() + "`" + q)
	}
}

func (o *AlertBot) SendLog(m ulog.Message) {
	var icon string
	switch m.Level {
	case ulog.LevelWarning:
		icon = "⚠️ "
	case ulog.LevelError:
		icon = "❗"
	case ulog.LevelCritical:
		icon = "💥 "
	default:
		return
	}
	text := strings.ReplaceAll(m.Text, "`", "'")
	o.Send(fmt.Sprintf("%s`%s`", icon, text))
}

func (o *AlertBot) send(text string) {
	if o.bot != nil {
		msg := botapi.NewMessage(o.chatId, text)
		msg.ParseMode = botapi.ModeMarkdown
		_, err := o.bot.Send(msg)
		if err != nil {
			o.log.Tracef("send[%d]: %v; text: %s", o.chatId, err, text)
		}
	}
}
