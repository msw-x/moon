package app

import (
	"time"

	"github.com/msw-x/moon/ulog"
	"github.com/msw-x/moon/usync"
)

type Job struct {
	do       *usync.Do
	log      *ulog.Log
	level    ulog.Level
	onInit   func() error
	onStart  func()
	onFinish func()
}

func NewJob() *Job {
	o := new(Job)
	o.do = usync.NewDo()
	o.log = ulog.Empty()
	o.level = ulog.LevelDebug
	return o
}

func (o *Job) Do() bool {
	return o.do.Do()
}

func (o *Job) Wait() {
	o.logPrint("wait")
	for o.Do() {
		o.Sleep(time.Millisecond)
	}
	o.logPrint("waited")
}

func (o *Job) Stop() {
	o.logPrint("stop")
	o.do.Stop()
	o.logPrint("stopped")
}

func (o *Job) Cancel() {
	o.do.Cancel()
}

func (o *Job) Sleep(timeout time.Duration) {
	o.do.Sleep(timeout)
}

func (o *Job) OnInit(fn func() error) *Job {
	o.onInit = fn
	return o
}

func (o *Job) OnStart(fn func()) *Job {
	o.onStart = fn
	return o
}

func (o *Job) OnFinish(fn func()) *Job {
	o.onFinish = fn
	return o
}

func (o *Job) WithLog(log *ulog.Log) *Job {
	o.log = log
	return o
}

func (o *Job) WithLogLevel(level ulog.Level) *Job {
	o.level = level
	return o
}

func (o *Job) Run(fn func()) {
	go func() {
		defer o.recover()
		defer o.do.Notify()
		if o.init() {
			o.start()
			defer o.finish()
			if o.Do() {
				fn()
			}
		}
	}()
}

func (o *Job) RunLoop(fn func()) {
	o.Run(func() {
		for o.Do() {
			fn()
		}
	})
}

func (o *Job) RunTicks(fn func(), interval time.Duration) {
	o.RunLoop(func() {
		fn()
		o.Sleep(interval)
	})
}

func (o *Job) logPrint(v ...any) {
	o.log.Print(o.level, v...)
}

func (o *Job) recover() {
	if r := recover(); r != nil {
		if o.log.Enabled() {
			o.log.Critical(r)
		} else {
			ulog.Critical(r)
		}
	}
}

func (o *Job) init() bool {
	if o.onInit != nil {
		o.logPrint("init")
		err := o.onInit()
		if err != nil {
			o.log.Error(err)
			return false
		}
		o.logPrint("inited")
	}
	return true
}

func (o *Job) start() {
	if o.onStart != nil {
		o.logPrint("start")
		o.onStart()
	}
	o.logPrint("started")
}

func (o *Job) finish() {
	if o.onFinish != nil {
		o.logPrint("finishing")
		o.onFinish()
	}
	o.logPrint("finished")
}
