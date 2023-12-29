package telegram

func Italic(s string) string {
	return "___" + s + "___"
}

func Bold(s string) string {
	return "***" + s + "***"
}

func Monospace(s string) string {
	return "```" + s + "```"
}
