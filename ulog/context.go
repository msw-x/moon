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
	fileTime time.Time
	fname    string
	timeLoc  *time.Location
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
	openFile := o.opts.File != opts.File || o.opts.Dir != opts.Dir
	o.opts = opts
	if opts.Timezone == "" {
		o.timeLoc = nil
	} else {
		o.timeLoc, _ = moon.TimezoneLocation(opts.Timezone)
	}
	if openFile {
		o.openFile(false)
	}
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
			text = fmt.Sprintf("%s %v[%s]", text, level.Laconic(), ufmt.Int(count, ufmt.IntCtx{Precision: 0, Dense: true}))
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

func (o *context) now() time.Time {
	now := time.Now()
	if o.timeLoc != nil {
		now = now.In(o.timeLoc)
	}
	return now
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

func (o *context) openFile(prolongation bool) {
	o.fname = ""
	if o.file != nil {
		o.file.Close()
		o.file = nil
	}
	o.fname = o.opts.File
	if o.opts.useDir() {
		appName := o.opts.AppName
		if appName == "" {
			appName = AppName()
		}
		if prolongation {
			appName += ".~"
		}
		o.fname = GenFilename(o.now(), o.opts.Dir, appName)
		o.rotate()
	}
	if o.fname != "" {
		o.file = OpenFile(o.fname, o.opts.Append)
		o.fileSize = 0
		o.fileTime = time.Now()
	}
}

func (o *context) trim(nextMessageSize int) {
	if o.file != nil && o.opts.useDir() && (o.fileSizeExceeded(nextMessageSize) || o.fileTimeExceeded()) {
		o.file.Close()
		o.openFile(true)
	}
}

func (o *context) fileSizeExceeded(nextMessageSize int) bool {
	return o.opts.FileSizeLimit != 0 && (o.fileSize+uint64(nextMessageSize)) > o.opts.FileSizeLimit
}

func (o *context) fileTimeExceeded() bool {
	return o.opts.FileTimeLimit != 0 && time.Since(o.fileTime) > o.opts.FileTimeLimit
}

func (o *context) rotateEnabled() bool {
	return o.opts.DaysCountLimit > 0 ||
		o.opts.TotalSizeLimit > 0
}

func (o *context) rotate() {
	if o.rotateEnabled() {
		dirs := scanDirs(o.opts.Dir)
		if o.opts.DaysCountLimit > 0 && dirs.count() > o.opts.DaysCountLimit {
			n := dirs.count() - o.opts.DaysCountLimit
			dirs.removeByCount(n)
		}
		totalSizeLimit := int64(o.opts.TotalSizeLimit)
		if totalSizeLimit > 0 && dirs.size() > totalSizeLimit {
			n := dirs.size() - totalSizeLimit
			dirs.removeBySize(n)
		}
	}
}
