package uws

import (
	"net/url"
	"time"

	"github.com/gorilla/websocket"
	"github.com/msw-x/moon/uerr"
	"github.com/msw-x/moon/uhttp"
)

type Options struct {
	Base               string
	Path               string
	Proxy              *url.URL
	InsecureSkipVerify bool
	PingInterval       time.Duration
	ReDialInterval     time.Duration
	HandshakeTimeout   time.Duration
	ReadTimeout        time.Duration
	SendTimeout        time.Duration
	OnPing             func()
	OnMessage          func(int, []byte)
	OnDial             func(string)
	OnDialError        func(error)
	OnConnected        func()
	OnDisconnected     func()
	OnCloseError       func(error)
	LogSentType        bool
	LogSentSize        bool
	LogSentData        bool
	LogReadType        bool
	LogReadSize        bool
	LogReadData        bool
}

func (o *Options) Default() {
	o.PingInterval = time.Second * 24
	o.ReDialInterval = time.Second * 10
	o.HandshakeTimeout = time.Second * 4
	o.ReadTimeout = o.PingInterval + time.Second*4
	o.SendTimeout = time.Second * 4
}

func (o *Options) Url() string {
	return uhttp.UrlJoin(o.Base, o.Path)
}

func (o *Options) EmptyUrl() bool {
	return o.Url() == ""
}

func (o *Options) AppendPath(path string) {
	o.Path = uhttp.UrlJoin(o.Path, path)
}

func (o *Options) SetProxy(proxy string) {
	if proxy == "" {
		o.Proxy = nil
	}
	var err error
	o.Proxy, err = url.Parse(proxy)
	uerr.Strictf(err, "parse proxy url: %s", o.Proxy)
}

func (o *Options) SetOnBinaryMessage(f func([]byte)) {
	o.OnMessage = func(messateType int, data []byte) {
		if messateType == websocket.BinaryMessage {
			f(data)
		}
	}
}

func (o *Options) SetOnTextMessage(f func(string)) {
	o.OnMessage = func(messateType int, data []byte) {
		if messateType == websocket.TextMessage {
			f(string(data))
		}
	}
}

func (o *Options) callOnMessage(messateType int, data []byte) {
	f := o.OnMessage
	if f != nil {
		f(messateType, data)
	}
}

func (o *Options) callOnDial(s string) {
	f := o.OnDial
	if f != nil {
		f(s)
	}
}

func (o *Options) callOnDialError(err error) {
	f := o.OnDialError
	if f != nil {
		f(err)
	}
}

func (o *Options) callOnConnected() {
	f := o.OnConnected
	if f != nil {
		f()
	}
}

func (o *Options) callOnDisconnected() {
	f := o.OnDisconnected
	if f != nil {
		f()
	}
}

func (o *Options) callOnCloseError(err error) {
	f := o.OnCloseError
	if f != nil {
		f(err)
	}
}
