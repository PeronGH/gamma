package main

import (
	"fmt"
	"strings"

	gamma "github.com/PeronGH/gamma"
	"github.com/PeronGH/gamma/screen"
	"github.com/charmbracelet/x/ansi"
)

var desktopStyle = gamma.Style{Bg: ansi.IndexedColor(236)}

// App is the window manager application state.
type App struct {
	scr  *gamma.Window
	root *gamma.Window

	menu   *MenuBar
	dialog *Dialog
	frames []*WinFrame

	active *WinFrame
	drag   *WinFrame
	resize *WinFrame
	dragOX int // offset from window origin to mouse at drag start
	dragOY int
	nextX  int // cascade position for new windows
	nextY  int
	quit   bool
	winSeq int
}

func NewApp(width, height int) *App {
	a := &App{}
	a.scr = gamma.NewScreen(width, height)
	a.scr.SetWidthMethod(ansi.GraphemeWidth)
	a.root = a.scr.NewWindow(0, 0, width, height)
	a.menu = NewMenuBar(a.root, width)
	a.dialog = NewDialog()
	a.nextX = 2
	a.nextY = menuBarHeight + 1
	return a
}

// Draw implements gamma.Drawable.
func (a *App) Draw(scr gamma.Screen, area gamma.Rectangle) {
	screen.Fill(a.root, &gamma.Cell{Content: " ", Width: 1, Style: desktopStyle})
	a.root.Draw(scr, area)
}

func (a *App) Resize(width, height int) {
	a.scr.Resize(width, height)
	a.root.Resize(width, height)
	a.menu.Resize(width)

	for _, frame := range a.frames {
		frame.FitToDesktop(width, height, menuBarHeight)
	}

	if a.dialog.Open {
		a.dialog.Center(a.root)
		a.dialog.paint()
	}
}

func (a *App) nextUntitledName() string {
	a.winSeq++
	return fmt.Sprintf("Untitled %d", a.winSeq)
}

func (a *App) advanceCascade() {
	a.nextX += 2
	a.nextY += 1

	rw, rh := a.root.Width(), a.root.Height()
	if a.nextX+minFrameW > rw {
		a.nextX = 2
	}
	if a.nextY+minFrameH > rh {
		a.nextY = menuBarHeight + 1
	}
}

func (a *App) normalizeChrome() {
	a.root.BringToFront(a.menu.Bar)
	if a.menu.Dropdown != nil {
		a.root.BringToFront(a.menu.Dropdown)
	}
	if a.dialog.Win != nil {
		a.root.BringToFront(a.dialog.Win)
	}
}

func (a *App) clearHover() {
	a.menu.SetHover(false, menuActionNone)
	for _, frame := range a.frames {
		frame.SetHover(false, false, false)
	}
}

func (a *App) topmostFrame() *WinFrame {
	children := a.root.Children()
	for i := len(children) - 1; i >= 0; i-- {
		child := children[i]
		for _, frame := range a.frames {
			if frame.Win == child {
				return frame
			}
		}
	}
	return nil
}

func (a *App) setActiveFrame(target *WinFrame) {
	if a.active == target {
		return
	}
	if a.active != nil {
		a.active.SetActive(false)
	}
	a.active = target
	if a.active != nil {
		a.active.SetActive(true)
	}
}

func (a *App) bringFrameToFront(frame *WinFrame) {
	if frame == nil {
		return
	}
	a.root.BringToFront(frame.Win)
	a.setActiveFrame(frame)
	a.normalizeChrome()
}

func (a *App) createWindow(name string) {
	if strings.TrimSpace(name) == "" {
		name = a.nextUntitledName()
	}

	f := NewWinFrame(a.root, name, a.nextX, a.nextY)
	f.ClampToDesktop(a.root.Width(), a.root.Height(), menuBarHeight)
	a.frames = append(a.frames, f)

	a.bringFrameToFront(f)
	a.advanceCascade()
	a.clearHover()
}

func (a *App) removeFrame(f *WinFrame) {
	if f == nil {
		return
	}

	if f.Win.Parent() == a.root {
		a.root.RemoveChild(f.Win)
	}
	for i, frame := range a.frames {
		if frame == f {
			a.frames = append(a.frames[:i], a.frames[i+1:]...)
			break
		}
	}
	if a.drag == f {
		a.drag = nil
	}
	if a.resize == f {
		a.resize = nil
	}
	if a.active == f {
		a.active = nil
		a.setActiveFrame(a.topmostFrame())
	}
}

func (a *App) frameAt(x, y int) *WinFrame {
	hit := a.root.WindowAt(x, y)
	for _, frame := range a.frames {
		if frame.Win == hit {
			return frame
		}
	}
	return nil
}

