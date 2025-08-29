package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	var pick = flag.Bool("p", false, "pick session from ghq repos")
	var switch_ = flag.Bool("s", false, "switch between existing sessions")
	flag.Parse()

	if !commandExists("ghq") {
		fmt.Fprintf(os.Stderr, "Error: ghq is not installed\n")
		os.Exit(1)
	}

	if !commandExists("tmux") {
		fmt.Fprintf(os.Stderr, "Error: tmux is not installed\n")
		os.Exit(1)
	}

	selections := 0
	if *pick {
		selections++
	}
	if *switch_ {
		selections++
	}

	if selections != 1 {
		fmt.Fprintln(os.Stderr, "Please make exactly one selection (-p or -s)")
		os.Exit(1)
	}

	var err error
	if *pick {
		err = pickSession()
	} else {
		err = switchSession()
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
