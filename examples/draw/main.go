package main

import (
	"context"
	"log"
	"unicode"
	"unicode/utf8"

	gamma "github.com/PeronGH/gamma"
	"github.com/PeronGH/gamma/screen"
	"github.com/charmbracelet/x/ansi"
	"github.com/mattn/go-runewidth"
)

func main() {
	t := gamma.DefaultTerminal()
	scr := t.Screen()

	// Start in altscreen mode
	scr.EnterAltScreen()

	if err := t.Start(); err != nil {
		log.Fatalf("failed to start program: %v", err)
	}

	defer t.Stop()

	modes := []ansi.Mode{
		ansi.ButtonEventMouseMode,
		ansi.SgrExtMouseMode,
		ansi.FocusEventMode,
	}

	scr.WriteString(ansi.SetMode(modes...))

	// Listen for input and mouse events.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	const help = `Welcome to Draw Example!

Use the mouse to draw on the screen.
Press ctrl+c to exit.
Press esc to clear the screen.
Press alt+esc to reset the pen character, color, and the screen.
Press 0-9 to set the foreground color.
Press any other key to set the pen character.
Press ctrl+h for this help message.

Press any key to continue...`

	helpComp := gamma.NewStyledString(help)
	helpArea := helpComp.Bounds()
	helpW, helpH := helpArea.Dx(), helpArea.Dy()

	var prevHelpBuf *gamma.Buffer
	showingHelp := true
	displayHelp := func(show bool) {
		bounds := scr.Bounds()
		midX, midY := bounds.Dx()/2, bounds.Dy()/2
		x, y := midX-helpW/2, midY-helpH/2
		midArea := gamma.Rect(x, y, helpW, helpH)
		if show {
			// Save the area under the help to restore it later.
			prevHelpBuf = screen.CloneArea(scr, midArea)
			helpComp.Draw(scr, midArea)
		} else if prevHelpBuf != nil {
			// Restore saved area under the help.
			prevHelpBuf.Draw(scr, midArea)
		}
		scr.Render()
		scr.Flush()
	}

	clearScreen := func() {
		screen.Clear(scr)
		scr.Render()
		scr.Flush()
	}

	// Display first frame.
	displayHelp(showingHelp)

	const defaultChar = "█"
	pen := gamma.EmptyCell
	pen.Content = defaultChar
	draw := func(ev gamma.MouseEvent) {
		m := ev.Mouse()
		cur := scr.CellAt(m.X, m.Y)
		if cur == nil {
			// Position out of bounds.
			return
		}

		if cur.IsZero() && pen.Width == 1 {
			// Find the previous wide cell.
			var wide *gamma.Cell
			var wideX, wideY int
			for i := 1; i < 5 && m.X-i >= 0; i++ {
				wide = scr.CellAt(m.X-i, m.Y)
				if wide != nil && !wide.IsZero() && wide.Width > 1 {
					wideX, wideY = m.X-i, m.Y
					break
				}
			}

			if wide != nil {
				// Found a wide cell, make all cells blank.
				wc := *wide
				wc.Empty()
				scr.SetCell(wideX, wideY, &wc)
			}
		}

		// Can we fit the cell?
		fit := true
		if w := pen.Width; w > 1 {
			if cur.IsZero() || cur.Width > 1 {
				fit = false
			} else {
				for i := 1; i < w; i++ {
					cur = scr.CellAt(m.X+i, m.Y)
					if cur == nil || cur.IsZero() || cur.Width > 1 {
						// Position out of bounds or not empty.
						fit = false
						break
					}
				}
			}
		}
		if !fit {
			// Cell is too wide, ignore it.
			return
		}

		scr.SetCell(m.X, m.Y, &pen)
		scr.Render()
		scr.Flush()
	}
	displayHelp(showingHelp)

LOOP:
	for {
		select {
		case <-ctx.Done():
			break LOOP
		case ev := <-t.Events():
			switch ev := ev.(type) {
			case gamma.WindowSizeEvent:
				if showingHelp {
					displayHelp(false)
				}
				scr.Resize(ev.Width, ev.Height)
				if showingHelp {
					displayHelp(showingHelp)
				}
			case gamma.KeyPressEvent:
				if showingHelp {
					showingHelp = false
					displayHelp(showingHelp)
					break
				}
				switch {
				case ev.MatchString("ctrl+c"):
					cancel()
				case ev.MatchString("alt+esc"):
					pen.Style = gamma.Style{}
					pen.Content = defaultChar
					fallthrough
				case ev.MatchString("esc"):
					clearScreen()
				case ev.MatchString("ctrl+h"):
					showingHelp = true
					displayHelp(showingHelp)
				default:
					text := ev.Text
					if len(text) == 0 {
						break
					}
					r, rw := utf8.DecodeRuneInString(text)
					if rw == 1 && unicode.IsDigit(r) {
						pen.Style.Fg = ansi.Black + ansi.BasicColor(r-'0')
						break
					}
					pen.Content = text
					pen.Width = runewidth.RuneWidth(r)
				}
			case gamma.MouseClickEvent:
				if showingHelp {
					break
				}
				draw(ev)
			case gamma.MouseMotionEvent:
				if showingHelp || ev.Button == gamma.MouseNone {
					break
				}
				draw(ev)
			}
		}
	}

	scr.WriteString(ansi.ResetMode(modes...))
}
