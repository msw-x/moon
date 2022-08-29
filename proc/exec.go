package proc

import (
	"os/exec"
	"strings"

	"github.com/msw-x/moon"
)

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
