package app

import "os"

func Args() (args []string) {
	if len(os.Args) > 1 {
		for n, v := range os.Args {
			if n > 0 {
				if n == 1 && v == "" {
					continue
				}
				args = append(args, v)
			}
		}
	}
	return
}
