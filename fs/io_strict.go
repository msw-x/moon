package fs

import "github.com/msw-x/moon/uerr"

func ReadStrict(path string) []byte {
	raw, err := Read(path)
	uerr.Strict(err, "fs read")
	return raw
}

func WriteStrict(path string, content []byte) {
	uerr.Strict(Write(path, content), "fs write")
}

func ReadStringStrict(path string) string {
	s, err := ReadString(path)
	uerr.Strict(err, "fs read string")
	return s
}

func ReadLinesStrict(path string) []string {
	l, err := ReadLines(path)
	uerr.Strict(err, "fs read lines")
	return l
}

func ReadCSVstrict(path string) [][]string {
	r, err := ReadCSV(path)
	uerr.Strict(err, "fs read csv")
	return r
}

func WriteCSVstrict(path string, records [][]string) {
	uerr.Strict(WriteCSV(path, records), "fs write csv")
}
