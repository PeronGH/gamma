package gamma

import (
	"image/color"
	"testing"

	"github.com/charmbracelet/x/ansi"
)

func TestStyledString(t *testing.T) {
	cases := []struct {
		name           string
		input          string
		expected       *Buffer
		expectedWidth  int
		expectedHeight int
	}{
		{
			name:           "single line",
			input:          "Hello, World!",
			expectedWidth:  13,
			expectedHeight: 1,
			expected: &Buffer{
				Lines: []Line{
					{
						newWcCell("H", nil, nil),
						newWcCell("e", nil, nil),
						newWcCell("l", nil, nil),
						newWcCell("l", nil, nil),
						newWcCell("o", nil, nil),
						newWcCell(",", nil, nil),
						newWcCell(" ", nil, nil),
						newWcCell("W", nil, nil),
						newWcCell("o", nil, nil),
						newWcCell("r", nil, nil),
						newWcCell("l", nil, nil),
						newWcCell("d", nil, nil),
						newWcCell("!", nil, nil),
					},
				},
			},
		},
		{
			name:           "multiple lines",
			input:          "Hello,\nWorld!",
			expectedWidth:  6,
			expectedHeight: 2,
			expected: &Buffer{
				Lines: []Line{
					{
						newWcCell("H", nil, nil),
						newWcCell("e", nil, nil),
						newWcCell("l", nil, nil),
						newWcCell("l", nil, nil),
						newWcCell("o", nil, nil),
						newWcCell(",", nil, nil),
					},
					{
						newWcCell("W", nil, nil),
						newWcCell("o", nil, nil),
						newWcCell("r", nil, nil),
						newWcCell("l", nil, nil),
						newWcCell("d", nil, nil),
						newWcCell("!", nil, nil),
					},
				},
			},
		},
		{
			name:           "empty string",
			input:          "",
			expectedWidth:  0,
			expectedHeight: 1,
			expected: &Buffer{
				Lines: []Line{{}},
			},
		},
		{
			name:           "multiple lines different width",
			input:          "Hello,\nWorld!\nThis is a test.",
			expectedWidth:  15,
			expectedHeight: 3,
			expected: &Buffer{
				Lines: []Line{
					{
						newWcCell("H", nil, nil),
						newWcCell("e", nil, nil),
						newWcCell("l", nil, nil),
						newWcCell("l", nil, nil),
						newWcCell("o", nil, nil),
						newWcCell(",", nil, nil),
						newWcCell(" ", nil, nil),
						newWcCell(" ", nil, nil),
						newWcCell(" ", nil, nil),
						newWcCell(" ", nil, nil),
						newWcCell(" ", nil, nil),
						newWcCell(" ", nil, nil),
						newWcCell(" ", nil, nil),
						newWcCell(" ", nil, nil),
						newWcCell(" ", nil, nil),
					},
					{
						newWcCell("W", nil, nil),
						newWcCell("o", nil, nil),
						newWcCell("r", nil, nil),
						newWcCell("l", nil, nil),
						newWcCell("d", nil, nil),
						newWcCell("!", nil, nil),
						newWcCell(" ", nil, nil),
						newWcCell(" ", nil, nil),
						newWcCell(" ", nil, nil),
						newWcCell(" ", nil, nil),
						newWcCell(" ", nil, nil),
						newWcCell(" ", nil, nil),
						newWcCell(" ", nil, nil),
						newWcCell(" ", nil, nil),
						newWcCell(" ", nil, nil),
					},
					{
						newWcCell("T", nil, nil),
						newWcCell("h", nil, nil),
						newWcCell("i", nil, nil),
						newWcCell("s", nil, nil),
						newWcCell(" ", nil, nil),
						newWcCell("i", nil, nil),
						newWcCell("s", nil, nil),
						newWcCell(" ", nil, nil),
						newWcCell("a", nil, nil),
						newWcCell(" ", nil, nil),
						newWcCell("t", nil, nil),
						newWcCell("e", nil, nil),
						newWcCell("s", nil, nil),
						newWcCell("t", nil, nil),
						newWcCell(".", nil, nil),
					},
				},
			},
		},
		{
			name:           "unicode characters",
			input:          "Hello, 世界!",
			expectedWidth:  12,
			expectedHeight: 1,
			expected: &Buffer{
				Lines: []Line{
					{
						newWcCell("H", nil, nil),
						newWcCell("e", nil, nil),
						newWcCell("l", nil, nil),
						newWcCell("l", nil, nil),
						newWcCell("o", nil, nil),
						newWcCell(",", nil, nil),
						newWcCell(" ", nil, nil),
						newWcCell("世", nil, nil),
						Cell{},
						newWcCell("界", nil, nil),
						Cell{},
						newWcCell("!", nil, nil),
					},
				},
			},
		},
		{
			name:           "styled hello world string",
			input:          "\x1b[31;1;4mHello, \x1b[32;22;4mWorld!\x1b[0m",
			expectedWidth:  13,
			expectedHeight: 1,
			expected: &Buffer{
				Lines: []Line{
					{
						newWcCell("H", &Style{Fg: ansi.Red, Underline: UnderlineSingle, Attrs: AttrBold}, nil),
						newWcCell("e", &Style{Fg: ansi.Red, Underline: UnderlineSingle, Attrs: AttrBold}, nil),
						newWcCell("l", &Style{Fg: ansi.Red, Underline: UnderlineSingle, Attrs: AttrBold}, nil),
						newWcCell("l", &Style{Fg: ansi.Red, Underline: UnderlineSingle, Attrs: AttrBold}, nil),
						newWcCell("o", &Style{Fg: ansi.Red, Underline: UnderlineSingle, Attrs: AttrBold}, nil),
						newWcCell(",", &Style{Fg: ansi.Red, Underline: UnderlineSingle, Attrs: AttrBold}, nil),
						newWcCell(" ", &Style{Fg: ansi.Red, Underline: UnderlineSingle, Attrs: AttrBold}, nil),
						newWcCell("W", &Style{Fg: ansi.Green, Underline: UnderlineSingle}, nil),
						newWcCell("o", &Style{Fg: ansi.Green, Underline: UnderlineSingle}, nil),
						newWcCell("r", &Style{Fg: ansi.Green, Underline: UnderlineSingle}, nil),
						newWcCell("l", &Style{Fg: ansi.Green, Underline: UnderlineSingle}, nil),
						newWcCell("d", &Style{Fg: ansi.Green, Underline: UnderlineSingle}, nil),
						newWcCell("!", &Style{Fg: ansi.Green, Underline: UnderlineSingle}, nil),
					},
				},
			},
		},
		{
			name:           "complex styling with multiple SGR sequences",
			input:          "\x1b[31;1;2;4mR\x1b[22;1med\x1b[0m \x1b[32;3mGreen\x1b[0m \x1b[34;9mBlue\x1b[0m \x1b[33;7mYellow\x1b[0m \x1b[35;5mPurple\x1b[0m",
			expectedWidth:  28,
			expectedHeight: 1,
			expected: &Buffer{
				Lines: []Line{
					{
						newWcCell("R", &Style{Fg: ansi.Red, Underline: UnderlineSingle, Attrs: AttrBold | AttrFaint}, nil),
						newWcCell("e", &Style{Fg: ansi.Red, Underline: UnderlineSingle, Attrs: AttrBold}, nil),
						newWcCell("d", &Style{Fg: ansi.Red, Underline: UnderlineSingle, Attrs: AttrBold}, nil),
						newWcCell(" ", nil, nil),
						newWcCell("G", &Style{Fg: ansi.Green, Attrs: AttrItalic}, nil),
						newWcCell("r", &Style{Fg: ansi.Green, Attrs: AttrItalic}, nil),
						newWcCell("e", &Style{Fg: ansi.Green, Attrs: AttrItalic}, nil),
						newWcCell("e", &Style{Fg: ansi.Green, Attrs: AttrItalic}, nil),
						newWcCell("n", &Style{Fg: ansi.Green, Attrs: AttrItalic}, nil),
						newWcCell(" ", nil, nil),
						newWcCell("B", &Style{Fg: ansi.Blue, Attrs: AttrStrikethrough}, nil),
						newWcCell("l", &Style{Fg: ansi.Blue, Attrs: AttrStrikethrough}, nil),
						newWcCell("u", &Style{Fg: ansi.Blue, Attrs: AttrStrikethrough}, nil),
						newWcCell("e", &Style{Fg: ansi.Blue, Attrs: AttrStrikethrough}, nil),
						newWcCell(" ", nil, nil),
						newWcCell("Y", &Style{Fg: ansi.Yellow, Attrs: AttrReverse}, nil),
						newWcCell("e", &Style{Fg: ansi.Yellow, Attrs: AttrReverse}, nil),
						newWcCell("l", &Style{Fg: ansi.Yellow, Attrs: AttrReverse}, nil),
						newWcCell("l", &Style{Fg: ansi.Yellow, Attrs: AttrReverse}, nil),
						newWcCell("o", &Style{Fg: ansi.Yellow, Attrs: AttrReverse}, nil),
						newWcCell("w", &Style{Fg: ansi.Yellow, Attrs: AttrReverse}, nil),
						newWcCell(" ", nil, nil),
						newWcCell("P", &Style{Fg: ansi.Magenta, Attrs: AttrBlink}, nil),
						newWcCell("u", &Style{Fg: ansi.Magenta, Attrs: AttrBlink}, nil),
						newWcCell("r", &Style{Fg: ansi.Magenta, Attrs: AttrBlink}, nil),
						newWcCell("p", &Style{Fg: ansi.Magenta, Attrs: AttrBlink}, nil),
						newWcCell("l", &Style{Fg: ansi.Magenta, Attrs: AttrBlink}, nil),
						newWcCell("e", &Style{Fg: ansi.Magenta, Attrs: AttrBlink}, nil),
					},
				},
			},
		},
		{
			name:           "different underline styles",
			input:          "\x1b[4:1mSingle\x1b[0m \x1b[4:2mDouble\x1b[0m \x1b[4:3mCurly\x1b[0m \x1b[4:4mDotted\x1b[0m \x1b[4:5mDashed\x1b[0m",
			expectedWidth:  33,
			expectedHeight: 1,
			expected: &Buffer{
				Lines: []Line{
					{
						newWcCell("S", &Style{Underline: UnderlineSingle}, nil),
						newWcCell("i", &Style{Underline: UnderlineSingle}, nil),
						newWcCell("n", &Style{Underline: UnderlineSingle}, nil),
						newWcCell("g", &Style{Underline: UnderlineSingle}, nil),
						newWcCell("l", &Style{Underline: UnderlineSingle}, nil),
						newWcCell("e", &Style{Underline: UnderlineSingle}, nil),
						newWcCell(" ", nil, nil),
						newWcCell("D", &Style{Underline: UnderlineDouble}, nil),
						newWcCell("o", &Style{Underline: UnderlineDouble}, nil),
						newWcCell("u", &Style{Underline: UnderlineDouble}, nil),
						newWcCell("b", &Style{Underline: UnderlineDouble}, nil),
						newWcCell("l", &Style{Underline: UnderlineDouble}, nil),
						newWcCell("e", &Style{Underline: UnderlineDouble}, nil),
						newWcCell(" ", nil, nil),
						newWcCell("C", &Style{Underline: UnderlineCurly}, nil),
						newWcCell("u", &Style{Underline: UnderlineCurly}, nil),
						newWcCell("r", &Style{Underline: UnderlineCurly}, nil),
						newWcCell("l", &Style{Underline: UnderlineCurly}, nil),
						newWcCell("y", &Style{Underline: UnderlineCurly}, nil),
						newWcCell(" ", nil, nil),
						newWcCell("D", &Style{Underline: UnderlineDotted}, nil),
						newWcCell("o", &Style{Underline: UnderlineDotted}, nil),
						newWcCell("t", &Style{Underline: UnderlineDotted}, nil),
						newWcCell("t", &Style{Underline: UnderlineDotted}, nil),
						newWcCell("e", &Style{Underline: UnderlineDotted}, nil),
						newWcCell("d", &Style{Underline: UnderlineDotted}, nil),
						newWcCell(" ", nil, nil),
						newWcCell("D", &Style{Underline: UnderlineDashed}, nil),
						newWcCell("a", &Style{Underline: UnderlineDashed}, nil),
						newWcCell("s", &Style{Underline: UnderlineDashed}, nil),
						newWcCell("h", &Style{Underline: UnderlineDashed}, nil),
						newWcCell("e", &Style{Underline: UnderlineDashed}, nil),
						newWcCell("d", &Style{Underline: UnderlineDashed}, nil),
					},
				},
			},
		},
		{
			name:           "truecolor and 256 color support",
			input:          "\x1b[38;2;255;0;0mRGB Red\x1b[0m \x1b[48;2;0;255;0mRGB Green BG\x1b[0m \x1b[38;5;33m256 Blue\x1b[0m",
			expectedWidth:  29,
			expectedHeight: 1,
			expected: &Buffer{
				Lines: []Line{
					{
						newWcCell("R", &Style{Fg: color.RGBA{R: 255, G: 0, B: 0, A: 255}}, nil),
						newWcCell("G", &Style{Fg: color.RGBA{R: 255, G: 0, B: 0, A: 255}}, nil),
						newWcCell("B", &Style{Fg: color.RGBA{R: 255, G: 0, B: 0, A: 255}}, nil),
						newWcCell(" ", &Style{Fg: color.RGBA{R: 255, G: 0, B: 0, A: 255}}, nil),
						newWcCell("R", &Style{Fg: color.RGBA{R: 255, G: 0, B: 0, A: 255}}, nil),
						newWcCell("e", &Style{Fg: color.RGBA{R: 255, G: 0, B: 0, A: 255}}, nil),
						newWcCell("d", &Style{Fg: color.RGBA{R: 255, G: 0, B: 0, A: 255}}, nil),
						newWcCell(" ", nil, nil),
						newWcCell("R", &Style{Bg: color.RGBA{R: 0, G: 255, B: 0, A: 255}}, nil),
						newWcCell("G", &Style{Bg: color.RGBA{R: 0, G: 255, B: 0, A: 255}}, nil),
						newWcCell("B", &Style{Bg: color.RGBA{R: 0, G: 255, B: 0, A: 255}}, nil),
						newWcCell(" ", &Style{Bg: color.RGBA{R: 0, G: 255, B: 0, A: 255}}, nil),
						newWcCell("G", &Style{Bg: color.RGBA{R: 0, G: 255, B: 0, A: 255}}, nil),
						newWcCell("r", &Style{Bg: color.RGBA{R: 0, G: 255, B: 0, A: 255}}, nil),
						newWcCell("e", &Style{Bg: color.RGBA{R: 0, G: 255, B: 0, A: 255}}, nil),
						newWcCell("e", &Style{Bg: color.RGBA{R: 0, G: 255, B: 0, A: 255}}, nil),
						newWcCell("n", &Style{Bg: color.RGBA{R: 0, G: 255, B: 0, A: 255}}, nil),
						newWcCell(" ", &Style{Bg: color.RGBA{R: 0, G: 255, B: 0, A: 255}}, nil),
						newWcCell("B", &Style{Bg: color.RGBA{R: 0, G: 255, B: 0, A: 255}}, nil),
						newWcCell("G", &Style{Bg: color.RGBA{R: 0, G: 255, B: 0, A: 255}}, nil),
						newWcCell(" ", nil, nil),
						newWcCell("2", &Style{Fg: ansi.IndexedColor(33)}, nil),
						newWcCell("5", &Style{Fg: ansi.IndexedColor(33)}, nil),
						newWcCell("6", &Style{Fg: ansi.IndexedColor(33)}, nil),
						newWcCell(" ", &Style{Fg: ansi.IndexedColor(33)}, nil),
						newWcCell("B", &Style{Fg: ansi.IndexedColor(33)}, nil),
						newWcCell("l", &Style{Fg: ansi.IndexedColor(33)}, nil),
						newWcCell("u", &Style{Fg: ansi.IndexedColor(33)}, nil),
						newWcCell("e", &Style{Fg: ansi.IndexedColor(33)}, nil),
						newWcCell("e", &Style{Fg: ansi.IndexedColor(33)}, nil),
					},
				},
			},
		},
		{
			name:           "hyperlink support",
			input:          "Normal \x1b]8;;https://charm.sh\x1b\\Charm\x1b]8;;\x1b\\ Text \x1b]8;;https://github.com/charmbracelet\x1b\\GitHub\x1b]8;;\x1b\\",
			expectedWidth:  24,
			expectedHeight: 1,
			expected: &Buffer{
				Lines: []Line{
					{
						newWcCell("N", nil, nil),
						newWcCell("o", nil, nil),
						newWcCell("r", nil, nil),
						newWcCell("m", nil, nil),
						newWcCell("a", nil, nil),
						newWcCell("l", nil, nil),
						newWcCell(" ", nil, nil),
						newWcCell("C", nil, &Link{URL: "https://charm.sh"}),
						newWcCell("h", nil, &Link{URL: "https://charm.sh"}),
						newWcCell("a", nil, &Link{URL: "https://charm.sh"}),
						newWcCell("r", nil, &Link{URL: "https://charm.sh"}),
						newWcCell("m", nil, &Link{URL: "https://charm.sh"}),
						newWcCell(" ", nil, nil),
						newWcCell("T", nil, nil),
						newWcCell("e", nil, nil),
						newWcCell("x", nil, nil),
						newWcCell("t", nil, nil),
						newWcCell(" ", nil, nil),
						newWcCell("G", nil, &Link{URL: "https://github.com/charmbracelet"}),
						newWcCell("i", nil, &Link{URL: "https://github.com/charmbracelet"}),
						newWcCell("t", nil, &Link{URL: "https://github.com/charmbracelet"}),
						newWcCell("H", nil, &Link{URL: "https://github.com/charmbracelet"}),
						newWcCell("u", nil, &Link{URL: "https://github.com/charmbracelet"}),
						newWcCell("b", nil, &Link{URL: "https://github.com/charmbracelet"}),
					},
				},
			},
		},
		{
			name:           "complex mixed styling with hyperlinks",
			input:          "\x1b[31;1;2;3mR\x1b[22;23;1med \x1b]8;;https://charm.sh\x1b\\\x1b[4mCharm\x1b]8;;\x1b\\\x1b[0m \x1b[38;5;33;48;2;0;100;0m\x1b]8;;https://github.com\x1b\\GitHub\x1b]8;;\x1b\\\x1b[0m",
			expectedWidth:  16,
			expectedHeight: 1,
			expected: &Buffer{
				Lines: []Line{
					{
						newWcCell("R", &Style{Fg: ansi.Red, Attrs: AttrBold | AttrFaint | AttrItalic}, nil),
						newWcCell("e", &Style{Fg: ansi.Red, Attrs: AttrBold}, nil),
						newWcCell("d", &Style{Fg: ansi.Red, Attrs: AttrBold}, nil),
						newWcCell(" ", &Style{Fg: ansi.Red, Attrs: AttrBold}, nil),
						newWcCell("C", &Style{Fg: ansi.Red, Underline: UnderlineSingle, Attrs: AttrBold}, &Link{URL: "https://charm.sh"}),
						newWcCell("h", &Style{Fg: ansi.Red, Underline: UnderlineSingle, Attrs: AttrBold}, &Link{URL: "https://charm.sh"}),
						newWcCell("a", &Style{Fg: ansi.Red, Underline: UnderlineSingle, Attrs: AttrBold}, &Link{URL: "https://charm.sh"}),
						newWcCell("r", &Style{Fg: ansi.Red, Underline: UnderlineSingle, Attrs: AttrBold}, &Link{URL: "https://charm.sh"}),
						newWcCell("m", &Style{Fg: ansi.Red, Underline: UnderlineSingle, Attrs: AttrBold}, &Link{URL: "https://charm.sh"}),
						newWcCell(" ", nil, nil),
						newWcCell("G", &Style{Fg: ansi.IndexedColor(33), Bg: color.RGBA{R: 0, G: 100, B: 0, A: 255}}, &Link{URL: "https://github.com"}),
						newWcCell("i", &Style{Fg: ansi.IndexedColor(33), Bg: color.RGBA{R: 0, G: 100, B: 0, A: 255}}, &Link{URL: "https://github.com"}),
						newWcCell("t", &Style{Fg: ansi.IndexedColor(33), Bg: color.RGBA{R: 0, G: 100, B: 0, A: 255}}, &Link{URL: "https://github.com"}),
						newWcCell("H", &Style{Fg: ansi.IndexedColor(33), Bg: color.RGBA{R: 0, G: 100, B: 0, A: 255}}, &Link{URL: "https://github.com"}),
						newWcCell("u", &Style{Fg: ansi.IndexedColor(33), Bg: color.RGBA{R: 0, G: 100, B: 0, A: 255}}, &Link{URL: "https://github.com"}),
						newWcCell("b", &Style{Fg: ansi.IndexedColor(33), Bg: color.RGBA{R: 0, G: 100, B: 0, A: 255}}, &Link{URL: "https://github.com"}),
					},
				},
			},
		},
	}

	for i, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Logf("Running case %d: %s for %q", i+1, tc.name, tc.input)
			ss := NewStyledString(tc.input)
			area := ss.Bounds()
			buf := NewScreenBuffer(area.Dx(), area.Dy())
			ss.Draw(buf, area)
			if buf.Width() != tc.expectedWidth {
				t.Errorf("case %d expected width %d, got %d", i+1, tc.expectedWidth, buf.Width())
			}
			if buf.Height() != tc.expectedHeight {
				t.Errorf("case %d expected height %d, got %d", i+1, tc.expectedHeight, buf.Height())
			}
			for y, line := range buf.Lines {
				for x, cell := range line {
					if !cellEqual(tc.expected.CellAt(x, y), &cell) {
						t.Errorf("case %d expected cell (%d, %d) %#v, got %#v", y+1, x, y, tc.expected.CellAt(x, y), &cell)
					}
				}
			}
		})
	}
}

