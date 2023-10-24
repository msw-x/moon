package ulog

type Statistics struct {
	Size     uint
	Trace    uint
	Debug    uint
	Info     uint
	Warning  uint
	Error    uint
	Critical uint
}

func (o *Statistics) Push(level Level, size int) {
	o.Size += uint(size)
	switch level {
	case LevelTrace:
		o.Trace++
	case LevelDebug:
		o.Debug++
	case LevelInfo:
		o.Info++
	case LevelWarning:
		o.Warning++
	case LevelError:
		o.Error++
	case LevelCritical:
		o.Critical++
	}
}
