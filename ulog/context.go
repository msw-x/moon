package ulog

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/msw-x/moon"
	"github.com/msw-x/moon/ufmt"
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

func (o *context) init(conf Conf) {
	o.mutex.Lock()
	defer o.mutex.Unlock()
	conf.init()
	if o.conf.File != conf.File || o.conf.Dir != conf.Dir {
		if o.file != nil {
			o.file.Close()
			o.file = nil
		}
		filename := conf.File
		if filename == "" && conf.Dir != "" {
			appName := conf.AppName
			if appName == "" {
				appName = AppName()
			}
			filename = GenFilename(conf.Dir, appName)
		}
		if filename != "" {
			o.file = OpenFile(filename, conf.Append)
		}
	}
	o.conf = conf
	o.maxid = 2
	o.mapid = make(map[int]bool)
	o.inited = time.Now()
}

func (o *context) statistics() string {
	o.mutex.Lock()
	defer o.mutex.Unlock()
	tm := time.Since(o.inited)
	dur := moon.DurationToTime(tm)
	var text string
	if tm < time.Second {
		text = fmt.Sprintf("%d ms", dur.Milliseconds)
	} else {
		text = dur.FormatDays()
	}
	text = fmt.Sprintf("%s | %s", text, ufmt.ByteSize(o.stat.Size))
	if o.conf.GoID {
		text = fmt.Sprintf("%s go[%s]", text, ufmt.WideInt(len(o.mapid)))
	}
	add := func(level Level, count uint) {
		if count > 0 {
			text = fmt.Sprintf("%s %v[%s]", text, level.Laconic(), ufmt.WideInt(count))
		}
	}
	add(LevelTrace, o.stat.Trace)
	add(LevelDebug, o.stat.Debug)
	add(LevelInfo, o.stat.Info)
	add(LevelWarning, o.stat.Warning)
	add(LevelError, o.stat.Error)
	add(LevelCritical, o.stat.Critical)
	return text
}

func (o *context) close() {
	o.mutex.Lock()
	defer o.mutex.Unlock()
	if o.file != nil {
		o.file.Close()
		o.file = nil
	}
}
