package gamma

import (
	"github.com/charmbracelet/x/ansi"
)

// Window represents a rectangular area on the screen. It can be a root window
// with no parent, or a sub-window with a parent window. A window can have its
// own buffer or share the buffer of its parent window (view).
type Window struct {
	*Buffer

	method   *WidthMethod
	parent   *Window
	children []*Window // ordered back-to-front (last = topmost)
	bounds   Rectangle
	view     bool // true if this window shares its parent's buffer
}

var (
	_ Screen   = (*Window)(nil)
	_ Drawable = (*Window)(nil)
)

// SetCell sets the cell at the given position, enforcing the window's bounds.
// Coordinates are relative to the window (0-based). Writes outside the
// window's bounds are silently dropped. If a wide cell would overflow the
// right edge, a space cell with the same style is written instead.
//
// For views (windows that share their parent's buffer), coordinates are
// translated to the parent buffer's coordinate space.
func (w *Window) SetCell(x, y int, c *Cell) {
	width := w.bounds.Dx()
	height := w.bounds.Dy()
	if x < 0 || x >= width || y < 0 || y >= height {
		return
	}

	if c != nil && c.Width > 1 && x+c.Width > width {
		// Wide cell would overflow the right edge; write a space instead.
		sc := *c
		sc.Content = " "
		sc.Width = 1
		c = &sc
	}

	if w.view {
		w.Buffer.SetCell(x+w.bounds.Min.X, y+w.bounds.Min.Y, c)
	} else {
		w.Buffer.SetCell(x, y, c)
	}
}

// CellAt returns the cell at the given position, enforcing the window's
// bounds. Coordinates are relative to the window (0-based). Returns nil for
// positions outside the window's bounds.
//
// For views (windows that share their parent's buffer), coordinates are
// translated to the parent buffer's coordinate space.
func (w *Window) CellAt(x, y int) *Cell {
	if x < 0 || x >= w.bounds.Dx() || y < 0 || y >= w.bounds.Dy() {
		return nil
	}

	if w.view {
		return w.Buffer.CellAt(x+w.bounds.Min.X, y+w.bounds.Min.Y)
	}
	return w.Buffer.CellAt(x, y)
}

// Fill fills the window with the given cell, respecting the window's bounds.
func (w *Window) Fill(c *Cell) {
	w.FillArea(c, Rect(0, 0, w.bounds.Dx(), w.bounds.Dy()))
}

// FillArea fills the specified area of the window with the given cell,
// respecting the window's bounds. The area is relative to the window.
func (w *Window) FillArea(c *Cell, area Rectangle) {
	// Clip area to window bounds.
	wbounds := Rect(0, 0, w.bounds.Dx(), w.bounds.Dy())
	area = area.Intersect(wbounds)
	if area.Empty() {
		return
	}

	cellWidth := 1
	if c != nil && c.Width > 1 {
		cellWidth = c.Width
	}
	for y := area.Min.Y; y < area.Max.Y; y++ {
		for x := area.Min.X; x < area.Max.X; x += cellWidth {
			w.SetCell(x, y, c)
		}
	}
}

// Clear clears the window with empty cells.
func (w *Window) Clear() {
	w.FillArea(nil, Rect(0, 0, w.bounds.Dx(), w.bounds.Dy()))
}

// ClearArea clears the specified area of the window with empty cells. The area
// is relative to the window.
func (w *Window) ClearArea(area Rectangle) {
	w.FillArea(nil, area)
}

// HasParent returns whether the window has a parent window. This can be used
// to determine if the window is a root window or a sub-window.
func (w *Window) HasParent() bool {
	return w.parent != nil
}

// Parent returns the parent window of the current window.
// If the window does not have a parent, it returns nil.
func (w *Window) Parent() *Window {
	return w.parent
}

// MoveTo moves the window to the specified x and y coordinates.
func (w *Window) MoveTo(x, y int) {
	size := w.bounds.Size()
	w.bounds.Min.X = x
	w.bounds.Min.Y = y
	w.bounds.Max.X = x + size.X
	w.bounds.Max.Y = y + size.Y
}

// MoveBy moves the window by the specified delta x and delta y.
func (w *Window) MoveBy(dx, dy int) {
	w.bounds.Min.X += dx
	w.bounds.Min.Y += dy
	w.bounds.Max.X += dx
	w.bounds.Max.Y += dy
}

// Children returns a copy of the window's children slice, ordered
// back-to-front (last element is topmost).
func (w *Window) Children() []*Window {
	return append([]*Window(nil), w.children...)
}

// RemoveChild removes a child from this window's children list and sets the
// child's parent to nil.
func (w *Window) RemoveChild(child *Window) {
	for i, c := range w.children {
		if c == child {
			w.children = append(w.children[:i], w.children[i+1:]...)
			child.parent = nil
			return
		}
	}
}

// BringToFront moves a child to the top of the z-order (end of children slice).
func (w *Window) BringToFront(child *Window) {
	for i, c := range w.children {
		if c == child {
			w.children = append(w.children[:i], w.children[i+1:]...)
			w.children = append(w.children, child)
			return
		}
	}
}

