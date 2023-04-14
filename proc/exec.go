package proc

import (
	"os/exec"
	"strings"
)

func Start(name string, arg ...string) error {
	return exec.Command(name, arg...).Start()
}

func Run(name string, arg ...string) error {
	return exec.Command(name, arg...).Run()
}

func ReadStdout(name string, arg ...string) (s string, err error) {
	var out []byte
	out, err = exec.Command(name, arg...).Output()
	s = string(out)
	s = strings.TrimSuffix(s, "\n")
	return
}
