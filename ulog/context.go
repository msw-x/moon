package ulog

import (
	"os"
	"sync"
	"time"
)

type context struct {
	inited time.Time
	conf   Conf
	stat   Statistics
	file   *os.File
	maxid  int
	mapid  map[int]bool
	mutex  sync.Mutex
}

var ctx context
