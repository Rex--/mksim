package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
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

	// File[path] to use as virtual paper tape
	// TapeFile  string
	// iTapeFile string
	// oTapeFile string

	// Lock memory viewer to page
	Page int
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

	flag.IntVar(&args.Page, "lock", -1, "Lock memory viewer to `page`")

	flag.BoolVar(&args.HALT, "halt", false, "HALT the machine before first instruction cycle")
	flag.BoolVar(&args.EXIT, "exit", false, "Exit the simulator on HALT")

	flag.BoolVar(&args.NoGui, "no-gui", false, "Do not display curses ui")

	flag.BoolVar(&args.Return, "print-return", false, "Print return code (AC) upon exiting")

	// flag.StringVar(&args.TapeFile, "tape", "mk-12.tape", "Specify `path` to file for virtual tape reader/punch")
	// flag.StringVar(&args.iTapeFile, "itape", "", "Specify `path` to file for virtual tape reader")
	// flag.StringVar(&args.oTapeFile, "otape", "", "Specify `path` to file for virtual tape punch")

	help := flag.Bool("help", false, "Print this message and exit")

	// Parse
	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	// if args.iTapeFile == "" {
	// 	args.iTapeFile = args.TapeFile
	// }
	// if args.oTapeFile == "" {
	// 	args.oTapeFile = args.TapeFile
	// }

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

func (fp *CLIFrontPanel) ReadSwitches() uint16 {
	return 0
}

type StdinKeyboard struct {
	Stdin             *os.File
	lastKey           []byte
	originalSttyState bytes.Buffer
}

func (sk *StdinKeyboard) getSttyState(state *bytes.Buffer) (err error) {
	cmd := exec.Command("stty", "-g")
	cmd.Stdin = sk.Stdin
	cmd.Stdout = state
	return cmd.Run()
}

func (sk *StdinKeyboard) saveSttyState() (err error) {
	cmd := exec.Command("stty", "-g")
	cmd.Stdin = sk.Stdin
	cmd.Stdout = &sk.originalSttyState
	return cmd.Run()
}

func (sk *StdinKeyboard) setSttyState(state *bytes.Buffer) (err error) {
	cmd := exec.Command("stty", state.String())
	cmd.Stdin = sk.Stdin
	cmd.Stdout = nil
	return cmd.Run()
}

func NewStdinKeyboard() (sk StdinKeyboard) {
	sk = StdinKeyboard{
		Stdin:   os.Stdin,
		lastKey: make([]byte, 1),
	}
	sk.saveSttyState()
	defer sk.ResetSttyState()

	// disable input buffering
	err := sk.setSttyState(bytes.NewBufferString("cbreak"))
	// err := exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
	if err != nil {
		panic(err)
	}
	// do not display entered characters on the screen
	err = sk.setSttyState(bytes.NewBufferString("-echo"))
	// err = exec.Command("stty", "-F", "/dev/tty", "-echo").Run()
	if err != nil {
		panic(err)
	}

	return sk
}

func (sk StdinKeyboard) Buffered() (buffered int) {
	buffered, err := sk.Stdin.Read(sk.lastKey)
	if err != nil {
		if err == io.EOF {
			// sk.Stdin.Seek(0, 0)
			return 0
		}
		panic(err)
	}
	return
}

func (sk StdinKeyboard) ReadByte() (char byte, err error) {
	char = sk.lastKey[0]
	return
}

func (sk *StdinKeyboard) ResetSttyState() {
	sk.setSttyState(&sk.originalSttyState)
}
