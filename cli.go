package main

import (
	"flag"
	"os"
)

type CLIArgs struct {
	ProgName string
	InFile   string

	// HALT on startup
	HALT bool
	// Exit on HALT
	EXIT bool

	// Clock Speed
	F_CPU int64
}

func parseArgs() CLIArgs {

	args := CLIArgs{}

	// Set program name
	args.ProgName = os.Args[0]

	// Add flags
	flag.Int64Var(&args.F_CPU, "F_CPU", 8000000, "simulated clock speed in Hz")

	flag.BoolVar(&args.HALT, "halt", false, "HALT the machine before first instruction cycle")
	flag.BoolVar(&args.EXIT, "exit", false, "Exit the simulator on HALT")

	// Parse
	flag.Parse()

	// Get remaining positional argument (infile)
	if len(flag.Args()) == 1 {
		args.InFile = flag.Arg(0)
	} else {
		flag.Usage()
		os.Exit(1)
	}

	return args
}
