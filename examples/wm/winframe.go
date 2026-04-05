package main

import (
	gamma "github.com/PeronGH/gamma"
	"github.com/PeronGH/gamma/screen"
	"github.com/charmbracelet/x/ansi"
)

const (
	minFrameW = 20
	minFrameH = 6
)

var (
	frameBodyStyle             = gamma.Style{Bg: ansi.IndexedColor(252)}
	frameInactiveBodyStyle     = gamma.Style{Bg: ansi.IndexedColor(250)}
	frameTitleStyle            = gamma.Style{Fg: ansi.IndexedColor(255), Bg: ansi.IndexedColor(60), Attrs: gamma.AttrBold}
	frameTitleHoverStyle       = gamma.Style{Fg: ansi.IndexedColor(255), Bg: ansi.IndexedColor(67), Attrs: gamma.AttrBold}
	frameTitleInactiveStyle    = gamma.Style{Fg: ansi.IndexedColor(252), Bg: ansi.IndexedColor(240), Attrs: gamma.AttrBold}
	frameTitleInactiveHotStyle = gamma.Style{Fg: ansi.IndexedColor(255), Bg: ansi.IndexedColor(243), Attrs: gamma.AttrBold}
	frameCloseStyle            = gamma.Style{Fg: ansi.IndexedColor(203), Bg: ansi.IndexedColor(60), Attrs: gamma.AttrBold}
	frameCloseInactiveStyle    = gamma.Style{Fg: ansi.IndexedColor(245), Bg: ansi.IndexedColor(240), Attrs: gamma.AttrBold}
	frameCloseHoverStyle       = gamma.Style{Fg: ansi.IndexedColor(255), Bg: ansi.IndexedColor(160), Attrs: gamma.AttrBold}
	frameResizeStyle           = gamma.Style{Fg: ansi.IndexedColor(244), Bg: ansi.IndexedColor(252)}
	frameResizeInactiveStyle   = gamma.Style{Fg: ansi.IndexedColor(245), Bg: ansi.IndexedColor(250)}
	frameResizeHoverStyle      = gamma.Style{Fg: ansi.IndexedColor(255), Bg: ansi.IndexedColor(67), Attrs: gamma.AttrBold}
)

// WinFrame is a decorated window with a title bar and border.
type WinFrame struct {
	Win       *gamma.Window
	Name      string
	Active    bool
	TitleHot  bool
	CloseHot  bool
	ResizeHot bool
}

// NewWinFrame creates a new decorated window as a child of parent.
func NewWinFrame(parent *gamma.Window, name string, x, y int) *WinFrame {
	f := &WinFrame{
		Win:  parent.NewWindow(x, y, minFrameW, minFrameH),
		Name: name,
	}
	f.paint()
	return f
}

func (f *WinFrame) titleStyle() gamma.Style {
	if f.CloseHot || f.TitleHot {
		if f.Active {
			return frameTitleHoverStyle
		}
		return frameTitleInactiveHotStyle
	}
	if f.Active {
		return frameTitleStyle
	}
	return frameTitleInactiveStyle
}

func (f *WinFrame) bodyStyle() gamma.Style {
	if f.Active {
		return frameBodyStyle
	}
	return frameInactiveBodyStyle
}

func (f *WinFrame) closeStyle(titleStyle gamma.Style) gamma.Style {
	if f.CloseHot {
		return frameCloseHoverStyle
	}
	if f.Active {
		style := frameCloseStyle
		style.Bg = titleStyle.Bg
		return style
	}
	style := frameCloseInactiveStyle
	style.Bg = titleStyle.Bg
	return style
}

func (f *WinFrame) resizeStyle(bodyStyle gamma.Style) gamma.Style {
	if f.ResizeHot {
		return frameResizeHoverStyle
	}
	if f.Active {
		style := frameResizeStyle
		style.Bg = bodyStyle.Bg
		return style
	}
	style := frameResizeInactiveStyle
	style.Bg = bodyStyle.Bg
	return style
}

// SetActive toggles the active state and repaints when needed.
func (f *WinFrame) SetActive(active bool) {
	if f.Active == active {
		return
	}
	f.Active = active
	f.paint()
}

// SetHover updates title and close hover state and repaints when needed.
func (f *WinFrame) SetHover(titleHot, closeHot, resizeHot bool) {
	if f.TitleHot == titleHot && f.CloseHot == closeHot && f.ResizeHot == resizeHot {
		return
	}
	f.TitleHot = titleHot
	f.CloseHot = closeHot
	f.ResizeHot = resizeHot
	f.paint()
}

