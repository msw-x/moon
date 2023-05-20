package uhttp

import (
	"net/http"
	"net/url"
	"time"

	"github.com/msw-x/moon/uerr"
	"github.com/msw-x/moon/ulog"
)

type Client struct {
	c     *http.Client
	base  string
	path  string
	trace func(Responce)
}

func NewClient() *Client {
	o := new(Client)
	o.c = new(http.Client)
	return o
}

func (o *Client) Clone() *Client {
	c := *o
	return &c
}

func (o *Client) WithBase(base string) *Client {
	o.base = base
	return o
}

func (o *Client) WithPath(path string) *Client {
	o.path = path
	return o
}

func (o *Client) WithAppendPath(path string) *Client {
	return o.WithPath(urlJoin(o.path, path))
}

func (o *Client) WithProxy(proxy string) *Client {
	if proxy == "" {
		return o.WithProxyUrl(nil)
	}
	proxyUrl, err := url.Parse(proxy)
	uerr.Strictf(err, "parse proxy url: %s", proxy)
	return o.WithProxyUrl(proxyUrl)
}

func (o *Client) WithProxyUrl(url *url.URL) *Client {
	if url == nil {
		o.c.Transport = nil
	} else {
		o.c.Transport = &http.Transport{
			Proxy: http.ProxyURL(url),
		}
	}
	return o
}

func (o *Client) WithTimeout(timeout time.Duration) *Client {
	o.c.Timeout = timeout
	return o
}

func (o *Client) WithTrace(trace func(Responce)) *Client {
	o.trace = trace
	return o
}

func (o *Client) WithTraceFormat(log *ulog.Log, f Format) *Client {
	return o.WithTrace(TraceFormat(log, f))
}

func (o *Client) Timeout() time.Duration {
	return o.c.Timeout
}

func (o *Client) Url(url string) string {
	return urlJoin(o.base, o.path, url)
}

func (o *Client) Request(method string, url string) *Performer {
	return &Performer{
		c: o.c,
		r: Request{
			Method: method,
			Url:    o.Url(url),
		},
		t: o.trace,
	}
}

func (o *Client) Get(url string) *Performer {
	return o.Request(http.MethodGet, url)
}

func (o *Client) Post(url string) *Performer {
	return o.Request(http.MethodPost, url)
}

func (o *Client) Put(url string) *Performer {
	return o.Request(http.MethodPut, url)
}

func (o *Client) Delete(url string) *Performer {
	return o.Request(http.MethodDelete, url)
}
