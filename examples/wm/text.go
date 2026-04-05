package main

import (
	"strings"

	gamma "github.com/PeronGH/gamma"
	"github.com/clipperhouse/uax29/v2/graphemes"
)

func truncateToWidth(method gamma.WidthMethod, s string, maxWidth int) string {
	if maxWidth <= 0 || s == "" {
		return ""
	}

	var b strings.Builder
	width := 0
	grs := graphemes.FromString(s)
	for grs.Next() {
		gr := string(grs.Value())
		grWidth := method.StringWidth(gr)
		if width+grWidth > maxWidth {
			break
		}
		width += grWidth
		b.WriteString(gr)
	}

	return b.String()
}

func appendWithinWidth(method gamma.WidthMethod, base, addition string, maxWidth int) string {
	if maxWidth <= 0 {
		return ""
	}

	base = truncateToWidth(method, base, maxWidth)
	width := method.StringWidth(base)
	if width >= maxWidth || addition == "" {
		return base
	}

	var b strings.Builder
	b.WriteString(base)

	grs := graphemes.FromString(addition)
	for grs.Next() {
		gr := string(grs.Value())
		grWidth := method.StringWidth(gr)
		if width+grWidth > maxWidth {
			break
		}
		width += grWidth
		b.WriteString(gr)
	}

	return b.String()
}

func dropLastCluster(s string) string {
	if s == "" {
		return ""
	}

	clusters := make([]string, 0, len(s))
	grs := graphemes.FromString(s)
	for grs.Next() {
		clusters = append(clusters, string(grs.Value()))
	}
	if len(clusters) == 0 {
		return ""
	}

	return strings.Join(clusters[:len(clusters)-1], "")
}
