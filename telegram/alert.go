package telegram

import (
	"fmt"
	"strings"
	"time"

	botapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/msw-x/moon/parse"
	"github.com/msw-x/moon/ufmt"
	"github.com/msw-x/moon/ulog"
	"github.com/msw-x/moon/utime"
)

type AlertBot struct {
	log         *ulog.Log
	bot         *botapi.BotAPI
	chatId      int64
	limiter     LimiterIf[string]
	closeAt     time.Time
	logMsgLimit int
}

func NewAlertBot(token string, chatId string) *AlertBot {
	o := &AlertBot{
		log:         ulog.New("alert-bot"),
		logMsgLimit: 1024,
	}
	if chatId != "" {
		var err error
		o.chatId, err = parse.Int64(chatId)
		if err == nil && token != "" {
			o.bot, err = botapi.NewBotAPI(token)
			if err == nil {
				o.log.Info("authorized on account:", o.bot.Self.UserName)
				o.limiter = NewLimiter[string](o.send)
			}
		}
		if err != nil {
			o.log.Error(err)
		}
	}
	return o
}

func (o *AlertBot) Close() {
	var t time.Duration
	if !o.closeAt.IsZero() {
		t = time.Since(o.closeAt)
	}
	o.Shutdown(t)
}

func (o *AlertBot) Closing() {
	o.closeAt = time.Now()
}

func (o *AlertBot) SetLogMessageLimit(limit int) {
	o.logMsgLimit = limit
}

func (o *AlertBot) SetQueue(queue QueueIf[string]) {
	if o.limiter != nil {
		o.limiter.SetQueue(queue)
	}
}

func (o *AlertBot) QueueSize() int {
	if o.limiter == nil {
		return 0
	}
	return o.limiter.Queue().Size()
}

func (o *AlertBot) QueueEmpty() bool {
	return o.QueueSize() == 0
}

func (o *AlertBot) Send(message string) {
	if o.limiter != nil {
		if !o.limiter.Push(message) {
			o.log.Tracef("limiter overloaded, the message is lost")
		}
	}
}

func (o *AlertBot) Sendf(f string, v ...any) {
	o.Send(fmt.Sprintf(f, v...))
}

func (o *AlertBot) Startup(version string) {
	v := version
	if v != "" {
		v = fmt.Sprintf("\n`v%s`", v)
	}
	o.send("🚀 *Startup*" + v)
}

func (o *AlertBot) Shutdown(ts time.Duration) {
	if o.limiter != nil {
		o.limiter.Close()
		s := ""
		if ts > 0 {
			ts = utime.PrettyTruncate(ts)
			s = fmt.Sprintf(" %v", ts)
			s = EscapeMarkdownV2(s)
		}
		o.predictiveSend("🏁 *Shutdown*"+s+"\n`"+ulog.Stat()+"`", 0)
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
	tail := ""
	text := m.Text
	limit := o.logMsgLimit
	n := len(text)
	if n > limit {
		text = text[0:limit]
		tail = fmt.Sprintf("\nmessage length limit exceeded: *%s / %s*",
			EscapeMarkdownV2(ufmt.ByteSize(n)),
			EscapeMarkdownV2(ufmt.ByteSize(limit)))
	}
	text = strings.ReplaceAll(text, "`", "'")
	o.Send(fmt.Sprintf("%s`%s`%s", icon, text, tail))
}

func (o *AlertBot) Discard(n int) int {
	return o.limiter.Discard(n)
}

func (o *AlertBot) pureSend(text string) {
	if o.bot != nil {
		msg := botapi.NewMessage(o.chatId, text)
		msg.ParseMode = botapi.ModeMarkdownV2
		_, err := o.bot.Send(msg)
		if err != nil {
			o.log.Errorf("send[%d]: %v", o.chatId, err)
			o.log.Tracef("send[%d]: %v; text: %s", o.chatId, err, text)
		}
	}
}

func (o *AlertBot) predictiveSend(text string, queueSizeLimit int) {
	if o.limiter != nil {
		n := o.limiter.Queue().Size()
		if n > queueSizeLimit {
			o.pureSend(text + fmt.Sprintf("\n📩 *%d*", n))
		} else {
			o.pureSend(text)
		}
	}
}

func (o *AlertBot) send(text string) {
	o.predictiveSend(text, 10)
}
