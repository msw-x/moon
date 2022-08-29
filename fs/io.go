package fs

import (
	"bytes"
	"encoding/csv"
	"io"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/msw-x/moon"
)

func Read(path string) []byte {
	raw, err := ioutil.ReadFile(path)
	moon.Check(err, "fs read")
	return raw
}

func Write(path string, content []byte) {
	dir := filepath.Dir(path)
	MakeDir(dir)
	err := ioutil.WriteFile(path, content, 0644)
	moon.Check(err, "fs write")
}

func ReadString(path string) string {
	return string(Read(path))
}

func WriteString(path string, content string) {
	Write(path, []byte(content))
}

func ReadLines(path string) []string {
	s := ReadString(path)
	return strings.Split(s, "\n")
}

func ReadCVS(path string) [][]string {
	in := ReadString(path)
	r := csv.NewReader(strings.NewReader(in))
	r.Comma = ';'
	r.Comment = '#'
	records, err := r.ReadAll()
	moon.Check(err, "read cvs")
	return records
}

func WriteCVS(path string, records [][]string) {
	var buf bytes.Buffer
	bw := io.Writer(&buf)
	w := csv.NewWriter(bw)
	w.WriteAll(records)
	moon.Check(w.Error(), "write cvs")
	Write(path, buf.Bytes())
}

func Start(name string, arg ...string) {
	err := exec.Command(name, arg...).Start()
	moon.Check(err, "exec:", name)
}

func Run(name string, arg ...string) {
	err := exec.Command(name, arg...).Run()
	moon.Check(err, "exec:", name)
}

func ReadStdout(name string, arg ...string) string {
	out, err := exec.Command(name, arg...).Output()
	moon.Check(err, "exec:", name)
	s := string(out)
	s = strings.TrimSuffix(s, "\n")
	return s
}
