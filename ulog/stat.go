package ulog

type Statistics struct {
	Debug    uint
	Info     uint
	Warning  uint
	Error    uint
	Critical uint
}

func (this *Statistics) Push(level Level) {
	switch level {
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
