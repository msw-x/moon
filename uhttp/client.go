package uhttp

import (
	"net/http"
	"net/url"
	"time"

	"github.com/msw-x/moon/uerr"
)

type Client struct {
	c      *http.Client
	base   string
	path   string
	trace  func(Response)
	errors OnErrors
}

func NewClient() *Client {
	o := new(Client)
	o.c = new(http.Client)
	return o
}

/// Close() - завершение выполнения всех запросов ???

func (o *Client) Copy() *Client {
	c := *o
	return &c
}

func (o *Client) Transport() http.RoundTripper {
	if o.c.Transport == nil {
		return nil
	}
	return o.c.Transport
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
	return o.WithPath(UrlJoin(o.path, path))
}

func (o *Client) WithTransport(transport http.RoundTripper) *Client {
	o.c.Transport = transport
	return o
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
	t := http.DefaultTransport.(*http.Transport).Clone()
	if url != nil {
		t.Proxy = http.ProxyURL(url)
	}
	o.WithTransport(t)
	return o
}

func (o *Client) WithTimeout(timeout time.Duration) *Client {
	o.c.Timeout = timeout
	return o
}

func (o *Client) WithTrace(trace func(Response)) *Client {
	o.trace = trace
	return o
}

func (o *Client) Timeout() time.Duration {
	return o.c.Timeout
}

func (o *Client) WithOnInitRequestError(f OnError) *Client {
	o.errors.InitRequest = f
	return o
}

func (o *Client) WithOnDoRequestError(f OnError) *Client {
	o.errors.DoRequest = f
	return o
}

func (o *Client) WithOnReadBodyError(f OnError) *Client {
	o.errors.ReadBody = f
	return o
}

func (o *Client) Base() string {
	return o.base
}

func (o *Client) Url(url string) string {
	return UrlJoin(o.base, o.path, url)
}

func (o *Client) Request(method string, url string) *Performer {
	return &Performer{
		Request: Request{
			Method: method,
			Url:    o.Url(url),
		},
		c:      o.c,
		trace:  o.trace,
		errors: o.errors,
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
