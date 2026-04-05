package main

import (
	"testing"

	gamma "github.com/PeronGH/gamma"
)

func TestBottomRightClickStartsResizeInsteadOfDrag(t *testing.T) {
	a := NewApp(80, 24)
	a.createWindow("Resize")

	frame := a.frames[0]
	b := frame.Win.AbsoluteBounds()

	a.HandleEvent(gamma.MouseClickEvent{X: b.Max.X - 1, Y: b.Max.Y - 1, Button: gamma.MouseLeft})
	if a.resize != frame {
		t.Fatal("expected resize handle click to start resizing")
	}
	if a.drag != nil {
		t.Fatal("expected resize handle click to avoid drag mode")
	}

	a.HandleEvent(gamma.MouseMotionEvent{X: b.Min.X + 27, Y: b.Min.Y + 9})
	a.HandleEvent(gamma.MouseReleaseEvent{X: b.Min.X + 27, Y: b.Min.Y + 9, Button: gamma.MouseLeft})

	if frame.Win.Width() != 28 || frame.Win.Height() != 10 {
		t.Fatalf("expected resize interaction to change frame size, got %dx%d", frame.Win.Width(), frame.Win.Height())
	}
	if a.resize != nil {
		t.Fatal("expected mouse release to end resize mode")
	}
}
