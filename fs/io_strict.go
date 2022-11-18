package fs

import "github.com/msw-x/moon"

func ReadStrict(path string) []byte {
	raw, err := Read(path)
	moon.Strict(err, "fs read")
	return raw
}

func WriteStrict(path string, content []byte) {
	moon.Strict(Write(path, content), "fs write")
}

func ReadStringStrict(path string) string {
	s, err := ReadString(path)
	moon.Strict(err, "fs read string")
	return s
}

func ReadLinesStrict(path string) []string {
	l, err := ReadLines(path)
	moon.Strict(err, "fs read lines")
	return l
}

func ReadCSVstrict(path string) [][]string {
	r, err := ReadCSV(path)
	moon.Strict(err, "fs read csv")
	return r
}

func WriteCSVstrict(path string, records [][]string) {
	moon.Strict(WriteCSV(path, records), "fs write csv")
}