// SendToBack moves a child to the bottom of the z-order (start of children slice).
func (w *Window) SendToBack(child *Window) {
	for i, c := range w.children {
		if c == child {
			w.children = append(w.children[:i], w.children[i+1:]...)
			w.children = append([]*Window{child}, w.children...)
			return
		}
	}
}

// WindowAt returns the deepest child window containing the point (x, y),
// which is relative to this window. Walks children front-to-back (topmost
// first). Returns w itself if no child contains the point.
func (w *Window) WindowAt(x, y int) *Window {
	// Walk front-to-back (reverse order).
	for i := len(w.children) - 1; i >= 0; i-- {
		child := w.children[i]
		p := Pos(x, y)
		if p.In(child.bounds) {
			// Recurse with coordinates translated to child's space.
			return child.WindowAt(x-child.bounds.Min.X, y-child.bounds.Min.Y)
		}
	}
	return w
}

// Clone creates an exact copy of the window, including its buffer and values.
// The cloned window will have the same parent and method as the original
// window.
func (w *Window) Clone() *Window {
	return w.CloneArea(w.Buffer.Bounds())
}

// CloneArea creates an exact copy of the window, including its buffer and
// values, but only within the specified area. The cloned window will have the
// same parent and method as the original window, but its bounds will be
// limited to the specified area.
func (w *Window) CloneArea(area Rectangle) *Window {
	clone := new(Window)
	clone.Buffer = w.Buffer.CloneArea(area)
	clone.parent = w.parent
	clone.method = w.method
	clone.bounds = area
	return clone
}

// Resize resizes the window to the specified width and height.
func (w *Window) Resize(width, height int) {
	// Only resize the buffer if this window owns its buffer.
	if !w.view {
		w.Buffer.Resize(width, height)
	}
	w.bounds.Max.X = w.bounds.Min.X + width
	w.bounds.Max.Y = w.bounds.Min.Y + height
}

// Width returns the width of the window.
func (w *Window) Width() int {
	return w.bounds.Dx()
}

// Height returns the height of the window.
func (w *Window) Height() int {
	return w.bounds.Dy()
}

// WidthMethod returns the method used to calculate the width of characters in
// the window.
func (w *Window) WidthMethod() WidthMethod {
	return *w.method
}

// Bounds returns the bounds of the window as a rectangle.
func (w *Window) Bounds() Rectangle {
	return w.bounds
}

// NewWindow creates a new window with its own buffer relative to the parent
// window at the specified position and size.
//
// This will panic if width or height is negative.
func (w *Window) NewWindow(x, y, width, height int) *Window {
	return newWindow(w, x, y, width, height, w.method, false)
}

// NewView creates a new view into the parent window at the specified position
// and size. Unlike [Window.NewWindow], this view shares the same buffer as the
// parent window.
func (w *Window) NewView(x, y, width, height int) *Window {
	return newWindow(w, x, y, width, height, w.method, true)
}

// NewScreen creates a new root [Window] with the given size and width method.
//
// This will panic if width or height is negative.
func NewScreen(width, height int) *Window {
	var method WidthMethod = ansi.WcWidth
	return newWindow(nil, 0, 0, width, height, &method, false)
}

// Draw draws the window's content to the given screen at the specified area.
// For views, it reads from the correct region of the parent buffer.
func (w *Window) Draw(scr Screen, area Rectangle) {
	if area.Empty() {
		return
	}

	bounds := scr.Bounds()
	if !area.Overlaps(bounds) {
		return
	}

	ww := w.bounds.Dx()
	wh := w.bounds.Dy()
	for y := area.Min.Y; y < area.Max.Y; y++ {
		sy := y - area.Min.Y
		if sy < 0 || sy >= wh {
			continue
		}
		for x := area.Min.X; x < area.Max.X; {
			sx := x - area.Min.X
			if sx < 0 || sx >= ww {
				x++
				continue
			}
			c := w.CellAt(sx, sy)
			if c == nil || c.IsZero() {
				x++
				continue
			}
			scr.SetCell(x, y, c)
			width := c.Width
			if width <= 0 {
				width = 1
			}
			x += width
		}
	}

	// Draw children in z-order (back-to-front).
	for _, child := range w.children {
		childArea := child.bounds.Add(area.Min)
		child.Draw(scr, childArea)
	}
}

// SetWidthMethod sets the width method for the window.
func (w *Window) SetWidthMethod(method WidthMethod) {
	w.method = &method
}

// newWindow creates a new [Window] with the specified parent, position,
// method, and size.
func newWindow(parent *Window, x, y, width, height int, method *WidthMethod, view bool) *Window {
	w := new(Window)
	if view {
		w.Buffer = parent.Buffer
	} else {
		w.Buffer = NewBuffer(width, height)
	}
	w.parent = parent
	w.method = method
	w.bounds = Rect(x, y, width, height)
	w.view = view
	if parent != nil {
		parent.children = append(parent.children, w)
	}
	return w
}
