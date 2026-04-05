package main

import (
	gamma "github.com/PeronGH/gamma"
	"github.com/PeronGH/gamma/screen"
	"github.com/charmbracelet/x/ansi"
)

const (
	dialogW = 32
	dialogH = 5
)

var (
	dialogStyle      = gamma.Style{Fg: ansi.IndexedColor(236), Bg: ansi.IndexedColor(253)}
	dialogBorder     = gamma.Style{Fg: ansi.IndexedColor(240), Bg: ansi.IndexedColor(253)}
	dialogTitleStyle = gamma.Style{Fg: ansi.IndexedColor(255), Bg: ansi.IndexedColor(24), Attrs: gamma.AttrBold}
	inputStyle       = gamma.Style{Fg: ansi.IndexedColor(235), Bg: ansi.IndexedColor(255), Underline: gamma.UnderlineSingle}
	hintStyle        = gamma.Style{Fg: ansi.IndexedColor(245), Bg: ansi.IndexedColor(253)}
)

// Dialog is a modal text-input dialog.
type Dialog struct {
	Win   *gamma.Window
	Input string
	Open  bool
}

func NewDialog() *Dialog {
	return &Dialog{}
}

// Show opens the dialog, centers it, and clears the input.
func (d *Dialog) Show(parent *gamma.Window) {
	if d.Win != nil {
		parent.RemoveChild(d.Win)
	}

	d.Win = parent.NewWindow(0, 0, dialogW, dialogH)
	d.Input = ""
	d.Open = true
	d.Center(parent)
	d.paint()
	parent.BringToFront(d.Win)
}

// Hide closes the dialog and removes it from the tree.
func (d *Dialog) Hide(parent *gamma.Window) {
	if d.Win != nil {
		parent.RemoveChild(d.Win)
		d.Win = nil
	}
	d.Open = false
	d.Input = ""
}

// Center moves the dialog to the center of the parent window.
func (d *Dialog) Center(parent *gamma.Window) {
	if d.Win == nil {
		return
	}

	x := max(0, (parent.Width()-dialogW)/2)
	y := max(menuBarHeight, (parent.Height()-dialogH)/2)
	d.Win.MoveTo(x, y)
}

func (d *Dialog) maxInputWidth() int {
	if d.Win == nil {
		return dialogW - 4
	}
	return max(0, d.Win.Width()-4)
}

func (d *Dialog) paint() {
	if d.Win == nil {
		return
	}

	w, h := d.Win.Width(), d.Win.Height()
	innerWidth := max(0, w-4)
	d.Win.Fill(&gamma.Cell{Content: " ", Width: 1, Style: dialogStyle})

	for x := 1; x < w-1; x++ {
		d.Win.SetCell(x, 0, &gamma.Cell{Content: " ", Width: 1, Style: dialogTitleStyle})
		d.Win.SetCell(x, h-1, &gamma.Cell{Content: "─", Width: 1, Style: dialogBorder})
	}
	for y := 1; y < h-1; y++ {
		d.Win.SetCell(0, y, &gamma.Cell{Content: "│", Width: 1, Style: dialogBorder})
		d.Win.SetCell(w-1, y, &gamma.Cell{Content: "│", Width: 1, Style: dialogBorder})
	}
	d.Win.SetCell(0, 0, &gamma.Cell{Content: "╭", Width: 1, Style: dialogBorder})
	d.Win.SetCell(w-1, 0, &gamma.Cell{Content: "╮", Width: 1, Style: dialogBorder})
	d.Win.SetCell(0, h-1, &gamma.Cell{Content: "╰", Width: 1, Style: dialogBorder})
	d.Win.SetCell(w-1, h-1, &gamma.Cell{Content: "╯", Width: 1, Style: dialogBorder})

	ctx := screen.NewContext(d.Win)
	ctx.SetStyle(dialogTitleStyle)
	ctx.DrawString(truncateToWidth(d.Win.WidthMethod(), " new window ", innerWidth), 2, 0)

	ctx.SetStyle(dialogStyle)
	ctx.DrawString(truncateToWidth(d.Win.WidthMethod(), "Name", innerWidth), 2, 1)

	fieldWidth := d.maxInputWidth()
	for x := 0; x < fieldWidth; x++ {
		d.Win.SetCell(x+2, 2, &gamma.Cell{Content: " ", Width: 1, Style: inputStyle})
	}
	ctx.SetStyle(inputStyle)
	ctx.DrawString(truncateToWidth(d.Win.WidthMethod(), d.Input, fieldWidth), 2, 2)

	ctx.SetStyle(hintStyle)
	ctx.DrawString(truncateToWidth(d.Win.WidthMethod(), "Enter create  Esc cancel", innerWidth), 2, 3)
}

// HandleKey processes a key event. Returns (confirmed, handled).
func (d *Dialog) HandleKey(ev gamma.KeyPressEvent) (bool, bool) {
	if !d.Open || d.Win == nil {
		return false, false
	}

	switch {
	case ev.MatchString("enter"):
		return true, true
	case ev.MatchString("escape"):
		return false, true
	case ev.MatchString("backspace"):
		d.Input = dropLastCluster(d.Input)
		d.paint()
		return false, true
	default:
		if ev.Text == "" {
			return false, true
		}

		d.Input = appendWithinWidth(d.Win.WidthMethod(), d.Input, ev.Text, d.maxInputWidth())
		d.paint()
		return false, true
	}
}
