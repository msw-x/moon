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
	return Router{
		router: mux.NewRouter(),
	}.Branch("")
}

func (this Router) Branch(path string) *Router {
	this.path = path
	if !strings.HasSuffix(this.path, "/") {
		this.path += "/"
	}
	return &this
}

func (this *Router) WithID(id any) *Router {
	this.id = fmt.Sprint(id)
	return this
}

func (this *Router) Handle(method string, path string, onRequest OnRequest) {
	if onRequest == nil {
		panic("router on-request func is nil")
	}
	this.init()
	path = this.path + path
	name := RouteName(method, path)
	this.log.Debug(name)
	this.router.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		defer moon.Recover(func(err string) {
			this.log.Error(name, err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err))
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
	name := RouteName("ws", path)
	this.log.Debug(name)
	this.Get(path, func(w http.ResponseWriter, r *http.Request) {
		defer moon.Recover(func(err string) {
			this.log.Error(name, err)
		})
		conn, err := up.Upgrade(w, r, nil)
		moon.Check(err, "upgrade")
		onWebsocket(conn)
	})
}

func (this *Router) Log() *ulog.Log {
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

func RouteName(method, path string) string {
	return fmt.Sprintf("%s:%s", strings.ToLower(method), path)
}
