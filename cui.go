package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jroimartin/gocui"
)

var lastKey byte
var switchRegister int

type CUIFrontPanel struct {
	g *gocui.Gui
}

func (fp *CUIFrontPanel) PowerOn(mk MK12) {
	// Initialize Console interface on powerup and save it
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	fp.g = g

	// Layout
	g.SetManagerFunc(layout)

	// Keybindings
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("", gocui.KeyEnter, gocui.ModNone, proceed); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("", gocui.KeySpace, gocui.ModNone, step); err != nil {
		log.Panicln(err)
	}

	// F1-F12 Keys for Switch register
	if err := g.SetKeybinding("", gocui.KeyF1, gocui.ModNone, switchRegister1); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyF2, gocui.ModNone, switchRegister2); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyF3, gocui.ModNone, switchRegister3); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyF4, gocui.ModNone, switchRegister4); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyF5, gocui.ModNone, switchRegister5); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyF6, gocui.ModNone, switchRegister6); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyF7, gocui.ModNone, switchRegister7); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyF8, gocui.ModNone, switchRegister8); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyF9, gocui.ModNone, switchRegister9); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyF10, gocui.ModNone, switchRegister10); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyF11, gocui.ModNone, switchRegister11); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyF12, gocui.ModNone, switchRegister12); err != nil {
		log.Panicln(err)
	}

	// Start CUI loop
	go g.MainLoop()
}

func (fp *CUIFrontPanel) PowerOff() {
	fp.g.Close()
}

func (fp *CUIFrontPanel) Update(mk MK12) {
	var status string
	var attr = gocui.AttrBold
	if mk.STATE.HALT {
		status = "HALT"
		attr |= gocui.ColorRed
	} else if mk.STATE.SSTEP {
		status = "STEP"
		attr |= gocui.ColorBlue
	} else {
		status = "RUN"
		attr |= gocui.ColorGreen
	}
	updateStatus(fp.g, status, attr)
	updateRegister(fp.g, "accumulator-register", mk.AC)
	updateRegister(fp.g, "counter-register", mk.PC)
	updateRegister(fp.g, "instruction-register", mk.IR)
	updateRegister(fp.g, "address-register", mk.MA)
	updateRegister(fp.g, "buffer-register", mk.MB)
	updateRegister(fp.g, "switch-register", mk.SR)
	debugPrint(fp.g, mk.IRd)
}

func (fp *CUIFrontPanel) ReadSwitches() int16 {
	return int16(switchRegister)
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	// Register size
	regNum := 6
	regWStart := 0
	regWidth := 15
	regWEnd := regWStart + regWidth
	regHeight := 2
	regHStop := maxY - 3
	regHStart := regHStop - ((regHeight + 1) * regNum)
	regHEnd := regHStart + regHeight

	// Console size
	consoleWStart := regWEnd + 1
	consoleWidth := (maxX - 1) - consoleWStart
	consoleWEnd := consoleWStart + consoleWidth
	consoleHStart := 1
	consoleHEnd := maxY - 4

	// Debug Console
	dconsoleWStart := consoleWStart
	dconsoleWEnd := consoleWEnd
	dconsoleHStart := consoleHEnd + 1
	dconsoleHEnd := maxY - 1

	// Status text
	statusWStart := regWStart
	statusWEnd := regWEnd
	statusHStart := regHStart - 3
	statusHEnd := regHStart - 1

	// Title Text
	titleWStart := regWStart
	titleWEnd := regWEnd
	titleHStart := regHStop
	titleHEnd := maxY - 1

	// Program Counter
	if v, err := g.SetView("counter-register", regWStart, regHStart, regWEnd, regHEnd); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = " PC "
	}
	regHStart = regHEnd + 1
	regHEnd = regHStart + regHeight
	// Instruction Register
	if v, err := g.SetView("instruction-register", regWStart, regHStart, regWEnd, regHEnd); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = " IR "
	}
	regHStart = regHEnd + 1
	regHEnd = regHStart + regHeight
	// Accumulator Register
	if v, err := g.SetView("accumulator-register", regWStart, regHStart, regWEnd, regHEnd); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = " AC "
	}
	regHStart = regHEnd + 1
	regHEnd = regHStart + regHeight
	// Memory Address Register
	if v, err := g.SetView("address-register", regWStart, regHStart, regWEnd, regHEnd); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = " MA "
	}
	regHStart = regHEnd + 1
	regHEnd = regHStart + regHeight
	// Memory Buffer Register
	if v, err := g.SetView("buffer-register", regWStart, regHStart, regWEnd, regHEnd); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = " MB "
	}
	regHStart = regHEnd + 1
	regHEnd = regHStart + regHeight
	// Switch Register
	if v, err := g.SetView("switch-register", regWStart, regHStart, regWEnd, regHEnd); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = " SR "
	}

	// Teletype printer + keyboard
	if v, err := g.SetView("teletype", consoleWStart, consoleHStart, consoleWEnd, consoleHEnd); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = " CONSOLE "
		v.Wrap = true
		v.Autoscroll = true
	}

	// Debug command console
	if v, err := g.SetView("dbg-console", dconsoleWStart, dconsoleHStart, dconsoleWEnd, dconsoleHEnd); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = " DEBUG "
		v.Autoscroll = true
	}

	// Program Title Text
	if v, err := g.SetView("title-text", titleWStart, titleHStart, titleWEnd, titleHEnd); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.FgColor = gocui.AttrBold
		s := "MKSIM"
		centerd := fmt.Sprintf("%*s", -regWidth, fmt.Sprintf("%*s", (regWidth+len(s))/2, s))
		fmt.Fprint(v, centerd)
	}
	// Status Text
	if v, err := g.SetView("status-text", statusWStart, statusHStart, statusWEnd, statusHEnd); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.FgColor = gocui.AttrBold
	}

	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	g.Close()
	os.Exit(0)
	return gocui.ErrQuit
}

