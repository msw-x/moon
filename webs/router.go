package webs

import (
	"fmt"
	"io/fs"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/msw-x/moon"
	"github.com/msw-x/moon/ufmt"
	"github.com/msw-x/moon/ulog"
)

type Router struct {
	log        *ulog.Log
	id         string
	path       string
	router     *mux.Router
	logRequest bool
}

type OnRequest func(http.ResponseWriter, *http.Request)
type OnWebsocket func(*websocket.Conn)

func NewRouter() (ret *Router) {
	return Router{
		router: mux.NewRouter(),
	}.Branch("")
}

func (this Router) Branch(path string) *Router {
	this.path += path
	if !strings.HasSuffix(this.path, "/") {
		this.path += "/"
	}
	return &this
}

func (this *Router) WithID(id any) *Router {
	this.id = fmt.Sprint(id)
	return this
}

func (this *Router) WithLogRequest(logRequest bool) *Router {
	this.logRequest = logRequest
	return this
}

func (this *Router) Handle(method string, path string, onRequest OnRequest) {
	if onRequest == nil {
		panic("router on-request func is nil")
	}
	this.init()
	uri := this.uri(path)
	this.log.Debug(RouteName(method, uri))
	this.router.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		name := RequestName(r)
		defer moon.Recover(func(err string) {
			this.log.Error(name, err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err))
		})
		r.ParseForm()
		if this.logRequest {
			if r.ContentLength > 0 {
				this.log.Debug(name, ufmt.ByteSize(r.ContentLength))
			} else {
				this.log.Debug(name)
			}
		}
		onRequest(w, r)
	}).Methods(method)
}

func (this *Router) Get(path string, onRequest OnRequest) {
	this.Handle(http.MethodGet, path, onRequest)
}

func (this *Router) Put(path string, onRequest OnRequest) {
	this.Handle(http.MethodPut, path, onRequest)
}

func (this *Router) Post(path string, onRequest OnRequest) {
	this.Handle(http.MethodPost, path, onRequest)
}

func (this *Router) Delete(path string, onRequest OnRequest) {
	this.Handle(http.MethodDelete, path, onRequest)
}

func (this *Router) Files(files fs.FS) {
	this.router.PathPrefix(this.path).Handler(http.FileServer(http.FS(files)))
}

func (this *Router) WebSocket(path string, onWebsocket OnWebsocket) {
	up := websocket.Upgrader{
		ReadBufferSize:  0,
		WriteBufferSize: 0,
	}
	method := http.MethodGet
	this.log.Debug(WebSocketName(RouteName(method, this.uri(path))))
	this.Handle(method, path, func(w http.ResponseWriter, r *http.Request) {
		defer moon.Recover(func(err string) {
			this.log.Error(WebSocketName(RequestName(r)), err)
		})
		conn, err := up.Upgrade(w, r, nil)
		moon.Check(err, "upgrade")
		onWebsocket(conn)
	})
}

func (this *Router) Log() *ulog.Log {
	this.init()
	return this.log
}

func (this *Router) Router() *mux.Router {
	return this.router
}

func (this *Router) init() {
	if this.log == nil {
		this.log = ulog.New("router").WithID(this.id)
	}
}

func (this *Router) uri(path string) string {
	return this.path + path
}

func RouteName(method, uri string) string {
	return ufmt.JoinWith(":", strings.ToUpper(method), uri)
}

func RequestName(r *http.Request) string {
	if r == nil {
		return "?"
	}
	return ufmt.JoinWith("?", r.RemoteAddr, RouteName(r.Method, r.RequestURI))
}

func WebSocketName(name string) string {
	return ufmt.JoinWith("@", name, "ws")
}
