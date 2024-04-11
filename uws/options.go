package uws

import (
	"net/url"
	"time"

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
	ReDialDelay        time.Duration
	HandshakeTimeout   time.Duration
	ReadTimeout        time.Duration
	SendTimeout        time.Duration
	LogSent            LogMessage
	LogRecv            LogMessage
}

func (o *Options) SetDefaultTimeouts() {
	o.PingInterval = time.Second * 24
	o.ReDialInterval = time.Second * 12
	o.ReDialDelay = time.Second * 1
	o.HandshakeTimeout = time.Second * 12
	o.ReadTimeout = o.PingInterval + time.Second*4
	o.SendTimeout = time.Second * 8
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
