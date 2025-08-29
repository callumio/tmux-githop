package main

import (
	"flag"
	"fmt"
	"os"
)

var assumeYes bool

func main() {
	var pick = flag.Bool("p", false, "pick session from ghq repos")
	var switch_ = flag.Bool("s", false, "switch between existing sessions")
	flag.BoolVar(&assumeYes, "y", false, "assume yes for all prompts")
	flag.Parse()

	if !commandExists("ghq") {
		fmt.Fprintf(os.Stderr, "Error: 'ghq' command not found in PATH. Please install GHQ (https://github.com/x-motemen/ghq) and ensure it's in your PATH\n")
		os.Exit(1)
	}

	if !commandExists("tmux") {
		fmt.Fprintf(os.Stderr, "Error: 'tmux' command not found in PATH. Please install tmux (https://github.com/tmux/tmux) and ensure it's in your PATH\n")
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
		fmt.Fprintln(os.Stderr, "Error: Please specify exactly one flag: -p (pick from GHQ repos) or -s (switch existing sessions)")
		fmt.Fprintln(os.Stderr, "Optional: -y to assume yes for all prompts")
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
