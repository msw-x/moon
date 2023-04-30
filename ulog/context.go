package ulog

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/msw-x/moon"
	"github.com/msw-x/moon/rt"
	"github.com/msw-x/moon/ufmt"
)

type context struct {
	inited   time.Time
	opts     Options
	stat     Statistics
	file     *os.File
	fileSize uint64
	fname    string
	maxid    int
	mapid    map[int]bool
	hook     func(Message)
	mutex    sync.Mutex
}

var ctx context

func (o *context) init(opts Options) {
	o.mutex.Lock()
	defer o.mutex.Unlock()
	opts.init()
	if o.opts.File != opts.File || o.opts.Dir != opts.Dir {
		o.openFile(opts, false)
	}
	o.opts = opts
	o.maxid = 2
	o.mapid = make(map[int]bool)
	o.inited = time.Now()
}

func (o *context) close() {
	o.mutex.Lock()
	defer o.mutex.Unlock()
	if o.file != nil {
		o.file.Close()
		o.file = nil
	}
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
	if o.opts.GoID {
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

func (o *context) query(f Filter) (lines []string, err error) {
	o.mutex.Lock()
	defer o.mutex.Unlock()
	if o.fname == "" {
		err = errors.New("file is not used")
		return
	}
	return QueryFromFile(o.fname, f)
}

func (o *context) goroutineID() int {
	if o.opts.GoID {
		return rt.GoroutineID()
	}
	return -1
}

func (o *context) fmtTime(ts time.Time) string {
	ms := ts.Sub(ts.Truncate(time.Second)).Milliseconds()
	return fmt.Sprintf("%s.%03d", ts.Format("2006-Jan-02 15:04:05"), ms)
}

func (o *context) fmtGoroutineID(id int) string {
	sid := strconv.Itoa(id)
	if len(sid) > o.maxid {
		o.maxid = len(sid)
	}
	for n := o.maxid - len(sid); n != 0; n-- {
		sid = " " + sid
	}
	o.mapid[id] = true
	return sid
}

func (o *context) openFile(opts Options, rotate bool) {
	o.fname = ""
	if o.file != nil {
		o.file.Close()
		o.file = nil
	}
	o.fname = opts.File
	if opts.useDir() {
		appName := opts.AppName
		if appName == "" {
			appName = AppName()
		}
		if rotate {
			appName += ".~"
		}
		o.fname = GenFilename(opts.Dir, appName)
	}
	if o.fname != "" {
		o.file = OpenFile(o.fname, opts.Append)
		o.fileSize = 0
	}
}

func (o *context) rotate(nextMessageSize int) {
	if o.file != nil && o.opts.useDir() && o.opts.FileSizeLimit != 0 && (o.fileSize+uint64(nextMessageSize)) > o.opts.FileSizeLimit {
		o.file.Close()
		o.openFile(o.opts, true)
	}
}