func TestStyledStringEmptyLines(t *testing.T) {
	// This test uses an input that results in empty lines when drawn to a smaller
	// screen buffer.
	input := "\x1b[31;1;4mHello, \x1b[32;22;4mWorld!\x1b[0m"
	ss := NewStyledString(input)
	scr := NewScreenBuffer(5, 3)
	ss.Draw(scr, scr.Bounds())
	expected := &Buffer{
		Lines: []Line{
			{
				newWcCell("H", &Style{Fg: ansi.Red, Underline: UnderlineSingle, Attrs: AttrBold}, nil),
				newWcCell("e", &Style{Fg: ansi.Red, Underline: UnderlineSingle, Attrs: AttrBold}, nil),
				newWcCell("l", &Style{Fg: ansi.Red, Underline: UnderlineSingle, Attrs: AttrBold}, nil),
				newWcCell("l", &Style{Fg: ansi.Red, Underline: UnderlineSingle, Attrs: AttrBold}, nil),
				newWcCell("o", &Style{Fg: ansi.Red, Underline: UnderlineSingle, Attrs: AttrBold}, nil),
			},
			NewLine(5),
			NewLine(5),
		},
	}
	for y, line := range scr.Lines {
		for x, cell := range line {
			if !cellEqual(expected.CellAt(x, y), &cell) {
				t.Errorf("expected cell (%d, %d) %#v, got %#v", x, y, expected.CellAt(x, y), &cell)
			}
		}
	}
}

func newWcCell(s string, style *Style, link *Link) Cell {
	c := NewCell(ansi.WcWidth, s)
	if style != nil {
		c.Style = *style
	}
	if link != nil {
		c.Link = *link
	}
	return *c
}
