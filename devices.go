package main

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
	In  int // Last teletype input char  (keyboard)
	Out int // Last teletype output char (printer)

	Keyboard TeleTypeKeyboard
	Printer  TeleTypePrinter

	ac  int16 // Local copy of AC
	dev int   // Device currently being interfaced with (Keyboard or printer)
}

const (
	TT_KEYBOARD = 0o03
	TT_PRINTER  = 0o04
)

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
