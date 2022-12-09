package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jroimartin/gocui"
)

// Instructions
const (
	AND = 0o0
	TAD = 0o1
	ISZ = 0o2
	DCA = 0o3
	JMS = 0o4
	JMP = 0o5
	IOT = 0o6
	OPR = 0o7
)

// OPR Instruction Groups
const (
	OPR_GROUP_1 = 0
	OPR_GROUP_2 = 1
	OPR_GROUP_3 = 2
)

// Special memory addresses
const (
	INT_vect   = 0o0
	RESET_vect = 0o200

	AUTO_begin = 0o10
	AUTO_end   = 0o17
)

// This structure contains the various components of a theoretical MK-12
// All registers are stored as int16 but have a valid range of -/+4096 (12-bit signed int)
type MK12 struct {
	// Program Counter
	PC int16

	// Instruction register
	IR int16

	// Accumulator Register
	AC int16

	// Link Flag [1-bit]
	L bool

	// Memory Address Register
	MA int16

	// Memory Buffer Register
	MB int16

	// Memory [4K x 12 (int16)]
	// Addresses 0o0 to 0o7777
	MEM [4096]int16

	// Switch Register
	// (unused)
	SR int16

	// The state structure holds the current state of the CPU
	STATE struct {
		// If halt is set, the computer is halted during the fetch phase
		HALT bool

		// If SSTEP is set, the computer halts after every instruction
		SSTEP bool
	}

	// The IOT struct holds IO devices
	IOT struct {
		// Teleprinter flag
		PRINTER bool
	}

	// Gui attached to this computer
	g *gocui.Gui
}

// This function handles the HALT state, listening for inputs
func (mk *MK12) halt() {
	// Listen for keyboard inputs
	if !mk.STATE.SSTEP {
		debugPrint(mk.g, "** SYSTEM HALTED **  [ENTER] CONTINUE  |  [CTRL] + [C] EXIT  |  [SPACE] SINGLE STEP")
		updateStatus(mk.g, "HALT", gocui.AttrBold|gocui.ColorRed)
	} else {
		updateStatus(mk.g, "STEP", gocui.AttrBold|gocui.ColorBlue)
	}

	for mk.STATE.HALT {
		c := getLastKey()
		switch *c {
		case '\n':
			// debugPrint(mk.g, "Detected [ENTER]")
			mk.STATE.SSTEP = false
			mk.STATE.HALT = false

		case ' ':
			// debugPrint(mk.g, "Detected [SPACE]")
			mk.STATE.SSTEP = true
			mk.STATE.HALT = false

		case 0: // No keypress

		default:
			debugPrint(mk.g, fmt.Sprintf("UNKNOWN\t'%c'", *c))
		}
	}
}

// This function implements the fetch process:
//  1. Load PC into MA, MB
//  2. Increment PC
//  3. Load instruction into IR using MA
//  4. Determines the Effective Address (EA) for memory reference instructions and loads it into MA
//     4a) If page bit is set, use the current page. If not set, use page 0
//     4b) If indirect bit is set, the EA contains the actual address to use
//  5. Fetches the Content of the Effective Address (CA) for instructions that require an operand
func (mk *MK12) fetch() {

	// Check if step button pressed
	if c := getLastKey(); *c == ' ' {
		mk.STATE.SSTEP = true
	}

	// Set HALT if single stepping
	if mk.STATE.SSTEP {
		mk.STATE.HALT = true
	}

	// Catch halt
	if mk.STATE.HALT {
		mk.halt()
	}

	// Update to RUN status after returning from HALT or STEP
	updateStatus(mk.g, "RUN", gocui.AttrBold|gocui.ColorGreen)

	// Load PC into MA to get next instruction,
	// Save PC into MB for later use (indirect addressing)
	mk.MA = mk.PC
	mk.MB = mk.PC

	// Increment PC to point to the next instruction to execute
	mk.PC, _ = MKadd(mk.PC, 1)

	// Load instruction register
	mk.IR = mk.MEM[mk.MA]

	// Shorthand variable for the current instruction operator
	inOpr := mk.IR >> 9
	// Load correct address and/or operand for memory reference instructions
	if inOpr <= JMP {
		var addr int16
		if (mk.IR & 0b0000000010000000) > 0 {
			// If page bit is set, we use the current page
			addr = mk.PC & 0b0000111110000000
		} else {
			// If bit is not set, we use the first page
			addr = 0
		}
		// Fill in word address in page
		addr = addr | (mk.IR & 0b0000000001111111)

		// Check if indirect bit is set
		if (mk.IR & 0b0000000100000000) > 0 {

			// Auto increment addresses 0o10 0o17
			if (addr >= AUTO_begin) && (addr <= AUTO_end) {
				mk.MEM[addr], _ = MKadd(mk.MEM[addr], 1)
			}

			// Get address stored at addr
			addr = mk.MEM[addr]
		}

		// Store address in MA
		mk.MA = addr
	}

	// Load data from address for data reference instructions
	if inOpr == AND || inOpr == TAD || inOpr == ISZ {
		mk.MB = mk.MEM[mk.MA]
	}
}

