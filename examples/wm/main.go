package main

import (
	"fmt"
	"log"
	"os"

	gamma "github.com/PeronGH/gamma"
	"github.com/charmbracelet/x/term"
)

func init() {
	f, err := os.OpenFile("gamma.log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		log.Fatalf("failed to open log file: %v", err)
	}
	log.SetOutput(f)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	stdin, stdout, environ := os.Stdin, os.Stdout, os.Environ()
	width, height, err := term.GetSize(stdout.Fd())
	if err != nil {
		log.Fatalf("failed to get terminal size: %v", err)
	}

	app := NewApp(width, height)
	if err := run(app, stdin, stdout, environ); err != nil {
		log.Fatalf("application error: %v", err)
	}
}

func run(app *App, input gamma.File, output gamma.File, environ []string) error {
	t := gamma.NewTerminal(gamma.NewConsole(input, output, environ), nil)
	scr := t.Screen()
	scr.EnterAltScreen()
	scr.HideCursor()

	if err := t.Start(); err != nil {
		return fmt.Errorf("failed to start terminal: %w", err)
	}
	defer t.Stop()

	scr.SetMouseMode(gamma.MouseModeMotion)

	// Initial render.
	if err := scr.Display(app); err != nil {
		return fmt.Errorf("failed to display: %w", err)
	}

	for !app.quit {
		ev := <-t.Events()
		switch ev := ev.(type) {
		case gamma.WindowSizeEvent:
			scr.Resize(ev.Width, ev.Height)
			app.Resize(ev.Width, ev.Height)
		}

		app.HandleEvent(ev)

		if err := scr.Display(app); err != nil {
			return fmt.Errorf("failed to display: %w", err)
		}
	}

	return nil
}
