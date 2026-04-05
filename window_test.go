package gamma

import (
	"testing"

	"github.com/charmbracelet/x/ansi"
)

func TestWindowSetCellClipping(t *testing.T) {
	root := NewScreen(10, 5)
	win := root.NewWindow(2, 1, 4, 3)

	cell := &Cell{Content: "A", Width: 1}

	// Write inside bounds -- should succeed.
	win.SetCell(0, 0, cell)
	if c := win.CellAt(0, 0); c == nil || c.Content != "A" {
		t.Fatal("expected cell A at (0,0)")
	}

	// Write outside bounds -- should be silently dropped.
	win.SetCell(4, 0, cell)  // x == width
	win.SetCell(-1, 0, cell) // negative x
	win.SetCell(0, 3, cell)  // y == height
	win.SetCell(0, -1, cell) // negative y

	// Verify out-of-bounds reads return nil.
	if c := win.CellAt(4, 0); c != nil {
		t.Fatal("expected nil for out-of-bounds CellAt")
	}
}

func TestWindowWideCellAtBoundary(t *testing.T) {
	root := NewScreen(10, 5)
	win := root.NewWindow(0, 0, 5, 3)

	// A 2-cell-wide character at x=4 would overflow (4+2 > 5).
	wide := &Cell{Content: "漢", Width: 2, Style: Style{Attrs: AttrBold}}
	win.SetCell(4, 0, wide)

	c := win.CellAt(4, 0)
	if c == nil {
		t.Fatal("expected cell at (4,0)")
	}
	if c.Content != " " || c.Width != 1 {
		t.Fatalf("expected space cell, got %q width=%d", c.Content, c.Width)
	}
	// Style should be preserved.
	if c.Style.Attrs&AttrBold == 0 {
		t.Fatal("expected bold style preserved on truncated wide cell")
	}

	// A 2-cell-wide character at x=3 should fit (3+2 == 5).
	win.SetCell(3, 1, wide)
	c = win.CellAt(3, 1)
	if c == nil || c.Content != "漢" {
		t.Fatalf("expected wide cell at (3,1), got %v", c)
	}
}

func TestViewCoordinateTranslation(t *testing.T) {
	root := NewScreen(20, 10)

	// Fill root with dots.
	dot := &Cell{Content: ".", Width: 1}
	root.Fill(dot)

	// Create a view at (5, 3) of size 4x2.
	view := root.NewView(5, 3, 4, 2)

	// Write to view at (0, 0) -- should appear at (5, 3) in root buffer.
	cell := &Cell{Content: "X", Width: 1}
	view.SetCell(0, 0, cell)

	// Read back through view.
	c := view.CellAt(0, 0)
	if c == nil || c.Content != "X" {
		t.Fatal("expected X at view (0,0)")
	}

	// Read back through root buffer directly.
	c = root.Buffer.CellAt(5, 3)
	if c == nil || c.Content != "X" {
		t.Fatalf("expected X at root buffer (5,3), got %v", c)
	}

	// Writing outside view bounds should be dropped.
	view.SetCell(4, 0, cell) // x == view width
	c = root.Buffer.CellAt(9, 3)
	if c != nil && c.Content == "X" {
		t.Fatal("expected view to clip write at x=4")
	}
}

func TestViewCellAtTranslation(t *testing.T) {
	root := NewScreen(20, 10)

	// Write to root buffer directly.
	cell := &Cell{Content: "Z", Width: 1}
	root.Buffer.SetCell(7, 4, cell)

	// Create a view that covers that position.
	view := root.NewView(5, 3, 10, 5)

	// (7, 4) in root = (2, 1) in view.
	c := view.CellAt(2, 1)
	if c == nil || c.Content != "Z" {
		t.Fatalf("expected Z at view (2,1), got %v", c)
	}
}

func TestWindowFill(t *testing.T) {
	root := NewScreen(10, 5)
	win := root.NewWindow(0, 0, 3, 2)

	cell := &Cell{Content: "#", Width: 1}
	win.Fill(cell)

	for y := 0; y < 2; y++ {
		for x := 0; x < 3; x++ {
			c := win.CellAt(x, y)
			if c == nil || c.Content != "#" {
				t.Fatalf("expected # at (%d,%d), got %v", x, y, c)
			}
		}
	}

	// Cells outside the window should not be filled (owned buffer is 3x2).
	c := win.CellAt(3, 0)
	if c != nil {
		t.Fatal("expected nil outside window bounds")
	}
}