func proceed(g *gocui.Gui, v *gocui.View) error {
	lastKey = '\n'
	return nil
}

func step(g *gocui.Gui, v *gocui.View) error {
	lastKey = ' '
	return nil
}

func switchRegister1(g *gocui.Gui, v *gocui.View) error {
	switchRegister ^= 0b100000000000
	updateRegister(g, "switch-register", int16(switchRegister))
	return nil
}
func switchRegister2(g *gocui.Gui, v *gocui.View) error {
	switchRegister ^= 0b010000000000
	updateRegister(g, "switch-register", int16(switchRegister))
	return nil
}
func switchRegister3(g *gocui.Gui, v *gocui.View) error {
	switchRegister ^= 0b001000000000
	updateRegister(g, "switch-register", int16(switchRegister))
	return nil
}
func switchRegister4(g *gocui.Gui, v *gocui.View) error {
	switchRegister ^= 0b000100000000
	updateRegister(g, "switch-register", int16(switchRegister))
	return nil
}
func switchRegister5(g *gocui.Gui, v *gocui.View) error {
	switchRegister ^= 0b000010000000
	updateRegister(g, "switch-register", int16(switchRegister))
	return nil
}
func switchRegister6(g *gocui.Gui, v *gocui.View) error {
	switchRegister ^= 0b000001000000
	updateRegister(g, "switch-register", int16(switchRegister))
	return nil
}
func switchRegister7(g *gocui.Gui, v *gocui.View) error {
	switchRegister ^= 0b000000100000
	updateRegister(g, "switch-register", int16(switchRegister))
	return nil
}
func switchRegister8(g *gocui.Gui, v *gocui.View) error {
	switchRegister ^= 0b000000010000
	updateRegister(g, "switch-register", int16(switchRegister))
	return nil
}
func switchRegister9(g *gocui.Gui, v *gocui.View) error {
	switchRegister ^= 0b000000001000
	updateRegister(g, "switch-register", int16(switchRegister))
	return nil
}
func switchRegister10(g *gocui.Gui, v *gocui.View) error {
	switchRegister ^= 0b000000000100
	updateRegister(g, "switch-register", int16(switchRegister))
	return nil
}
func switchRegister11(g *gocui.Gui, v *gocui.View) error {
	switchRegister ^= 0b000000000010
	updateRegister(g, "switch-register", int16(switchRegister))
	return nil
}
func switchRegister12(g *gocui.Gui, v *gocui.View) error {
	switchRegister ^= 0b000000000001
	updateRegister(g, "switch-register", int16(switchRegister))
	return nil
}

var b byte

func getLastKey() byte {
	if lastKey != 0 {
		b = lastKey
		lastKey = 0
	} else {
		b = 0
	}
	return b
}

func updateRegister(g *gocui.Gui, registerName string, registerVal int16) {
	g.Update(func(g *gocui.Gui) error {
		v, err := g.View(registerName)
		if err != nil {
			return err
		}
		v.Clear()
		uVal := uint16(registerVal)
		fmt.Fprintf(v, " %.12b ", uVal)
		return nil
	})
}

func debugPrint(g *gocui.Gui, msg string) {
	g.Update(func(g *gocui.Gui) error {
		v, err := g.View("dbg-console")
		if err != nil {
			return err
		}
		fmt.Fprint(v, "\n"+msg)
		return nil
	})
}

func updateStatus(g *gocui.Gui, status string, atr gocui.Attribute) {
	regWidth := 15
	g.Update(func(g *gocui.Gui) error {
		v, err := g.View("status-text")
		if err != nil {
			return err
		}
		v.Clear()
		v.FgColor = atr
		centerd := fmt.Sprintf("%*s", -regWidth, fmt.Sprintf("%*s", (regWidth+len(status))/2, status))
		fmt.Fprint(v, centerd)
		return nil
	})
}

type CursedTeletype struct {
	g *gocui.Gui
}

// Printer
func (p *CursedTeletype) WriteByte(c byte) error {
	p.g.Update(func(g *gocui.Gui) error {
		v, err := g.View("teletype")
		if err != nil {
			return err
		}
		fmt.Fprintf(v, "%c", c)
		return nil
	})
	return nil
}

func (ct *CursedTeletype) Flush() error {
	return nil
}

func (ct *CursedTeletype) Available() int {
	return 1
}

// Keyboard
func (ct *CursedTeletype) ReadByte() (byte, error) {
	return getLastKey(), nil
}

func (ct *CursedTeletype) Buffered() int {
	if lastKey != 0 {
		return 1
	}
	return 0
}
