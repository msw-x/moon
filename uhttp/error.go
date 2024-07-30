package uhttp

type OnError func(error, *Response) error

type OnErrors struct {
	InitRequest OnError
	DoRequest   OnError
	ReadBody    OnError

	r *Response
}

func (o *OnErrors) init(r *Response) {
	o.r = r
}

func (o *OnErrors) initRequest(err error) {
	o.call(o.InitRequest, "init request", err)
}

func (o *OnErrors) doRequest(err error) {
	o.call(o.DoRequest, "do request", err)
}

func (o *OnErrors) readBody(err error) {
	o.call(o.ReadBody, "read body", err)
}

func (o *OnErrors) call(on OnError, name string, err error) {
	if on != nil {
		err = on(err, o.r)
	}
	if err != nil {
		o.r.RefineError(name, err)
	}
}
