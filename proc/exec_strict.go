package proc

import "github.com/msw-x/moon/uerr"

func StartStrict(name string, arg ...string) {
	uerr.Strict(Start(name, arg...), "exec(start):", name)
}

func RunStrict(name string, arg ...string) {
	uerr.Strict(Run(name, arg...), "exec(run):", name)
}

func ReadStdoutStrict(name string, arg ...string) string {
	s, err := ReadStdout(name, arg...)
	uerr.Strict(err, "exec(stdout):", name)
	return s
}