// MoveToClamped moves the frame while keeping it within the desktop bounds.
func (f *WinFrame) MoveToClamped(x, y, rootWidth, rootHeight, topInset int) {
	maxX := max(0, rootWidth-f.Win.Width())
	maxY := rootHeight - f.Win.Height()
	if maxY < topInset {
		maxY = topInset
	}

	x = max(0, min(x, maxX))
	y = max(topInset, min(y, maxY))
	f.Win.MoveTo(x, y)
}

// ClampToDesktop keeps the frame within the visible desktop area.
func (f *WinFrame) ClampToDesktop(rootWidth, rootHeight, topInset int) {
	b := f.Win.AbsoluteBounds()
	f.MoveToClamped(b.Min.X, b.Min.Y, rootWidth, rootHeight, topInset)
}

// FitToDesktop shrinks and repositions the frame only as much as needed to keep it visible.
func (f *WinFrame) FitToDesktop(rootWidth, rootHeight, topInset int) {
	availableWidth := max(1, rootWidth)
	availableHeight := max(1, rootHeight-topInset)

	newWidth := min(f.Win.Width(), availableWidth)
	newHeight := min(f.Win.Height(), availableHeight)

	b := f.Win.AbsoluteBounds()
	x := max(0, min(b.Min.X, rootWidth-newWidth))
	y := max(topInset, min(b.Min.Y, rootHeight-newHeight))

	f.Win.MoveTo(x, y)
	f.Win.Resize(newWidth, newHeight)
	f.paint()
}

// ResizeToClamped resizes the frame from its bottom-right corner while keeping it on screen.
func (f *WinFrame) ResizeToClamped(mouseX, mouseY, rootWidth, rootHeight int) {
	b := f.Win.AbsoluteBounds()
	maxWidth := max(1, rootWidth-b.Min.X)
	maxHeight := max(1, rootHeight-b.Min.Y)

	newWidth := mouseX - b.Min.X + 1
	newHeight := mouseY - b.Min.Y + 1

	newWidth = max(minFrameW, min(newWidth, maxWidth))
	newHeight = max(minFrameH, min(newHeight, maxHeight))

	if maxWidth < minFrameW {
		newWidth = maxWidth
	}
	if maxHeight < minFrameH {
		newHeight = maxHeight
	}

	f.Win.Resize(newWidth, newHeight)
	f.paint()
}

// paint fills the window buffer with a simple title bar and border.
func (f *WinFrame) paint() {
	w := f.Win.Width()
	h := f.Win.Height()

	bodyStyle := f.bodyStyle()
	titleStyle := f.titleStyle()

	f.Win.Fill(&gamma.Cell{Content: " ", Width: 1, Style: bodyStyle})

	for x := 0; x < w; x++ {
		f.Win.SetCell(x, 0, &gamma.Cell{Content: " ", Width: 1, Style: titleStyle})
	}

	ctx := screen.NewContext(f.Win)
	ctx.SetStyle(titleStyle)
	ctx.DrawString(truncateToWidth(f.Win.WidthMethod(), " "+f.Name+" ", max(0, w-6)), 1, 0)

	closeStart := w - 4
	ctx.SetStyle(f.closeStyle(titleStyle))
	ctx.DrawString("[X]", closeStart, 0)

	resizeStyle := f.resizeStyle(bodyStyle)
	f.Win.SetCell(w-1, h-1, &gamma.Cell{Content: "◢", Width: 1, Style: resizeStyle})
}

// HitTitleBar returns true if (x, y) in root coordinates is on this window's title bar.
func (f *WinFrame) HitTitleBar(x, y int) bool {
	b := f.Win.AbsoluteBounds()
	return y == b.Min.Y && x >= b.Min.X && x < b.Max.X
}

// HitClose returns true if (x, y) in root coordinates hits the close button.
func (f *WinFrame) HitClose(x, y int) bool {
	b := f.Win.AbsoluteBounds()
	closeStart := b.Max.X - 4
	return y == b.Min.Y && x >= closeStart && x < closeStart+3
}

// HitResize returns true if (x, y) in root coordinates hits the bottom-right resize handle.
func (f *WinFrame) HitResize(x, y int) bool {
	b := f.Win.AbsoluteBounds()
	return x >= b.Max.X-2 && x < b.Max.X && y >= b.Max.Y-2 && y < b.Max.Y
}
