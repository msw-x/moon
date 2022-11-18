package proc

import (
	"github.com/msw-x/moon"
)

func StartStrict(name string, arg ...string) {
	moon.Strict(Start(name, arg...), "exec(start):", name)
}

func RunStrict(name string, arg ...string) {
	moon.Strict(Run(name, arg...), "exec(run):", name)
}

func ReadStdoutStrict(name string, arg ...string) string {
	s, err := ReadStdout(name, arg...)
	moon.Strict(err, "exec(stdout):", name)
	return s
}
