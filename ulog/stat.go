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

func (this *Statistics) Push(level Level, size int) {
	this.Size += uint(size)
	switch level {
	case LevelTrace:
		this.Trace++
	case LevelDebug:
		this.Debug++
	case LevelInfo:
		this.Info++
	case LevelWarning:
		this.Warning++
	case LevelError:
		this.Error++
	case LevelCritical:
		this.Critical++
	}
}
