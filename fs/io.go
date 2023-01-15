package fs

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"io"
	"io/ioutil"
	"os"
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

func WriteString(path string, content string) error {
	return Write(path, []byte(content))
}

func ReadLines(path string) ([]string, error) {
	s, err := ReadString(path)
	return strings.Split(s, "\n"), err
}

func ForEachLine(path string, fn func(string)) error {
	file, err := os.OpenFile(path, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return err
	}
	defer file.Close()
	rd := bufio.NewReader(file)
	for {
		line, err := rd.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		fn(line)
	}
	return err
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
