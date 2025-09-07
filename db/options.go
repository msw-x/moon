package db

import (
	"runtime"
	"time"
)

type Options struct {
	User            string
	Pass            string
	Host            string
	Name            string
	Timeout         time.Duration
	MaxConnFactor   float32
	MinOpenConns    int
	Strict          bool
	Insecure        bool
	DisablePrepared bool
	LogErrors       bool
	LogQueries      bool
	LogLongQueries  bool
	WarnLongQueries bool
	LongQueriesTime time.Duration
	ReadOnly        bool
}

func (o Options) MaxOpenConnections() int {
	v := int(o.MaxConnFactor * float32(runtime.GOMAXPROCS(0)))
	if v == 0 {
		v = 1
	}
	if v < o.MinOpenConns {
		v = o.MinOpenConns
	}
	return v
}