func TestWindowClear(t *testing.T) {
	root := NewScreen(10, 5)
	win := root.NewWindow(0, 0, 3, 2)

	cell := &Cell{Content: "#", Width: 1}
	win.Fill(cell)
	win.Clear()

	for y := 0; y < 2; y++ {
		for x := 0; x < 3; x++ {
			c := win.CellAt(x, y)
			if c == nil {
				t.Fatalf("expected non-nil cell at (%d,%d)", x, y)
			}
			// Clear fills with nil which SetCell converts to EmptyCell.
			if c.Content != " " {
				t.Fatalf("expected space at (%d,%d), got %q", x, y, c.Content)
			}
		}
	}
}

func TestWindowDraw(t *testing.T) {
	root := NewScreen(20, 10)
	win := root.NewWindow(0, 0, 3, 2)

	cell := &Cell{Content: "W", Width: 1}
	win.Fill(cell)

	// Draw the window at position (5, 3) on the root.
	win.Draw(root, Rect(5, 3, 3, 2))

	c := root.CellAt(5, 3)
	if c == nil || c.Content != "W" {
		t.Fatalf("expected W at root (5,3), got %v", c)
	}
	c = root.CellAt(7, 4)
	if c == nil || c.Content != "W" {
		t.Fatalf("expected W at root (7,4), got %v", c)
	}
}

func TestViewFill(t *testing.T) {
	root := NewScreen(20, 10)

	// Fill root with dots.
	dot := &Cell{Content: ".", Width: 1}
	root.Fill(dot)

	// Create a view and fill it.
	view := root.NewView(5, 3, 4, 2)
	cell := &Cell{Content: "#", Width: 1}
	view.Fill(cell)

	// View area in root should be filled.
	for y := 3; y < 5; y++ {
		for x := 5; x < 9; x++ {
			c := root.Buffer.CellAt(x, y)
			if c == nil || c.Content != "#" {
				t.Fatalf("expected # at root (%d,%d), got %v", x, y, c)
			}
		}
	}

	// Area outside view should still be dots.
	c := root.Buffer.CellAt(4, 3)
	if c == nil || c.Content != "." {
		t.Fatalf("expected . at root (4,3), got %v", c)
	}
	c = root.Buffer.CellAt(9, 3)
	if c == nil || c.Content != "." {
		t.Fatalf("expected . at root (9,3), got %v", c)
	}
}

func TestWindowWidthHeight(t *testing.T) {
	root := NewScreen(20, 10)
	win := root.NewWindow(3, 4, 7, 5)

	if w := win.Width(); w != 7 {
		t.Fatalf("expected width 7, got %d", w)
	}
	if h := win.Height(); h != 5 {
		t.Fatalf("expected height 5, got %d", h)
	}
}

func TestWindowBoundsReturnsLocalRect(t *testing.T) {
	root := NewScreen(20, 10)

	// Bounds should reflect position and size.
	win := root.NewWindow(3, 4, 7, 5)
	b := win.Bounds()
	if b.Min.X != 3 || b.Min.Y != 4 || b.Dx() != 7 || b.Dy() != 5 {
		t.Fatalf("unexpected bounds: %v", b)
	}
}

func TestWindowWideCellFitsExactly(t *testing.T) {
	root := NewScreen(10, 5)
	win := root.NewWindow(0, 0, 4, 1)

	// Width=2 cell at x=2 fits exactly (2+2 == 4).
	wide := &Cell{Content: "字", Width: 2}
	win.SetCell(2, 0, wide)

	c := win.CellAt(2, 0)
	if c == nil || c.Content != "字" || c.Width != 2 {
		t.Fatalf("expected wide cell at (2,0), got %v", c)
	}
}

func TestNewScreenIsNotView(t *testing.T) {
	root := NewScreen(10, 5)
	// Root window SetCell should work as before.
	cell := &Cell{Content: "R", Width: 1}
	root.SetCell(0, 0, cell)

	c := root.CellAt(0, 0)
	if c == nil || c.Content != "R" {
		t.Fatal("expected R at root (0,0)")
	}
}

// Ensure the Screen interface is satisfied.
func TestWindowImplementsScreen(t *testing.T) {
	var _ Screen = (*Window)(nil)
}

func TestWindowWidthMethod(t *testing.T) {
	root := NewScreen(10, 5)
	root.SetWidthMethod(ansi.GraphemeWidth)

	wm := root.WidthMethod()
	if wm == nil {
		t.Fatal("expected non-nil WidthMethod")
	}
}
