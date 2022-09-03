package webs

import (
	"fmt"
	"io/fs"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/msw-x/moon"
	"github.com/msw-x/moon/ulog"
)

type Router struct {
	log    *ulog.Log
	id     string
	path   string
	router *mux.Router
}

type OnRequest func(http.ResponseWriter, *http.Request)
type OnWebsocket func(*websocket.Conn)

func NewRouter() (ret *Router) {
	ret = &Router{
		router: mux.NewRouter(),
	}
	ret.WithUrl("")
	return
}

func (this *Router) WithUrl(path string) *Router {
	this.path = path
	if !strings.HasSuffix(this.path, "/") {
		this.path += "/"
	}
	return this
}

func (this *Router) WithID(id any) *Router {
	this.id = fmt.Sprint(id)
	return this
}

func (this *Router) Handle(method string, path string, onRequest OnRequest) {
	this.init()
	path = this.path + path
	this.router.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		defer moon.Recover(func(err string) {
			this.log.Errorf("%s:%s %s", method, path, err)
		})
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
	this.Get(path, func(w http.ResponseWriter, r *http.Request) {
		defer moon.Recover(func(err string) {
			this.log.Errorf("ws:%s %s", path, err)
		})
		conn, err := up.Upgrade(w, r, nil)
		moon.Check(err, "upgrade")
		onWebsocket(conn)
	})
}

func (this *Router) init() {
	if this.log == nil {
		this.log = ulog.New("router").WithID(this.id)
	}
}
