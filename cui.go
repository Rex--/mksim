package main

import (
	"fmt"
	"log"

	"github.com/jroimartin/gocui"
)

var lastKey rune

func CUI() *gocui.Gui {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	// defer g.Close()

	g.SetManagerFunc(layout)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("", gocui.KeyEnter, gocui.ModNone, proceed); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("", gocui.KeySpace, gocui.ModNone, step); err != nil {
		log.Panicln(err)
	}

	// if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
	// 	log.Panicln(err)
	// }

	return g
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
		// fmt.Fprintln(v, " 000000000000")
	}
	regHStart = regHEnd + 1
	regHEnd = regHStart + regHeight
	// Instruction Register
	if v, err := g.SetView("instruction-register", regWStart, regHStart, regWEnd, regHEnd); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = " IR "
		// fmt.Fprintln(v, " 000000000000")
	}
	regHStart = regHEnd + 1
	regHEnd = regHStart + regHeight
	// Accumulator Register
	if v, err := g.SetView("accumulator-register", regWStart, regHStart, regWEnd, regHEnd); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = " AC "
		// fmt.Fprintln(v, " 000000000000")
	}
	regHStart = regHEnd + 1
	regHEnd = regHStart + regHeight
	// Memory Address Register
	if v, err := g.SetView("address-register", regWStart, regHStart, regWEnd, regHEnd); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = " MA "
		// fmt.Fprintln(v, " 000000000000")
	}
	regHStart = regHEnd + 1
	regHEnd = regHStart + regHeight
	// Memory Buffer Register
	if v, err := g.SetView("buffer-register", regWStart, regHStart, regWEnd, regHEnd); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = " MB "
		// fmt.Fprintln(v, " 000000000000")
	}
	regHStart = regHEnd + 1
	regHEnd = regHStart + regHeight
	// Switch Register
	if v, err := g.SetView("switch-register", regWStart, regHStart, regWEnd, regHEnd); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = " SR "
		// fmt.Fprintln(v, " 000000000000")
	}

	// Teletype printer + keyboard
	if v, err := g.SetView("teletype", consoleWStart, consoleHStart, consoleWEnd, consoleHEnd); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = " CONSOLE "
		v.Wrap = true
		v.Autoscroll = true
		// v.SelBgColor = gocui.ColorCyan
		// v.Highlight = true
		// fmt.Fprintln(v, "000000000000")
		// fmt.Fprintln(v, " Really long line.Really long line. Really long line. Really long line.Really long line. Really long line. Really long line. Really long line. Really long line. Really long line. Really long line. Really long line. Really long line. Really long line. Really long line. Really long line. Really long line. Really long line. Really long line. Really long line. Really long line. Really long line. Really long line. Really long line. Really long line. Really long line. Really long line. Really long line. Really long line. Really long line. Really long line. Really long line. Really long line. Really long line. Really long line. Really long line. Really long line. Really long line. Really long line. Really long line. Really long line. Really long line. Really long line. Really long line. Really long line. Really long line. Really long line. Really long line. Really long line. Really long line. Really long line. Really long line. Really long line. Really long line. Really long line. Really long line. Really long line. Really long line. Really long line. Really long line. Really long line. Really long line. Really long line. Really long line. Really long line. Really long line. Really long line. Really long line. Really long line. Really long line. Really long line.Really long line.Really long line.Really long line.Really long line.Really long line.Really long line. Really long line. Really long line. Really long line. Really long line. Really long line. Really long line.Really long line. Really long line.")
	}

	// Debug command console
	if v, err := g.SetView("dbg-console", dconsoleWStart, dconsoleHStart, dconsoleWEnd, dconsoleHEnd); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = " DEBUG "
		v.Autoscroll = true
		// fmt.Fprintf(v, "Hello, world!")
	}

	// Program Title Text
	if v, err := g.SetView("title-text", titleWStart, titleHStart, titleWEnd, titleHEnd); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.FgColor = gocui.AttrBold
		s := "MKSIM v0.0"
		centerd := fmt.Sprintf("%*s", -regWidth, fmt.Sprintf("%*s", (regWidth+len(s))/2, s))
		fmt.Fprint(v, centerd)
	}
	// Status Text
	if v, err := g.SetView("status-text", statusWStart, statusHStart, statusWEnd, statusHEnd); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.FgColor = gocui.AttrBold
		// centerd := fmt.Sprintf("%*s", -regWidth, fmt.Sprintf("%*s", (regWidth+len(s))/2, s))
		// fmt.Fprint(v, centerd)
	}

	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
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

var k rune

func getLastKey() *rune {
	if lastKey != 0 {
		k = lastKey
		lastKey = 0
	} else {
		k = 0
	}
	return &k
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

func consolePrint(g *gocui.Gui, msg string) {
	g.Update(func(g *gocui.Gui) error {
		v, err := g.View("teletype")
		if err != nil {
			return err
		}
		fmt.Fprint(v, msg)
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
