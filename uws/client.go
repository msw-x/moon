package uws

import (
	"crypto/tls"
	"errors"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/msw-x/moon/app"
	"github.com/msw-x/moon/ufmt"
	"github.com/msw-x/moon/ujson"
	"github.com/msw-x/moon/ulog"
)

type Client struct {
	Options Options

	log      *ulog.Log
	mutex    sync.Mutex
	ws       *websocket.Conn
	job      *app.Job
	lastSend time.Time
}

func NewClient(url string) *Client {
	o := new(Client)
	o.Options.Default()
	o.Options.Base = url
	o.log = ulog.Empty()
	o.job = app.NewJob()
	return o
}

func (o *Client) Close() {
	o.job.Cancel()
	o.closeSocket()
	o.job.Stop()
}

func (o *Client) WithLog(log *ulog.Log) *Client {
	o.log = log
	o.job.WithLog(log)
	return o
}

func (o *Client) WithBase(base string) *Client {
	o.Options.Base = base
	return o
}

func (o *Client) WithPath(path string) *Client {
	o.Options.Path = path
	return o
}

func (o *Client) WithAppendPath(path string) *Client {
	o.Options.AppendPath(path)
	return o
}

func (o *Client) WithProxy(proxy string) *Client {
	o.Options.SetProxy(proxy)
	return o
}

func (o *Client) WithProxyUrl(url *url.URL) *Client {
	o.Options.Proxy = url
	return o
}

func (o *Client) WithOnMessage(f func(int, []byte)) {
	o.Options.OnMessage = f
}

func (o *Client) WithOnBinaryMessage(f func([]byte)) {
	o.Options.SetOnBinaryMessage(f)
}

func (o *Client) WithOnTextMessage(f func(string)) {
	o.Options.SetOnTextMessage(f)
}

func (o *Client) WithOnDial(f func(string)) {
	o.Options.OnDial = f
}

func (o *Client) WithOnDialError(f func(error)) {
	o.Options.OnDialError = f
}

func (o *Client) WithOnConnected(f func()) {
	o.Options.OnConnected = f
}

func (o *Client) WithOnDisconnected(f func()) {
	o.Options.OnDisconnected = f
}

func (o *Client) WithOnCloseError(f func(error)) {
	o.Options.OnCloseError = f
}

func (o *Client) Run() {
	o.job.OnStart(func() {
		if o.Options.EmptyUrl() {
			o.job.Cancel()
		}
	})
	o.job.RunLoop(o.connectAndRun)
}

func (o *Client) Send(messageType int, data []byte) error {
	o.mutex.Lock()
	defer o.mutex.Unlock()
	ws, err := o.socket()
	if err != nil {
		return err
	}
	o.lastSend = time.Now()
	if o.Options.SendTimeout != 0 {
		ws.SetWriteDeadline(time.Now().Add(o.Options.SendTimeout))
	}
	err = ws.WriteMessage(messageType, data)
	if err == nil {
		o.dump("sent", messageType, data, o.Options.LogReadSize, o.Options.LogReadData)
	} else {
		o.log.Error(err)
		o.closeSocket()
	}
	return err
}

func (o *Client) SendBinary(data []byte) error {
	err := o.Send(websocket.BinaryMessage, data)
	return err
}

func (o *Client) SendText(data string) error {
	err := o.Send(websocket.TextMessage, []byte(data))
	return err
}

func (o *Client) SendJson(v any) (err error) {
	var data []byte
	data, err = ujson.MarshalLowerCase(v)
	if err == nil {
		return o.Send(websocket.TextMessage, data)
	}
	return
}

func (o *Client) socket() (ws *websocket.Conn, err error) {
	ws = o.ws
	if ws == nil {
		err = errors.New("empty socket")
	}
	return
}

func (o *Client) closeSocket() {
	if ws, err := o.socket(); err == nil {
		ws.Close()
	}
}

func (o *Client) connectAndRun() {
	url := o.Options.Url()
	o.log.Info("dial:", url)
	err := o.dial(url)
	if err != nil {
		o.log.Error("dial:", err)
		o.Options.callOnDialError(err)
		o.job.Sleep(o.Options.ReDialInterval)
		return
	}
	defer o.onDisconnected()
	o.onConnected()
	o.run()
}

func (o *Client) dial(url string) (err error) {
	if o.Options.OnDial != nil {
		o.Options.OnDial(url)
	}
	dialer := websocket.Dialer{
		HandshakeTimeout: o.Options.HandshakeTimeout,
	}
	if strings.HasPrefix(url, "wss") {
		dialer.TLSClientConfig = new(tls.Config)
		dialer.TLSClientConfig.InsecureSkipVerify = o.Options.InsecureSkipVerify
	}
	if o.Options.Proxy != nil {
		dialer.Proxy = http.ProxyURL(o.Options.Proxy)
	}
	o.ws, _, err = dialer.Dial(url, nil)
	return
}

func (o *Client) onConnected() {
	o.log.Info("connected")
	o.Options.callOnConnected()
}

func (o *Client) onDisconnected() {
	o.log.Info("disconnected")
	defer o.closeSocket()
	defer o.Options.callOnDisconnected()
}

func (o *Client) run() {
	if o.Options.PingInterval != 0 {
		pingJob := app.NewJob()
		pingJob.RunTicks(o.ping, o.Options.PingInterval)
		defer pingJob.Stop()
	}
	for o.job.Do() {
		if o.Options.ReadTimeout != 0 {
			o.ws.SetReadDeadline(time.Now().Add(o.Options.ReadTimeout))
		}
		messageType, data, err := o.ws.ReadMessage()
		if o.job.Do() {
			if err != nil {
				o.closeSocket()
				if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway, websocket.CloseNoStatusReceived) {
					o.log.Info("recv:", err)
					o.Options.callOnCloseError(err)
				} else {
					o.log.Error("recv:", err)
				}
				return
			}
			if o.job.Do() {
				o.dump("recv", messageType, data, o.Options.LogReadSize, o.Options.LogReadData)
				o.Options.callOnMessage(messageType, data)
			}
		}
	}
}

func (o *Client) ping() {
	if o.lastSend.IsZero() {
		o.lastSend = time.Now()
	} else {
		if time.Since(o.lastSend) > o.Options.PingInterval/2 {
			if o.Options.OnPing == nil {
				o.Send(websocket.PingMessage, nil)
			} else {
				o.Options.OnPing()
			}
		}
	}
}

func (o *Client) dump(action string, messageType int, data []byte, viewSize bool, viewData bool) {
	if viewSize || viewData {
		var s string
		switch messageType {
		case websocket.TextMessage:
			s = "text"
		case websocket.BinaryMessage:
			s = "binary"
		case websocket.CloseMessage:
			s = "close"
		case websocket.PingMessage:
			s = "ping"
		case websocket.PongMessage:
			s = "pong"
		}
		size := len(data)
		if size > 0 {
			if viewSize {
				s += ": "
				s += ufmt.ByteSizeDense(len(data))
			}
			if viewData {
				s += ": "
				if messageType == websocket.TextMessage {
					s += string(data)
				} else {
					s += ufmt.Hex(data)
				}
			}
		}
		o.log.Debugf("%s: %s", action, s)
	}
}