// Executes the fetched instruction
func (mk *MK12) execute() {

	switch mk.IR >> 9 {

	case AND:
		// AND data with AC and store it back in AC
		tAC := mk.AC & mk.MB
		debugPrint(mk.g, fmt.Sprintf("AND %o & %o = %o --> AC", mk.AC, mk.MB, tAC))
		mk.AC = tAC
	case TAD:
		tAC, c := MKadd(mk.AC, mk.MB)
		debugPrint(mk.g, fmt.Sprintf("TAD %d + %d = %d --> AC", mk.AC, mk.MB, tAC))
		mk.L = c
		mk.AC = tAC
	case ISZ:
		// Increment MB and store it in MEM
		mk.MB, _ = MKadd(mk.MB, 1)
		mk.MEM[mk.MA] = mk.MB
		// If MB is zero, skip next instruction
		if mk.MB == 0 {
			debugPrint(mk.g, fmt.Sprintf("ISZ %o + 1 = %o --> %o; SKP %o", mk.MB-1, mk.MB, mk.MA, mk.PC))
			mk.PC = mk.PC + 1
		} else {
			debugPrint(mk.g, fmt.Sprintf("ISZ %o + 1 = %o --> %o", mk.MB-1, mk.MB, mk.MA))
		}
	case DCA:
		mk.MB = mk.AC
		mk.MEM[mk.MA] = mk.MB
		mk.AC = 0
		debugPrint(mk.g, fmt.Sprintf("DCA %o --> %o ; 0 --> AC", mk.MB, mk.MA))
	case JMS:
		mk.MEM[mk.MA] = mk.PC
		debugPrint(mk.g, fmt.Sprintf("JMS %o ; RET %o", mk.MA, mk.PC))
		mk.PC = mk.MA + 1
	case JMP:
		// Jump to the address stored in MA by storing it in the PC
		mk.PC = mk.MA
		debugPrint(mk.g, fmt.Sprintf("JMP %o", mk.MA))
	case IOT:
		devAddr := (mk.IR >> 3) & 0o77
		devReq := mk.IR & 0o7

		var debugInst string
		switch devAddr {
		case 4: // Teletype teleprinter/punch (stdout)
			debugInst += "IOT TELEPRINTER "
			if devReq == 1 { // Skip if flag is true
				debugInst += "SKIP"
				if mk.IOT.PRINTER {
					mk.PC = mk.PC + 1
					debugInst += "(T) "
				} else {
					debugInst += "(F) "
				}
			}
			if devReq == 2 { // Clear flag
				mk.IOT.PRINTER = false
				debugInst += "CLF "
			}
			if devReq == 4 { // Load char from AC and print it
				consolePrint(mk.g, fmt.Sprintf("%c", mk.AC))
				mk.IOT.PRINTER = true
				debugInst += "PRINT"
			}

			if devReq == 6 { // Combo of the two above: print char and clear flag
				mk.IOT.PRINTER = false
				consolePrint(mk.g, fmt.Sprintf("%c", mk.AC))
				mk.IOT.PRINTER = true
				debugInst += "PRINT CLF"
			}
		default:
			panic("IOT device not implemented")
		}

		if debugInst != "" {
			debugPrint(mk.g, debugInst)
		}

	case OPR:
		// if mk.IR == 0o7402 {
		// 	mk.STATE.HALT = true
		// 	fmt.Println("HALT")
		// }
		// if mk.IR == 0o7000 {
		// 	time.Sleep(time.Millisecond)
		// }

		group := (mk.IR >> 8) & 1
		if group > 0 {
			group += mk.IR & 1
		}

		switch group {
		case OPR_GROUP_1:
			var debugInst string = "OPR "

			if ((mk.IR >> 7) & 1) == 1 { // CLA - Clear Accumulator
				mk.AC = 0
				debugInst += "CLA "
			}
			if ((mk.IR >> 6) & 1) == 1 { // CLL - Clear Link
				mk.L = false
				debugInst += "CLL "
			}

			if ((mk.IR >> 5) & 1) == 1 { // CMA - Complement Accumulator
				mk.AC = MKcomplement(mk.AC)
				debugInst += "CMA "
			}
			if ((mk.IR >> 4) & 1) == 1 { // CML - Complement Link
				if mk.L {
					mk.L = false
				} else {
					mk.L = true
				}
				debugInst += "CML "
			}

			if ((mk.IR) & 1) == 1 { // IAC - Increment Accumulator
				mk.AC, mk.L = MKadd(mk.AC, 1)
				debugInst += "IAC"
			}

			if ((mk.IR >> 3) & 1) == 1 { // RAR
				panic("rotate right not implemented")
			}
			if ((mk.IR >> 2) & 1) == 1 { // RAL
				panic("rotate left not implemented")
			}
			if ((mk.IR >> 1) & 1) == 1 { // Rotate twice
				panic("rotate not implemented")
			}

			if debugInst != "" {
				debugPrint(mk.g, debugInst)
			}

		case OPR_GROUP_2:
			var debugInst string = "OPR "
			if ((mk.IR >> 7) & 1) == 1 { // CLA - Clear AC
				mk.AC = 0
				debugInst += "CLA "
			}

			// Determine state of skip conditions
			debugInst += "("
			skip := false
			if ((mk.IR >> 6) & 1) == 1 { // SMA - Skip on AC < 0
				if mk.AC < 0 {
					skip = true
				}
				debugInst += "SMA "
			}
			if ((mk.IR >> 5) & 1) == 1 { // SZA - Skip on AC == 0
				if mk.AC == 0 {
					skip = true
				}
				debugInst += "SZA "
			}
			if ((mk.IR >> 4) & 1) == 1 { // SNL - Skip on L == 1
				if mk.L {
					skip = true
				}
				debugInst += "SNL "
			}
			debugInst += ")"
			// Do the actual skip
			if ((mk.IR >> 3) & 1) == 1 { // Sense of skip (any or none)
				// If bit is set, no skip occurs if any condition has been satisfied (skip=true)
				if !skip {
					mk.PC = mk.PC + 1
					debugInst += "SKIP[NOR] "
				}
			} else {
				// If bit is not set, skip occurs if any condition is satisfied
				if skip {
					debugInst += "SKIP[OR] "
					mk.PC = mk.PC + 1
				}
			}

			if ((mk.IR >> 2) & 1) == 1 { // OSR - OR switch register with AC
				mk.AC |= mk.SR
				debugInst += "OSR "
			}
			if ((mk.IR >> 1) & 1) == 1 { // HLT - Halt the system
				debugInst += "HALT"
				mk.STATE.HALT = true
			}

			if debugInst != "" {
				debugPrint(mk.g, debugInst)
			}

		case OPR_GROUP_3:
			panic("group 3 operate instructions not implemented")
		}

	default:
		panic("ope")
	}
}

