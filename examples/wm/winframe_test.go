package main

import (
	"testing"

	gamma "github.com/PeronGH/gamma"
	"github.com/charmbracelet/x/ansi"
)

func newTestFrame(x, y int) *WinFrame {
	root := gamma.NewScreen(80, 30)
	root.SetWidthMethod(ansi.GraphemeWidth)
	return NewWinFrame(root, "test", x, y)
}

func TestHitResizeOnlyBottomRightCorner(t *testing.T) {
	frame := newTestFrame(5, 4)
	b := frame.Win.AbsoluteBounds()

	if !frame.HitResize(b.Max.X-1, b.Max.Y-1) {
		t.Fatal("expected bottom-right corner to hit resize handle")
	}
	if frame.HitResize(b.Min.X+1, b.Max.Y-1) {
		t.Fatal("expected bottom edge away from corner to not hit resize handle")
	}
	if frame.HitResize(b.Max.X-1, b.Min.Y) {
		t.Fatal("expected top-right corner to not hit resize handle")
	}
}

func TestResizeToClampedAppliesMinAndScreenBounds(t *testing.T) {
	frame := newTestFrame(5, 4)

	frame.ResizeToClamped(0, 0, 40, 20)
	if frame.Win.Width() != minFrameW || frame.Win.Height() != minFrameH {
		t.Fatalf("expected resize to clamp to minimum size, got %dx%d", frame.Win.Width(), frame.Win.Height())
	}

	frame.ResizeToClamped(200, 200, 40, 20)
	if frame.Win.Width() != 35 || frame.Win.Height() != 16 {
		t.Fatalf("expected resize to clamp to available space, got %dx%d", frame.Win.Width(), frame.Win.Height())
	}
}

func TestFitToDesktopShrinksOnlyAsNeeded(t *testing.T) {
	frame := newTestFrame(70, 25)
	frame.Win.Resize(50, 20)

	frame.FitToDesktop(30, 10, menuBarHeight)

	b := frame.Win.AbsoluteBounds()
	if b.Min.X != 0 || b.Min.Y != menuBarHeight {
		t.Fatalf("expected frame to move back onto screen, got (%d,%d)", b.Min.X, b.Min.Y)
	}
	if frame.Win.Width() != 30 || frame.Win.Height() != 9 {
		t.Fatalf("expected frame to shrink only to available size, got %dx%d", frame.Win.Width(), frame.Win.Height())
	}
}
