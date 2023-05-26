package uws

import "github.com/gorilla/websocket"

type Events struct {
	OnPing         func()
	OnMessage      func(int, []byte)
	OnDial         func(string)
	OnDialError    func(error)
	OnConnected    func()
	OnDisconnected func()
	OnCloseError   func(error)
}

func (o *Events) SetOnBinaryMessage(f func([]byte)) {
	o.OnMessage = func(messateType int, data []byte) {
		if messateType == websocket.BinaryMessage {
			f(data)
		}
	}
}

func (o *Events) SetOnTextMessage(f func(string)) {
	o.OnMessage = func(messateType int, data []byte) {
		if messateType == websocket.TextMessage {
			f(string(data))
		}
	}
}

func (o *Events) callOnMessage(messateType int, data []byte) {
	f := o.OnMessage
	if f != nil {
		f(messateType, data)
	}
}

func (o *Events) callOnDial(s string) {
	f := o.OnDial
	if f != nil {
		f(s)
	}
}

func (o *Events) callOnDialError(err error) {
	f := o.OnDialError
	if f != nil {
		f(err)
	}
}

func (o *Events) callOnConnected() {
	f := o.OnConnected
	if f != nil {
		f()
	}
}

func (o *Events) callOnDisconnected() {
	f := o.OnDisconnected
	if f != nil {
		f()
	}
}

func (o *Events) callOnCloseError(err error) {
	f := o.OnCloseError
	if f != nil {
		f(err)
	}
}
