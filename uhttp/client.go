package uhttp

import (
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	c     *http.Client
	url   string
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

func (o *Client) WithUrl(url string) *Client {
	o.url = url
	return o
}

func (o *Client) WithPath(path string) *Client {
	return o.Clone().WithUrl(urlJoin(o.url, path))
}

func (o *Client) WithProxy(url *url.URL) *Client {
	o.c.Transport = &http.Transport{
		Proxy: http.ProxyURL(url),
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

func (o *Client) Timeout() time.Duration {
	return o.c.Timeout
}

func (o *Client) Request(method string, url string) *Performer {
	return &Performer{
		c: o.c,
		r: Request{
			Method: method,
			Url:    urlJoin(o.url, url),
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
