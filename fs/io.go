package fs

import (
	"bytes"
	"encoding/csv"
	"io"
	"io/ioutil"
	"path/filepath"
	"strings"
)

func Read(path string) ([]byte, error) {
	return ioutil.ReadFile(path)
}

func Write(path string, content []byte) error {
	dir := filepath.Dir(path)
	err := MakeDir(dir)
	if err == nil {
		err = ioutil.WriteFile(path, content, 0644)
	}
	return err
}

func ReadString(path string) (string, error) {
	bytes, err := Read(path)
	return string(bytes), err
}

func WriteString(path string, content string) {
	Write(path, []byte(content))
}

func ReadLines(path string) ([]string, error) {
	s, err := ReadString(path)
	return strings.Split(s, "\n"), err
}

func ReadCSV(path string) (records [][]string, err error) {
	in, err := ReadString(path)
	if err == nil {
		r := csv.NewReader(strings.NewReader(in))
		r.Comma = ';'
		r.Comment = '#'
		return r.ReadAll()
	}
	return
}

func WriteCSV(path string, records [][]string) error {
	var buf bytes.Buffer
	bw := io.Writer(&buf)
	w := csv.NewWriter(bw)
	w.Comma = ';'
	err := w.WriteAll(records)
	if err != nil {
		return err
	}
	return Write(path, buf.Bytes())
}