func (mk *MK12) run() {
	// Init the registers with some default value because we get stuck in the first fetch halt loop
	updateRegister(mk.g, "accumulator-register", mk.AC)
	updateRegister(mk.g, "counter-register", mk.PC)
	updateRegister(mk.g, "instruction-register", mk.IR)
	updateRegister(mk.g, "address-register", mk.MA)
	updateRegister(mk.g, "buffer-register", mk.MB)
	updateRegister(mk.g, "switch-register", mk.SR)
	// Loop forever
	for {
		mk.fetch()
		// Update SR after we fetch because we might be returning from a HALT, so
		// the switches might have changed. Update it before execute for same reason
		mk.SR = int16(switchRegister)

		mk.execute()

		updateRegister(mk.g, "accumulator-register", mk.AC)
		updateRegister(mk.g, "counter-register", mk.PC)
		updateRegister(mk.g, "instruction-register", mk.IR)
		updateRegister(mk.g, "address-register", mk.MA)
		updateRegister(mk.g, "buffer-register", mk.MB)

		time.Sleep(10 * time.Millisecond)
	}
}

func main() {
	// Check arguments
	if len(os.Args) == 1 {
		fmt.Println("Usage:", os.Args[0], "<infile>")
		os.Exit(1)
	}

	// Create our console user interface for the computer
	gui := CUI()
	defer gui.Close()

	// Create a new MK-12 computer
	myMK12 := MK12{g: gui}

	// Load our compiled object file into memory
	inFile := os.Args[1]
	m, err := LoadPObjFile(inFile)
	if err != nil {
		panic(err)
	}
	myMK12.MEM = m

	// Set PC to RESET vector and start computer HALTed
	myMK12.PC = 0o200
	myMK12.STATE.HALT = true
	go myMK12.run()

	// CUI Loop
	if err := gui.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}
