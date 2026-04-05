package main

import (
	"strings"

	gamma "github.com/PeronGH/gamma"
	"github.com/PeronGH/gamma/screen"
	"github.com/charmbracelet/x/ansi"
)

const (
	menuBarHeight = 1
	menuLabel     = " File "
)

var (
	menuBarStyle      = gamma.Style{Fg: ansi.IndexedColor(252), Bg: ansi.IndexedColor(237)}
	menuLabelStyle    = gamma.Style{Fg: ansi.IndexedColor(254), Bg: ansi.IndexedColor(239)}
	menuHotStyle      = gamma.Style{Fg: ansi.IndexedColor(232), Bg: ansi.IndexedColor(180), Attrs: gamma.AttrBold}
	menuPopupStyle    = gamma.Style{Fg: ansi.IndexedColor(236), Bg: ansi.IndexedColor(252)}
	menuPopupHotStyle = gamma.Style{Fg: ansi.IndexedColor(232), Bg: ansi.IndexedColor(180), Attrs: gamma.AttrBold}
)

type MenuAction int

const (
	menuActionNone MenuAction = iota
	menuActionNewWindow
	menuActionQuit
)

type menuItem struct {
	Action MenuAction
	Label  string
}

var menuItems = []menuItem{
	{Action: menuActionNewWindow, Label: " + New Window "},
	{Action: menuActionQuit, Label: " x Quit "},
}

// MenuBar is the top menu bar and its dropdown.
type MenuBar struct {
	Bar      *gamma.Window
	Dropdown *gamma.Window
	Open     bool
	HotFile  bool
	HotItem  MenuAction
}

// NewMenuBar creates the persistent menu bar as a child of parent.
func NewMenuBar(parent *gamma.Window, width int) *MenuBar {
	m := &MenuBar{
		Bar: parent.NewWindow(0, 0, width, menuBarHeight),
	}
	m.paintBar()
	return m
}

// Resize updates the menu bar width and repaints any open dropdown.
func (m *MenuBar) Resize(width int) {
	m.Bar.Resize(width, menuBarHeight)
	m.paintBar()
	if m.Dropdown != nil {
		m.paintDropdown()
	}
}

func (m *MenuBar) SetHover(fileHot bool, item MenuAction) {
	if item != menuActionNone && (!m.Open || m.Dropdown == nil) {
		item = menuActionNone
	}

	if m.HotFile == fileHot && m.HotItem == item {
		return
	}

	m.HotFile = fileHot
	m.HotItem = item
	m.paintBar()
	if m.Dropdown != nil {
		m.paintDropdown()
	}
}

func (m *MenuBar) paintBar() {
	m.Bar.Fill(&gamma.Cell{Content: " ", Width: 1, Style: menuBarStyle})

	labelStyle := menuLabelStyle
	if m.Open || m.HotFile {
		labelStyle = menuHotStyle
	}

	ctx := screen.NewContext(m.Bar)
	ctx.SetStyle(labelStyle)
	ctx.DrawString(menuLabel, 1, 0)
}

func (m *MenuBar) paintDropdown() {
	if m.Dropdown == nil {
		return
	}

	m.Dropdown.Fill(&gamma.Cell{Content: " ", Width: 1, Style: menuPopupStyle})

	ctx := screen.NewContext(m.Dropdown)
	width := m.dropdownWidth()
	for row, item := range menuItems {
		style := menuPopupStyle
		if item.Action == m.HotItem {
			style = menuPopupHotStyle
		}
		ctx.SetStyle(style)
		ctx.DrawString(padRight(item.Label, width), 0, row)
	}
}

func (m *MenuBar) dropdownWidth() int {
	width := 0
	for _, item := range menuItems {
		width = max(width, len(item.Label))
	}
	return width
}

func padRight(s string, width int) string {
	if len(s) >= width {
		return s
	}
	return s + strings.Repeat(" ", width-len(s))
}

// HitFileLabel returns true if (x, y) in root coordinates hits the "File" label.
func (m *MenuBar) HitFileLabel(x, y int) bool {
	return y == 0 && x >= 1 && x < 1+len(menuLabel)
}

// HitDropdownItem returns the action hit, or menuActionNone.
func (m *MenuBar) HitDropdownItem(x, y int) MenuAction {
	if !m.Open || m.Dropdown == nil {
		return menuActionNone
	}

	b := m.Dropdown.AbsoluteBounds()
	if y < b.Min.Y || y >= b.Max.Y || x < b.Min.X || x >= b.Max.X {
		return menuActionNone
	}

	row := y - b.Min.Y
	if row < 0 || row >= len(menuItems) {
		return menuActionNone
	}
	return menuItems[row].Action
}

// ShowDropdown shows or hides the dropdown by adding or removing it from the tree.
func (m *MenuBar) ShowDropdown(parent *gamma.Window, show bool) {
	if show {
		if m.Dropdown == nil {
			m.Dropdown = parent.NewWindow(1, menuBarHeight, m.dropdownWidth(), len(menuItems))
		}
		m.Open = true
		m.paintBar()
		m.paintDropdown()
		parent.BringToFront(m.Bar)
		parent.BringToFront(m.Dropdown)
		return
	}

	if m.Dropdown != nil {
		parent.RemoveChild(m.Dropdown)
		m.Dropdown = nil
	}
	m.Open = false
	m.HotItem = menuActionNone
	m.paintBar()
}
