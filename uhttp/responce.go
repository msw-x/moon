package uhttp

import (
	"encoding/json"
	"fmt"
	"time"
)

type Responce struct {
	Request    Request
	Time       time.Duration
	Status     string
	StatusCode int
	Body       []byte
	Error      error
}

func (o *Responce) Json(v any) error {
	return json.Unmarshal(o.Body, v)
}

func (o *Responce) RefineError(text string, err error) {
	o.Error = fmt.Errorf("%s: %v", text, err)
}

func (o *Responce) Format(f Format) string {
	/*
		"POST[url] 200 OK 350ms 12.3KB"
		if f.Params {
			"?id=1245&user_id=9283577"
		}
		if f.Header {
			"Content-Type: application/json"
		}
		if f.RequestBody {
			"RequestBody"
		}
		if f.ResponceBody {
			"ResponceBody"
		}
	*/
	return ""
}

type Format struct {
	Params       bool
	Header       bool
	RequestBody  bool
	ResponceBody bool
}
