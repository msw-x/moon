package uhttp

import (
	"bytes"
	"io"
	"net/http"
)

func DumpBody(r *http.Request) (v []byte, err error) {
	if r.Body != nil {
		v, r.Body, err = drainBody(r.Body)
	}
	return
}

func drainBody(b io.ReadCloser) ([]byte, io.ReadCloser, error) {
	if b == nil || b == http.NoBody {
		return nil, http.NoBody, nil
	}
	var err error
	var buf bytes.Buffer
	if _, err = buf.ReadFrom(b); err != nil {
		return nil, b, err
	}
	if err = b.Close(); err != nil {
		return nil, b, err
	}
	return buf.Bytes(), io.NopCloser(bytes.NewReader(buf.Bytes())), nil
}
