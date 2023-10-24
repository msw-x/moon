package uhttp

type Format struct {
	RequestParams     bool
	RequestHeader     bool
	ResponceHeader    bool
	RequestBody       bool
	ResponceBody      bool
	RequestBodyTrim   bool
	ResponceBodyTrim  bool
	RequestBodyLimit  int
	ResponceBodyLimit int
}