func (a *App) updateHover(x, y int) {
	if a.dialog.Open {
		a.clearHover()
		return
	}

	a.menu.SetHover(a.menu.HitFileLabel(x, y), a.menu.HitDropdownItem(x, y))

	hit := a.frameAt(x, y)
	for _, frame := range a.frames {
		titleHot := false
		closeHot := false
		resizeHot := false
		if frame == hit {
			resizeHot = frame.HitResize(x, y)
			closeHot = frame.HitClose(x, y)
			titleHot = frame.HitTitleBar(x, y) && !closeHot && !resizeHot
		}
		frame.SetHover(titleHot, closeHot, resizeHot)
	}
}

// HandleEvent processes one input event. Returns true if the event was consumed.
func (a *App) HandleEvent(ev gamma.Event) bool {
	if a.dialog.Open {
		switch ev := ev.(type) {
		case gamma.KeyPressEvent:
			confirmed, handled := a.dialog.HandleKey(ev)
			if handled {
				if confirmed {
					name := a.dialog.Input
					a.dialog.Hide(a.root)
					a.createWindow(name)
				} else if ev.MatchString("escape") {
					a.dialog.Hide(a.root)
				}
				return true
			}
		case gamma.MouseClickEvent:
			if a.dialog.Win == nil {
				return true
			}
			p := gamma.Pos(ev.X, ev.Y)
			if !p.In(a.dialog.Win.AbsoluteBounds()) {
				a.dialog.Hide(a.root)
			}
			return true
		case gamma.MouseMotionEvent:
			a.clearHover()
			return true
		default:
			return true
		}

		return true
	}

	switch ev := ev.(type) {
	case gamma.KeyPressEvent:
		switch {
		case ev.MatchString("ctrl+c"):
			a.quit = true
			return true
		case ev.MatchString("escape"):
			if a.menu.Open {
				a.menu.ShowDropdown(a.root, false)
				a.updateHover(-1, -1)
				return true
			}
			a.quit = true
			return true
		}

	case gamma.MouseClickEvent:
		if ev.Button != gamma.MouseLeft {
			return false
		}

		if a.menu.HitFileLabel(ev.X, ev.Y) {
			a.menu.ShowDropdown(a.root, !a.menu.Open)
			a.normalizeChrome()
			a.updateHover(ev.X, ev.Y)
			return true
		}

		switch a.menu.HitDropdownItem(ev.X, ev.Y) {
		case menuActionNewWindow:
			a.menu.ShowDropdown(a.root, false)
			a.dialog.Show(a.root)
			a.clearHover()
			a.normalizeChrome()
			return true
		case menuActionQuit:
			a.menu.ShowDropdown(a.root, false)
			a.quit = true
			return true
		}

		if a.menu.Open {
			a.menu.ShowDropdown(a.root, false)
		}

		if frame := a.frameAt(ev.X, ev.Y); frame != nil {
			if frame.HitClose(ev.X, ev.Y) {
				a.removeFrame(frame)
				a.normalizeChrome()
				a.updateHover(ev.X, ev.Y)
				return true
			}

			a.bringFrameToFront(frame)

			if frame.HitResize(ev.X, ev.Y) {
				a.resize = frame
			} else if frame.HitTitleBar(ev.X, ev.Y) {
				b := frame.Win.AbsoluteBounds()
				a.drag = frame
				a.dragOX = ev.X - b.Min.X
				a.dragOY = ev.Y - b.Min.Y
			}
			a.updateHover(ev.X, ev.Y)
			return true
		}

		a.setActiveFrame(nil)
		a.updateHover(ev.X, ev.Y)

	case gamma.MouseMotionEvent:
		if a.resize != nil {
			a.resize.ResizeToClamped(ev.X, ev.Y, a.root.Width(), a.root.Height())
			a.updateHover(ev.X, ev.Y)
			return true
		}
		if a.drag != nil {
			newX := ev.X - a.dragOX
			newY := ev.Y - a.dragOY
			a.drag.MoveToClamped(newX, newY, a.root.Width(), a.root.Height(), menuBarHeight)
			a.updateHover(ev.X, ev.Y)
			return true
		}
		a.updateHover(ev.X, ev.Y)
		return false

	case gamma.MouseReleaseEvent:
		if a.resize != nil {
			a.resize = nil
			a.updateHover(ev.X, ev.Y)
			return true
		}
		if a.drag != nil {
			a.drag = nil
			a.updateHover(ev.X, ev.Y)
			return true
		}
	}

	return false
}
