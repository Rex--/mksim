package main

import (
	"io"
	"os"
)

type Device interface {
	// Select returns true if addr is addressed to this device, false otherwise.
	// Th device should get any values it needs from the registers (such as AC).
	Select(addr int16, mk *MK12) bool

	// Get is called to receive 12-bits of input from the device when the ORAC
	// control signal is asserted.
	Get() (data int16)

	// The IOT instruction is broken down into three stages, the device should
	// implement the proper functionality at each stage by returning 3 booleans
	// that represent control signals:
	// skp - IOSKIP - Skip the next instruction if this is true at any stage
	// clr - ACCLR  - Clear the AC register if this is true at any stage
	// or  - ORAC   - OR input data with AC and store it in AC
	Iop1() (skip bool, clr bool, or bool)
	Iop2() (skip bool, clr bool, or bool)
	Iop4() (skip bool, clr bool, or bool)
}

const (
	PT_READER   = 0o01
	PT_PUNCH    = 0o02
	TT_KEYBOARD = 0o03
	TT_PRINTER  = 0o04
)

/////////////////////
// TeleType Device
//

type TeleTypeKeyboard interface {
	ReadByte() (byte, error)
	Buffered() int
}

type TeleTypePrinter interface {
	WriteByte(byte) error
	Flush() error
	Available() int
}

type TeleTypeDevice struct {
	Device
	In  int // Last teletype input char  (keyboard)
	Out int // Last teletype output char (printer)

	Keyboard TeleTypeKeyboard
	Printer  TeleTypePrinter

	ac  int16 // Local copy of AC
	dev int   // Device currently being interfaced with (Keyboard or printer)
}

func (tt *TeleTypeDevice) Select(addr int16, mk *MK12) bool {

	if addr == TT_KEYBOARD || addr == TT_PRINTER {
		if addr == TT_PRINTER {
			tt.dev = TT_PRINTER
			// Save AC in case we need to print it
			tt.ac = mk.AC
		} else {
			tt.dev = TT_KEYBOARD
		}
		return true
	}
	return false
}

func (tt *TeleTypeDevice) Get() (data int16) {
	cData, err := tt.Keyboard.ReadByte()
	if err != nil {
		panic("Read error")
	}
	return int16(cData)
}

func (tt *TeleTypeDevice) Iop1() (skip bool, clr bool, or bool) {
	if tt.dev == TT_KEYBOARD {
		if tt.Keyboard.Buffered() > 0 { // KSF - Skip if incoming data available
			skip = true
		}
	} else if tt.dev == TT_PRINTER {
		if tt.Printer.Available() > 0 { // TSF - Skip if available to send
			skip = true
		}
	}
	return skip, clr, or
}

func (tt *TeleTypeDevice) Iop2() (skip bool, clr bool, or bool) {
	if tt.dev == TT_KEYBOARD { // KCC
		// Clear internal indicator of character ready (will be done when we read the byte)
		clr = true // Signal to clear AC in preparation of data exchange
	}
	return skip, clr, or
}

func (tt *TeleTypeDevice) Iop4() (skip bool, clr bool, or bool) {
	if tt.dev == TT_KEYBOARD { // KRS - Read keyboard buffer
		or = true
	} else if tt.dev == TT_PRINTER { // TPC - Print the contents of AC 4-11
		tt.Printer.WriteByte(byte(tt.ac & 0b000011111111))
		tt.Printer.Flush()
	}
	return skip, clr, or
}

///////////////////////////////////
// Paper Tape Reader/Punch Device (PC8-E)
//

type PaperTapeDevice struct {
	Device

	inTape  *os.File
	outTape *os.File

	// Read Buffer - Register to hold the character read from the tape
	RB int16
	// Reader Flag - Flag to signify if character is available to read
	RF bool

	// Punch buffer - register to hold the character to punch
	PB int16
	// Punch Flag - Denote a punch operation is complete
	PF bool

	// Device currently being interfaced with
	dev int
	// Local copy of AC
	ac int16
}

func (pt *PaperTapeDevice) Select(addr int16, mk *MK12) bool {
	if addr == PT_READER {
		pt.dev = PT_READER
		return true
	} else if addr == PT_PUNCH {
		pt.dev = PT_PUNCH
		// Flag starts true
		pt.PF = true
		// Save AC for use in PPC instruction
		pt.ac = mk.AC
		return true
	}
	return false
}

func (pt *PaperTapeDevice) Get() (data int16) {
	return pt.RB
}

func (pt *PaperTapeDevice) Iop1() (skip bool, clr bool, or bool) {
	if pt.dev == PT_READER {
		skip = pt.RF
	}
	if pt.dev == PT_PUNCH {
		skip = pt.PF
	}
	return
}

func (pt *PaperTapeDevice) Iop2() (skip bool, clr bool, or bool) {
	if pt.dev == PT_READER {
		or = true
	}
	if pt.dev == PT_PUNCH {
		// Clear punch flag
		pt.PF = false
		// Clear PB?
		// pt.PB = 0
	}
	return
}

func (pt *PaperTapeDevice) Iop4() (skip bool, clr bool, or bool) {
	if pt.dev == PT_READER {
		// Clear reader flag
		pt.RF = false
		// Read character(byte) into RB
		nextByte := make([]byte, 1)
		_, err := pt.inTape.Read(nextByte)
		if err != nil {
			if err == io.EOF {
				// End of file
				pt.RB = 0
			} else {
				println("papertape jam")
				panic(err)
			}
		} else {
			pt.RB = int16(nextByte[0])
		}
		// Set reader flag
		pt.RF = true
	}
	if pt.dev == PT_PUNCH {
		// Punch character
		_, err := pt.outTape.Write([]byte{byte(pt.ac)})
		if err != nil {
			panic(err)
		}
		// Set Punch flag
		pt.PF = true
	}
	return
}
