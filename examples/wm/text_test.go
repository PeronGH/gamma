package main

import (
	"testing"

	"github.com/charmbracelet/x/ansi"
)

func TestTruncateToWidthKeepsWholeGraphemes(t *testing.T) {
	got := truncateToWidth(ansi.GraphemeWidth, "ab🙂c", 3)
	if got != "ab" {
		t.Fatalf("expected wide grapheme to be excluded when it does not fit, got %q", got)
	}
}

func TestAppendWithinWidthKeepsWholeGraphemes(t *testing.T) {
	got := appendWithinWidth(ansi.GraphemeWidth, "A", "e\u0301🙂", 3)
	if got != "Ae\u0301" {
		t.Fatalf("expected combining cluster to fit and emoji to be excluded, got %q", got)
	}
}

func TestDropLastClusterRemovesWholeCluster(t *testing.T) {
	got := dropLastCluster("Ae\u0301🙂")
	if got != "Ae\u0301" {
		t.Fatalf("expected trailing emoji grapheme to be removed as one cluster, got %q", got)
	}

	got = dropLastCluster("e\u0301")
	if got != "" {
		t.Fatalf("expected combining cluster to be removed as one unit, got %q", got)
	}
}
