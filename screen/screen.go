// Package screen provides functions and helpers to manipulate a [gamma.Screen].
package screen

import (
	gamma "github.com/PeronGH/gamma"
)

// Clear clears the screen with empty cells. This is equivalent to filling the
// screen with empty cells.
//
// If the screen implements a [Clear] method, it will be called instead of
// filling the screen with empty cells.
func Clear(scr gamma.Screen) {
	if c, ok := scr.(interface {
		Clear()
	}); ok {
		c.Clear()
		return
	}
	Fill(scr, nil)
}

// ClearArea clears the given area of the screen with empty cells. This is
// equivalent to filling the area with empty cells.
//
// If the screen implements a [ClearArea] method, it will be called instead of
// filling the area with empty cells.
func ClearArea(scr gamma.Screen, area gamma.Rectangle) {
	if c, ok := scr.(interface {
		ClearArea(area gamma.Rectangle)
	}); ok {
		c.ClearArea(area)
		return
	}
	FillArea(scr, nil, area)
}

// Fill fills the screen with the given cell. If the cell is nil, it fills the
// screen with empty cells.
//
// If the screen implements a [Fill] method, it will be called instead of
// filling the screen with empty cells.
func Fill(scr gamma.Screen, cell *gamma.Cell) {
	if f, ok := scr.(interface {
		Fill(cell *gamma.Cell)
	}); ok {
		f.Fill(cell)
		return
	}
	FillArea(scr, cell, scr.Bounds())
}

// FillArea fills the given area of the screen with the given cell. If the cell
// is nil, it fills the area with empty cells.
//
// If the screen implements a [FillArea] method, it will be called instead of
// filling the area with empty cells.
func FillArea(scr gamma.Screen, cell *gamma.Cell, area gamma.Rectangle) {
	if f, ok := scr.(interface {
		FillArea(cell *gamma.Cell, area gamma.Rectangle)
	}); ok {
		f.FillArea(cell, area)
		return
	}
	cellWidth := 1
	if cell != nil && cell.Width > 1 {
		cellWidth = cell.Width
	}
	for y := area.Min.Y; y < area.Max.Y; y++ {
		for x := area.Min.X; x < area.Max.X; x += cellWidth {
			scr.SetCell(x, y, cell)
		}
	}
}

// CloneArea clones the given area of the screen and returns a new buffer
// with the same size as the area. The new buffer will contain the same cells
// as the area in the screen.
// Use [gamma.Buffer.Draw] to draw the cloned buffer to a screen again.
//
// If the screen implements a [CloneArea] method, it will be called instead of
// cloning the area manually.
func CloneArea(scr gamma.Screen, area gamma.Rectangle) *gamma.Buffer {
	if c, ok := scr.(interface {
		CloneArea(area gamma.Rectangle) *gamma.Buffer
	}); ok {
		return c.CloneArea(area)
	}
	buf := gamma.NewBuffer(area.Dx(), area.Dy())
	for y := area.Min.Y; y < area.Max.Y; y++ {
		for x := area.Min.X; x < area.Max.X; {
			cell := scr.CellAt(x, y)
			if cell == nil || cell.IsZero() {
				x++
				continue
			}
			buf.SetCell(x-area.Min.X, y-area.Min.Y, cell)
			x += max(cell.Width, 1)
		}
	}
	return buf
}

// Clone creates a new [gamma.Buffer] clone of the given screen. The new buffer will
// have the same size as the screen and will contain the same cells.
// Use [gamma.Buffer.Draw] to draw the cloned buffer to a screen again.
//
// If the screen implements a [Clone] method, it will be called instead of
// cloning the entire screen manually.
func Clone(scr gamma.Screen) *gamma.Buffer {
	if c, ok := scr.(interface {
		Clone() *gamma.Buffer
	}); ok {
		return c.Clone()
	}
	return CloneArea(scr, scr.Bounds())
}
