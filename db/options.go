package db

import "time"

type Options struct {
	User          string
	Pass          string
	Host          string
	Name          string
	Timeout       time.Duration
	MaxConnFactor float32
	MinOpenConns  int
	Strict        bool
	Insecure      bool
	LogErrors     bool
	LogQueries    bool
}
