package color

const (
	Reset      = "\033[0m"
	Bold       = "\033[1m"
	Red        = "\033[31m"
	Green      = "\033[32m"
	Yellow     = "\033[33m"
	Blue       = "\033[34m"
	Purple     = "\033[35m"
	Cyan       = "\033[36m"
	Gray       = "\033[37m"
	White      = "\033[97m"
	BoldRed    = "\033[1;31m"
	BoldGreen  = "\033[1;32m"
	BoldYellow = "\033[1;33m"
	BoldBlue   = "\033[1;34m"
	BoldPurple = "\033[1;35m"
	BoldCyan   = "\033[1;36m"
	BoldGray   = "\033[1;37m"
	BoldWhite  = "\033[1;97m"
)

func Colorize(color, s string) string {
	return color + s + Reset
}

func InBold(s string) string {
	return Colorize(Bold, s)
}

func InRed(s string) string {
	return Colorize(Red, s)
}

func InGreen(s string) string {
	return Colorize(Green, s)
}

func InYellow(s string) string {
	return Colorize(Yellow, s)
}

func InBlue(s string) string {
	return Colorize(Blue, s)
}

func InPurple(s string) string {
	return Colorize(Purple, s)
}

func InCyan(s string) string {
	return Colorize(Cyan, s)
}

func InGray(s string) string {
	return Colorize(Gray, s)
}

func InWhite(s string) string {
	return Colorize(White, s)
}

func InBoldRed(s string) string {
	return Colorize(BoldRed, s)
}

func InBoldGreen(s string) string {
	return Colorize(BoldGreen, s)
}

func InBoldYellow(s string) string {
	return Colorize(BoldYellow, s)
}

func InBoldBlue(s string) string {
	return Colorize(BoldBlue, s)
}

func InBoldPurple(s string) string {
	return Colorize(BoldPurple, s)
}

func InBoldCyan(s string) string {
	return Colorize(BoldCyan, s)
}

func InBoldGray(s string) string {
	return Colorize(BoldGray, s)
}

func InBoldWhite(s string) string {
	return Colorize(BoldWhite, s)
}
