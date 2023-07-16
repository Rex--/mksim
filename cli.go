package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
)

type CLIArgs struct {
	ProgName string
	InFile   string

	NoGui bool

	Return bool // Print AC before exiting

	// HALT on startup
	HALT bool
	// Exit on HALT
	EXIT bool

	// Clock Speed
	F_CPU int64
}

func printUsage() {
	fmt.Println("Usage:", os.Args[0], "[options] <in_file>")
	fmt.Printf("\nOptions:\n")
	flag.PrintDefaults()
}

func parseArgs() CLIArgs {

	args := CLIArgs{}

	// Set program name and usage print function
	args.ProgName = os.Args[0]
	flag.Usage = printUsage

	// Add flags
	flag.Int64Var(&args.F_CPU, "F_CPU", 8000000, "simulated clock `speed`")

	flag.BoolVar(&args.HALT, "halt", false, "HALT the machine before first instruction cycle")
	flag.BoolVar(&args.EXIT, "exit", false, "Exit the simulator on HALT")

	flag.BoolVar(&args.NoGui, "no-gui", false, "Do not display curses ui")

	flag.BoolVar(&args.Return, "print-return", false, "Print return code (AC) upon exiting")

	help := flag.Bool("help", false, "Print this message and exit")

	// Parse
	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	// Get remaining positional argument (infile)
	if len(flag.Args()) == 1 {
		args.InFile = flag.Arg(0)
	} else {
		flag.Usage()
		os.Exit(1)
	}

	return args
}

type CLIFrontPanel struct {
}

func (fp *CLIFrontPanel) PowerOn(mk MK12) {

}

func (fp *CLIFrontPanel) PowerOff() {

}

func (fp *CLIFrontPanel) Update(mk MK12) {

}

func (fp *CLIFrontPanel) ReadSwitches() int16 {
	return 0
}

type StdinKeyboard struct {
	Stdin   *os.File
	lastKey []byte
}

func NewStdinKeyboard() (sk StdinKeyboard) {
	// disable input buffering
	err := exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
	if err != nil {
		panic(err)
	}
	// do not display entered characters on the screen
	err = exec.Command("stty", "-F", "/dev/tty", "-echo").Run()
	if err != nil {
		panic(err)
	}

	return StdinKeyboard{
		Stdin:   os.Stdin,
		lastKey: make([]byte, 1),
	}
}

func (sk StdinKeyboard) Buffered() (buffered int) {
	buffered, err := sk.Stdin.Read(sk.lastKey)
	if err != nil {
		panic(err)
	}
	return
}

func (sk StdinKeyboard) ReadByte() (char byte, err error) {
	char = sk.lastKey[0]
	return
}
