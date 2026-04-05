package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"runtime/debug"
	"time"

	gamma "github.com/PeronGH/gamma"
	"github.com/PeronGH/gamma/screen"
)

type tickEvent struct{}

func main() {
	t := gamma.DefaultTerminal()
	scr := t.Screen()

	if err := t.Start(); err != nil {
		log.Fatalf("failed to start terminal: %v", err)
	}

	defer func() {
		if r := recover(); r != nil {
			_ = t.Stop()
			fmt.Fprintf(os.Stderr, "\r\nrecovered from panic: %v", r)
			debug.PrintStack()
		}
	}()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	counter := 5
	ticker := time.NewTicker(time.Second)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				t.SendEvent(tickEvent{})
			}
		}
	}()

	view := func(c int) string {
		return fmt.Sprintf("Panicing after %d seconds...\nPress 'q' or 'Ctrl+C' to exit.", c)
	}

OUT:
	for {
		select {
		case <-ctx.Done():
			break OUT
		case ev := <-t.Events():
			switch e := ev.(type) {
			case gamma.WindowSizeEvent:
				scr.Resize(e.Width, 2)
			case gamma.KeyPressEvent:
				switch {
				case e.MatchString("q", "ctrl+c"):
					cancel()
					break OUT
				}
			case tickEvent:
				counter--
				if counter < 0 {
					panic("Time's up!\n")
				}
			}
		}

		ss := gamma.NewStyledString(view(counter))
		screen.Clear(scr)
		scr.Display(ss)
	}

	ss := gamma.NewStyledString(view(counter) + "\n")
	screen.Clear(scr)
	scr.Display(ss)

	if err := t.Stop(); err != nil {
		log.Fatalf("failed to stop terminal: %v", err)
	}
}
