package main

import (
	"fmt"
	"log"
	"math/rand"

	gamma "github.com/PeronGH/gamma"
	"github.com/PeronGH/gamma/screen"
	"github.com/charmbracelet/x/ansi"
)

func main() {
	t := gamma.DefaultTerminal()
	scr := t.Screen()

	if err := t.Start(); err != nil {
		log.Fatalf("failed to start program: %v", err)
	}

	defer t.Stop()

	scr.WriteString(ansi.SetWindowTitle("Hello, World!"))

	var st gamma.Style
	bg := 1
	st.Bg = ansi.BasicColor(bg)
	st.Fg = ansi.Black

	display := func() {
		const hw = "Hello, World!"
		bg := gamma.EmptyCell
		bg.Style = st
		screen.FillArea(scr, &bg, gamma.Rect(0, 0, scr.Bounds().Dx(), 1))
		for i, r := range hw {
			scr.SetCell(i, 0, &gamma.Cell{
				Content: string(r),
				Style:   st,
				Width:   1,
			})
		}
		scr.Render()
		scr.Flush()
	}

	// initial render
	display()

	var width int
	for ev := range t.Events() {
		switch ev := ev.(type) {
		case gamma.WindowSizeEvent:
			width = ev.Width
			scr.Resize(width, 1)
			display()
		case gamma.KeyPressEvent:
			switch {
			case ev.MatchString("q", "ctrl+c"):
				return
			}

			st.Bg = ansi.BasicColor(rand.Intn(16))
		}

		// Log event (this will appear above when we exit altscreen)
		scr.InsertAbove(fmt.Sprintf("%T %v", ev, ev))

		rd := rand.Intn(8)
		st.Bg = ansi.BasicColor(rd)
		display()
	}

	scr.WriteString(ansi.SetWindowTitle(""))
}
