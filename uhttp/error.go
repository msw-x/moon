package uhttp

type OnError func(error) error
type OnRefineError func(string, error)

type OnErrors struct {
	InitRequest OnError
	DoRequest   OnError
	ReadBody    OnError
}

func (o *OnErrors) initRequest(err error, f OnRefineError) {
	o.call(o.InitRequest, "init request", err, f)
}

func (o *OnErrors) doRequest(err error, f OnRefineError) {
	o.call(o.DoRequest, "do request", err, f)
}

func (o *OnErrors) readBody(err error, f OnRefineError) {
	o.call(o.ReadBody, "read body", err, f)
}

func (o *OnErrors) call(on OnError, name string, err error, f OnRefineError) {
	if on != nil {
		err = on(err)
	}
	if err != nil {
		f(name, err)
	}
}
