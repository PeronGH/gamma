package screen

import (
	"fmt"
	"testing"

	gamma "github.com/PeronGH/gamma"
	"github.com/charmbracelet/x/ansi"
)

// mockScreen is a basic screen implementation for testing
type mockScreen struct {
	*gamma.Buffer
	method ansi.Method
}

func newMockScreen(width, height int) *mockScreen {
	return &mockScreen{
		Buffer: gamma.NewBuffer(width, height),
		method: ansi.WcWidth,
	}
}

func (m *mockScreen) Bounds() gamma.Rectangle {
	return m.Buffer.Bounds()
}

func (m *mockScreen) CellAt(x, y int) *gamma.Cell {
	return m.Buffer.CellAt(x, y)
}

func (m *mockScreen) SetCell(x, y int, c *gamma.Cell) {
	m.Buffer.SetCell(x, y, c)
}

func (m *mockScreen) WidthMethod() gamma.WidthMethod {
	return m.method
}

// minimalMockScreen is a minimal screen that only implements the required Screen interface
// and doesn't have any of the optional methods (Clear, Fill, Clone, etc.)
type minimalMockScreen struct {
	cells  [][]gamma.Cell
	width  int
	height int
	method ansi.Method
}

func newMinimalMockScreen(width, height int) *minimalMockScreen {
	cells := make([][]gamma.Cell, height)
	for i := range cells {
		cells[i] = make([]gamma.Cell, width)
		for j := range cells[i] {
			cells[i][j] = gamma.EmptyCell
		}
	}
	return &minimalMockScreen{
		cells:  cells,
		width:  width,
		height: height,
		method: ansi.WcWidth,
	}
}

func (m *minimalMockScreen) Bounds() gamma.Rectangle {
	return gamma.Rect(0, 0, m.width, m.height)
}

func (m *minimalMockScreen) CellAt(x, y int) *gamma.Cell {
	if x < 0 || x >= m.width || y < 0 || y >= m.height {
		return nil
	}
	return &m.cells[y][x]
}

func (m *minimalMockScreen) SetCell(x, y int, c *gamma.Cell) {
	if x < 0 || x >= m.width || y < 0 || y >= m.height {
		return
	}
	if c == nil {
		m.cells[y][x] = gamma.EmptyCell
	} else {
		m.cells[y][x] = *c
	}
}

func (m *minimalMockScreen) WidthMethod() gamma.WidthMethod {
	return m.method
}

// nilCellMockScreen is a screen that can return nil for certain cells
type nilCellMockScreen struct {
	cells        map[string]*gamma.Cell
	width        int
	height       int
	method       ansi.Method
	nilPositions map[string]bool
}

func newNilCellMockScreen(width, height int) *nilCellMockScreen {
	return &nilCellMockScreen{
		cells:        make(map[string]*gamma.Cell),
		width:        width,
		height:       height,
		method:       ansi.WcWidth,
		nilPositions: make(map[string]bool),
	}
}

func (m *nilCellMockScreen) Bounds() gamma.Rectangle {
	return gamma.Rect(0, 0, m.width, m.height)
}

func (m *nilCellMockScreen) CellAt(x, y int) *gamma.Cell {
	if x < 0 || x >= m.width || y < 0 || y >= m.height {
		return nil
	}
	key := fmt.Sprintf("%d,%d", x, y)
	if m.nilPositions[key] {
		return nil
	}
	if cell, ok := m.cells[key]; ok {
		return cell
	}
	// Return empty cell by default
	empty := gamma.EmptyCell
	return &empty
}

func (m *nilCellMockScreen) SetCell(x, y int, c *gamma.Cell) {
	if x < 0 || x >= m.width || y < 0 || y >= m.height {
		return
	}
	key := fmt.Sprintf("%d,%d", x, y)
	if c == nil {
		m.cells[key] = &gamma.EmptyCell
	} else {
		cellCopy := *c
		m.cells[key] = &cellCopy
	}
}

func (m *nilCellMockScreen) WidthMethod() gamma.WidthMethod {
	return m.method
}

func (m *nilCellMockScreen) SetNilAt(x, y int) {
	key := fmt.Sprintf("%d,%d", x, y)
	m.nilPositions[key] = true
}

// mockScreenWithClear implements Clear method
type mockScreenWithClear struct {
	*mockScreen
	clearCalled bool
}

func (m *mockScreenWithClear) Clear() {
	m.clearCalled = true
	m.Buffer.Clear()
}

// mockScreenWithClearArea implements ClearArea method
type mockScreenWithClearArea struct {
	*mockScreen
	clearAreaCalled bool
	lastArea        gamma.Rectangle
}

func (m *mockScreenWithClearArea) ClearArea(area gamma.Rectangle) {
	m.clearAreaCalled = true
	m.lastArea = area
	m.Buffer.ClearArea(area)
}

// mockScreenWithFill implements Fill method
type mockScreenWithFill struct {
	*mockScreen
	fillCalled bool
	lastCell   *gamma.Cell
}

func (m *mockScreenWithFill) Fill(cell *gamma.Cell) {
	m.fillCalled = true
	m.lastCell = cell
	m.Buffer.Fill(cell)
}

// mockScreenWithFillArea implements FillArea method
type mockScreenWithFillArea struct {
	*mockScreen
	fillAreaCalled bool
	lastCell       *gamma.Cell
	lastArea       gamma.Rectangle
}

func (m *mockScreenWithFillArea) FillArea(cell *gamma.Cell, area gamma.Rectangle) {
	m.fillAreaCalled = true
	m.lastCell = cell
	m.lastArea = area
	m.Buffer.FillArea(cell, area)
}

// mockScreenWithClone implements Clone method
type mockScreenWithClone struct {
	*mockScreen
	cloneCalled bool
}

func (m *mockScreenWithClone) Clone() *gamma.Buffer {
	m.cloneCalled = true
	return m.Buffer.Clone()
}

// mockScreenWithCloneArea implements CloneArea method
type mockScreenWithCloneArea struct {
	*mockScreen
	cloneAreaCalled bool
	lastArea        gamma.Rectangle
}

func (m *mockScreenWithCloneArea) CloneArea(area gamma.Rectangle) *gamma.Buffer {
	m.cloneAreaCalled = true
	m.lastArea = area
	return m.Buffer.CloneArea(area)
}

func TestClear(t *testing.T) {
	t.Run("with Clear method", func(t *testing.T) {
		scr := &mockScreenWithClear{
			mockScreen: newMockScreen(10, 5),
		}

		// Set some cells
		testCell := &gamma.Cell{Content: "X", Width: 1}
		scr.SetCell(0, 0, testCell)
		scr.SetCell(5, 2, testCell)

		Clear(scr)

		if !scr.clearCalled {
			t.Error("Clear method was not called")
		}

		// Verify cells are cleared
		if cell := scr.CellAt(0, 0); cell != nil && cell.Content != " " {
			t.Errorf("Cell at (0,0) not cleared, got %v", cell)
		}
		if cell := scr.CellAt(5, 2); cell != nil && cell.Content != " " {
			t.Errorf("Cell at (5,2) not cleared, got %v", cell)
		}
	})

	t.Run("without Clear method", func(t *testing.T) {
		scr := newMockScreen(10, 5)

		// Set some cells
		testCell := &gamma.Cell{Content: "X", Width: 1}
		scr.SetCell(0, 0, testCell)
		scr.SetCell(5, 2, testCell)

		Clear(scr)

		// Verify cells are cleared (filled with nil/empty)
		if cell := scr.CellAt(0, 0); cell != nil && cell.Content != " " {
			t.Errorf("Cell at (0,0) not cleared, got %v", cell)
		}
		if cell := scr.CellAt(5, 2); cell != nil && cell.Content != " " {
			t.Errorf("Cell at (5,2) not cleared, got %v", cell)
		}
	})
}

func TestClearArea(t *testing.T) {
	t.Run("with ClearArea method", func(t *testing.T) {
		scr := &mockScreenWithClearArea{
			mockScreen: newMockScreen(10, 5),
		}

		// Set cells across the screen
		testCell := &gamma.Cell{Content: "X", Width: 1}
		for y := 0; y < 5; y++ {
			for x := 0; x < 10; x++ {
				scr.SetCell(x, y, testCell)
			}
		}

		area := gamma.Rect(2, 1, 4, 2)
		ClearArea(scr, area)

		if !scr.clearAreaCalled {
			t.Error("ClearArea method was not called")
		}

		if scr.lastArea != area {
			t.Errorf("ClearArea called with wrong area, expected %v, got %v", area, scr.lastArea)
		}

		// Verify only the area is cleared
		for y := 0; y < 5; y++ {
			for x := 0; x < 10; x++ {
				cell := scr.CellAt(x, y)
				if x >= 2 && x < 6 && y >= 1 && y < 3 {
					// Inside cleared area
					if cell != nil && cell.Content != " " {
						t.Errorf("Cell at (%d,%d) should be cleared, got %v", x, y, cell)
					}
				} else {
					// Outside cleared area
					if cell == nil || cell.Content != "X" {
						t.Errorf("Cell at (%d,%d) should not be cleared, got %v", x, y, cell)
					}
				}
			}
		}
	})

	t.Run("without ClearArea method", func(t *testing.T) {
		scr := newMockScreen(10, 5)

		// Set cells across the screen
		testCell := &gamma.Cell{Content: "X", Width: 1}
		for y := 0; y < 5; y++ {
			for x := 0; x < 10; x++ {
				scr.SetCell(x, y, testCell)
			}
		}

		area := gamma.Rect(2, 1, 4, 2)
		ClearArea(scr, area)

		// Verify only the area is cleared
		for y := 0; y < 5; y++ {
			for x := 0; x < 10; x++ {
				cell := scr.CellAt(x, y)
				if x >= 2 && x < 6 && y >= 1 && y < 3 {
					// Inside cleared area
					if cell != nil && cell.Content != " " {
						t.Errorf("Cell at (%d,%d) should be cleared, got %v", x, y, cell)
					}
				} else {
					// Outside cleared area
					if cell == nil || cell.Content != "X" {
						t.Errorf("Cell at (%d,%d) should not be cleared, got %v", x, y, cell)
					}
				}
			}
		}
	})
}

func TestFill(t *testing.T) {
	t.Run("with Fill method", func(t *testing.T) {
		scr := &mockScreenWithFill{
			mockScreen: newMockScreen(10, 5),
		}

		fillCell := &gamma.Cell{Content: "F", Width: 1}
		Fill(scr, fillCell)

		if !scr.fillCalled {
			t.Error("Fill method was not called")
		}

		if scr.lastCell != fillCell {
			t.Errorf("Fill called with wrong cell, expected %v, got %v", fillCell, scr.lastCell)
		}

		// Verify all cells are filled
		for y := 0; y < 5; y++ {
			for x := 0; x < 10; x++ {
				cell := scr.CellAt(x, y)
				if cell == nil || cell.Content != "F" {
					t.Errorf("Cell at (%d,%d) not filled correctly, got %v", x, y, cell)
				}
			}
		}
	})

	t.Run("without Fill method", func(t *testing.T) {
		scr := newMockScreen(10, 5)

		fillCell := &gamma.Cell{Content: "F", Width: 1}
		Fill(scr, fillCell)

		// Verify all cells are filled
		for y := 0; y < 5; y++ {
			for x := 0; x < 10; x++ {
				cell := scr.CellAt(x, y)
				if cell == nil || cell.Content != "F" {
					t.Errorf("Cell at (%d,%d) not filled correctly, got %v", x, y, cell)
				}
			}
		}
	})

	t.Run("with nil cell", func(t *testing.T) {
		scr := newMockScreen(10, 5)

		// Set some cells first
		testCell := &gamma.Cell{Content: "X", Width: 1}
		scr.SetCell(0, 0, testCell)
		scr.SetCell(5, 2, testCell)

		Fill(scr, nil)

		// Verify all cells are cleared (filled with empty)
		for y := 0; y < 5; y++ {
			for x := 0; x < 10; x++ {
				cell := scr.CellAt(x, y)
				if cell != nil && cell.Content != " " {
					t.Errorf("Cell at (%d,%d) not cleared, got %v", x, y, cell)
				}
			}
		}
	})
}

func TestFillArea(t *testing.T) {
	t.Run("with FillArea method", func(t *testing.T) {
		scr := &mockScreenWithFillArea{
			mockScreen: newMockScreen(10, 5),
		}

		fillCell := &gamma.Cell{Content: "A", Width: 1}
		area := gamma.Rect(2, 1, 4, 2)
		FillArea(scr, fillCell, area)

		if !scr.fillAreaCalled {
			t.Error("FillArea method was not called")
		}

		if scr.lastCell != fillCell {
			t.Errorf("FillArea called with wrong cell, expected %v, got %v", fillCell, scr.lastCell)
		}

		if scr.lastArea != area {
			t.Errorf("FillArea called with wrong area, expected %v, got %v", area, scr.lastArea)
		}

		// Verify only the area is filled
		for y := 0; y < 5; y++ {
			for x := 0; x < 10; x++ {
				cell := scr.CellAt(x, y)
				if x >= 2 && x < 6 && y >= 1 && y < 3 {
					// Inside filled area
					if cell == nil || cell.Content != "A" {
						t.Errorf("Cell at (%d,%d) should be filled with 'A', got %v", x, y, cell)
					}
				}
			}
		}
	})

	t.Run("with FillArea method and wide cell", func(t *testing.T) {
		scr := &mockScreenWithFillArea{
			mockScreen: newMockScreen(10, 5),
		}

		fillCell := &gamma.Cell{Content: "混", Width: 2}
		area := gamma.Rect(0, 1, 4, 2)
		FillArea(scr, fillCell, area)

		if !scr.fillAreaCalled {
			t.Error("FillArea method was not called")
		}

		if scr.lastCell != fillCell {
			t.Errorf("FillArea called with wrong cell, expected %v, got %v", fillCell, scr.lastCell)
		}

		if scr.lastArea != area {
			t.Errorf("FillArea called with wrong area, expected %v, got %v", area, scr.lastArea)
		}

		// Verify only the area is filled with the wide cell
		for y := 0; y < 5; y++ {
			for x := 0; x < 10; x += fillCell.Width {
				cell := scr.CellAt(x, y)
				if x >= 0 && x < 4 && y >= 1 && y < 3 {
					// Inside filled area
					if cell == nil || cell.Content != "混" || cell.Width != 2 {
						t.Errorf("Cell at (%d,%d) should be filled with '混', got %q", x, y, cell)
					}
				}
			}
		}
	})

	t.Run("without FillArea method", func(t *testing.T) {
		scr := newMockScreen(10, 5)

		fillCell := &gamma.Cell{Content: "B", Width: 1}
		area := gamma.Rect(2, 1, 4, 2)
		FillArea(scr, fillCell, area)

		// Verify only the area is filled
		for y := 0; y < 5; y++ {
			for x := 0; x < 10; x++ {
				cell := scr.CellAt(x, y)
				if x >= 2 && x < 6 && y >= 1 && y < 3 {
					// Inside filled area
					if cell == nil || cell.Content != "B" {
						t.Errorf("Cell at (%d,%d) should be filled with 'B', got %v", x, y, cell)
					}
				}
			}
		}
	})

	t.Run("with nil cell", func(t *testing.T) {
		scr := newMockScreen(10, 5)

		// Set cells across the screen
		testCell := &gamma.Cell{Content: "X", Width: 1}
		for y := 0; y < 5; y++ {
			for x := 0; x < 10; x++ {
				scr.SetCell(x, y, testCell)
			}
		}

		area := gamma.Rect(2, 1, 4, 2)
		FillArea(scr, nil, area)

		// Verify only the area is cleared
		for y := 0; y < 5; y++ {
			for x := 0; x < 10; x++ {
				cell := scr.CellAt(x, y)
				if x >= 2 && x < 6 && y >= 1 && y < 3 {
					// Inside filled area (should be empty)
					if cell != nil && cell.Content != " " {
						t.Errorf("Cell at (%d,%d) should be cleared, got %v", x, y, cell)
					}
				} else {
					// Outside filled area
					if cell == nil || cell.Content != "X" {
						t.Errorf("Cell at (%d,%d) should not be changed, got %v", x, y, cell)
					}
				}
			}
		}
	})
}

func TestClone(t *testing.T) {
	t.Run("with Clone method", func(t *testing.T) {
		scr := &mockScreenWithClone{
			mockScreen: newMockScreen(10, 5),
		}

		// Set some cells
		scr.SetCell(0, 0, &gamma.Cell{Content: "A", Width: 1})
		scr.SetCell(5, 2, &gamma.Cell{Content: "B", Width: 1})
		scr.SetCell(9, 4, &gamma.Cell{Content: "C", Width: 1})

		cloned := Clone(scr)

		if !scr.cloneCalled {
			t.Error("Clone method was not called")
		}

		if cloned == nil {
			t.Fatal("Clone returned nil")
		}

		// Verify cloned buffer has same content
		if cloned.Width() != 10 || cloned.Height() != 5 {
			t.Errorf("Cloned buffer has wrong dimensions: %dx%d", cloned.Width(), cloned.Height())
		}

		// Check specific cells
		if cell := cloned.CellAt(0, 0); cell == nil || cell.Content != "A" {
			t.Errorf("Cloned cell at (0,0) incorrect, got %v", cell)
		}
		if cell := cloned.CellAt(5, 2); cell == nil || cell.Content != "B" {
			t.Errorf("Cloned cell at (5,2) incorrect, got %v", cell)
		}
		if cell := cloned.CellAt(9, 4); cell == nil || cell.Content != "C" {
			t.Errorf("Cloned cell at (9,4) incorrect, got %v", cell)
		}
	})

	t.Run("without Clone method", func(t *testing.T) {
		scr := newMockScreen(10, 5)

		// Set some cells
		scr.SetCell(0, 0, &gamma.Cell{Content: "A", Width: 1})
		scr.SetCell(5, 2, &gamma.Cell{Content: "B", Width: 1})
		scr.SetCell(9, 4, &gamma.Cell{Content: "C", Width: 1})

		cloned := Clone(scr)

		if cloned == nil {
			t.Fatal("Clone returned nil")
		}

		// Verify cloned buffer has same content
		if cloned.Width() != 10 || cloned.Height() != 5 {
			t.Errorf("Cloned buffer has wrong dimensions: %dx%d", cloned.Width(), cloned.Height())
		}

		// Check specific cells
		if cell := cloned.CellAt(0, 0); cell == nil || cell.Content != "A" {
			t.Errorf("Cloned cell at (0,0) incorrect, got %v", cell)
		}
		if cell := cloned.CellAt(5, 2); cell == nil || cell.Content != "B" {
			t.Errorf("Cloned cell at (5,2) incorrect, got %v", cell)
		}
		if cell := cloned.CellAt(9, 4); cell == nil || cell.Content != "C" {
			t.Errorf("Cloned cell at (9,4) incorrect, got %v", cell)
		}
	})
}

func TestCloneArea(t *testing.T) {
	t.Run("with CloneArea method", func(t *testing.T) {
		scr := &mockScreenWithCloneArea{
			mockScreen: newMockScreen(10, 5),
		}

		// Set cells across the screen
		for y := 0; y < 5; y++ {
			for x := 0; x < 10; x++ {
				content := string(rune('A' + y*10 + x))
				scr.SetCell(x, y, &gamma.Cell{Content: content, Width: 1})
			}
		}

		area := gamma.Rect(2, 1, 4, 2)
		cloned := CloneArea(scr, area)

		if !scr.cloneAreaCalled {
			t.Error("CloneArea method was not called")
		}

		if scr.lastArea != area {
			t.Errorf("CloneArea called with wrong area, expected %v, got %v", area, scr.lastArea)
		}

		if cloned == nil {
			t.Fatal("CloneArea returned nil")
		}

		// Verify cloned buffer has correct dimensions
		if cloned.Width() != 4 || cloned.Height() != 2 {
			t.Errorf("Cloned buffer has wrong dimensions: %dx%d, expected 4x2", cloned.Width(), cloned.Height())
		}

		// Verify content matches the cloned area
		for y := 0; y < 2; y++ {
			for x := 0; x < 4; x++ {
				expectedContent := string(rune('A' + (y+1)*10 + (x + 2)))
				cell := cloned.CellAt(x, y)
				if cell == nil || cell.Content != expectedContent {
					t.Errorf("Cloned cell at (%d,%d) incorrect, expected %s, got %v", x, y, expectedContent, cell)
				}
			}
		}
	})

	t.Run("without CloneArea method", func(t *testing.T) {
		scr := newMockScreen(10, 5)

		// Set cells across the screen
		for y := 0; y < 5; y++ {
			for x := 0; x < 10; x++ {
				content := string(rune('A' + y*10 + x))
				scr.SetCell(x, y, &gamma.Cell{Content: content, Width: 1})
			}
		}

		area := gamma.Rect(2, 1, 4, 2)
		cloned := CloneArea(scr, area)

		if cloned == nil {
			t.Fatal("CloneArea returned nil")
		}

		// Verify cloned buffer has correct dimensions
		if cloned.Width() != 4 || cloned.Height() != 2 {
			t.Errorf("Cloned buffer has wrong dimensions: %dx%d, expected 4x2", cloned.Width(), cloned.Height())
		}

		// Verify content matches the cloned area
		for y := 0; y < 2; y++ {
			for x := 0; x < 4; x++ {
				expectedContent := string(rune('A' + (y+1)*10 + (x + 2)))
				cell := cloned.CellAt(x, y)
				if cell == nil || cell.Content != expectedContent {
					t.Errorf("Cloned cell at (%d,%d) incorrect, expected %s, got %v", x, y, expectedContent, cell)
				}
			}
		}
	})

	t.Run("with zero cells", func(t *testing.T) {
		scr := newMockScreen(10, 5)

		// Set some cells but leave some as zero
		scr.SetCell(2, 1, &gamma.Cell{Content: "A", Width: 1})
		scr.SetCell(4, 2, &gamma.Cell{Content: "B", Width: 1})
		// Leave (3,1) and (5,2) as zero cells

		area := gamma.Rect(2, 1, 4, 2)
		cloned := CloneArea(scr, area)

		if cloned == nil {
			t.Fatal("CloneArea returned nil")
		}

		// Check that non-zero cells are cloned
		if cell := cloned.CellAt(0, 0); cell == nil || cell.Content != "A" {
			t.Errorf("Cell at (0,0) should be 'A', got %v", cell)
		}
		if cell := cloned.CellAt(2, 1); cell == nil || cell.Content != "B" {
			t.Errorf("Cell at (2,1) should be 'B', got %v", cell)
		}

		// Check that zero cells are not cloned (should be empty)
		if cell := cloned.CellAt(1, 0); cell != nil && !cell.IsZero() && cell.Content != " " {
			t.Errorf("Cell at (1,0) should be zero or empty, got %v", cell)
		}
		if cell := cloned.CellAt(3, 1); cell != nil && !cell.IsZero() && cell.Content != " " {
			t.Errorf("Cell at (3,1) should be zero or empty, got %v", cell)
		}
	})
}

func TestScreenBufferIntegration(t *testing.T) {
	t.Run("using ScreenBuffer", func(t *testing.T) {
		scr := gamma.NewScreenBuffer(10, 5)

		// Test Clear
		scr.SetCell(0, 0, &gamma.Cell{Content: "X", Width: 1})
		Clear(scr)
		if cell := scr.CellAt(0, 0); cell != nil && cell.Content != " " {
			t.Errorf("Clear failed, cell at (0,0) is %v", cell)
		}

		// Test Fill
		fillCell := &gamma.Cell{Content: "F", Width: 1}
		Fill(scr, fillCell)
		for y := 0; y < 5; y++ {
			for x := 0; x < 10; x++ {
				if cell := scr.CellAt(x, y); cell == nil || cell.Content != "F" {
					t.Errorf("Fill failed at (%d,%d), got %v", x, y, cell)
				}
			}
		}

		// Test ClearArea
		area := gamma.Rect(2, 1, 3, 2)
		ClearArea(scr, area)
		for y := 1; y < 3; y++ {
			for x := 2; x < 5; x++ {
				if cell := scr.CellAt(x, y); cell != nil && cell.Content != " " {
					t.Errorf("ClearArea failed at (%d,%d), got %v", x, y, cell)
				}
			}
		}

		// Test FillArea
		fillCell2 := &gamma.Cell{Content: "A", Width: 1}
		area2 := gamma.Rect(1, 1, 2, 2)
		FillArea(scr, fillCell2, area2)
		for y := 1; y < 3; y++ {
			for x := 1; x < 3; x++ {
				if cell := scr.CellAt(x, y); cell == nil || cell.Content != "A" {
					t.Errorf("FillArea failed at (%d,%d), got %v", x, y, cell)
				}
			}
		}

		// Test Clone
		cloned := Clone(scr)
		if cloned.Width() != 10 || cloned.Height() != 5 {
			t.Errorf("Clone dimensions wrong: %dx%d", cloned.Width(), cloned.Height())
		}

		// Test CloneArea
		area3 := gamma.Rect(0, 0, 3, 3)
		clonedArea := CloneArea(scr, area3)
		if clonedArea.Width() != 3 || clonedArea.Height() != 3 {
			t.Errorf("CloneArea dimensions wrong: %dx%d", clonedArea.Width(), clonedArea.Height())
		}
	})
}

func TestEdgeCases(t *testing.T) {
	t.Run("empty screen", func(t *testing.T) {
		scr := newMockScreen(0, 0)

		// These should not panic
		Clear(scr)
		Fill(scr, &gamma.Cell{Content: "X", Width: 1})
		ClearArea(scr, gamma.Rect(0, 0, 1, 1))
		FillArea(scr, &gamma.Cell{Content: "X", Width: 1}, gamma.Rect(0, 0, 1, 1))
		Clone(scr)
		CloneArea(scr, gamma.Rect(0, 0, 1, 1))
	})

	t.Run("wide cells", func(t *testing.T) {
		scr := gamma.NewScreenBuffer(10, 5)

		// Test with wide cell (e.g., emoji or CJK character)
		wideCell := &gamma.Cell{Content: "😀", Width: 2}
		scr.SetCell(0, 0, wideCell)

		t.Logf("Set wide cell at (0,0): %+v", scr.CellAt(0, 0))

		cloned := Clone(scr)
		t.Logf("Cloned cell at (0,0): %+v", cloned.CellAt(0, 0))

		if cell := cloned.CellAt(0, 0); cell == nil || cell.Content != "😀" || cell.Width != 2 {
			t.Errorf("Wide cell not cloned correctly, got %#v", cell)
		}

		// Test FillArea with wide cell
		FillArea(scr, wideCell, gamma.Rect(0, 1, 4, 1))
		for x := 0; x < 4; x += 2 {
			if cell := scr.CellAt(x, 1); cell == nil || cell.Content != "😀" || cell.Width != 2 {
				t.Errorf("Wide cell at (%d,1) not filled correctly, got %#v", x, cell)
			}
		}
	})

	t.Run("styled cells", func(t *testing.T) {
		scr := gamma.NewScreenBuffer(10, 5)

		// Test with styled cell
		styledCell := &gamma.Cell{
			Content: "S",
			Width:   1,
			Style:   gamma.Style{Attrs: gamma.AttrBold | gamma.AttrItalic},
		}
		scr.SetCell(0, 0, styledCell)

		cloned := Clone(scr)
		if cell := cloned.CellAt(0, 0); cell == nil || cell.Content != "S" || (cell.Style.Attrs&gamma.AttrBold == 0) || (cell.Style.Attrs&gamma.AttrItalic == 0) {
			t.Errorf("Styled cell not cloned correctly, got %v", cell)
		}
	})

	t.Run("cells with links", func(t *testing.T) {
		scr := gamma.NewScreenBuffer(10, 5)

		// Test with cell containing hyperlink
		linkedCell := &gamma.Cell{
			Content: "L",
			Width:   1,
			Link:    gamma.NewLink("https://example.com", "id=test"),
		}
		scr.SetCell(0, 0, linkedCell)

		cloned := Clone(scr)
		if cell := cloned.CellAt(0, 0); cell == nil || cell.Content != "L" || cell.Link.URL != "https://example.com" {
			t.Errorf("Cell with link not cloned correctly, got %v", cell)
		}
	})
}

func TestMinimalScreenFallbacks(t *testing.T) {
	t.Run("Clear fallback", func(t *testing.T) {
		scr := newMinimalMockScreen(5, 3)

		// Set some cells
		testCell := &gamma.Cell{Content: "X", Width: 1}
		scr.SetCell(0, 0, testCell)
		scr.SetCell(2, 1, testCell)
		scr.SetCell(4, 2, testCell)

		// Clear should use Fill fallback
		Clear(scr)

		// Verify all cells are cleared
		for y := 0; y < 3; y++ {
			for x := 0; x < 5; x++ {
				cell := scr.CellAt(x, y)
				if cell != nil && cell.Content != " " {
					t.Errorf("Cell at (%d,%d) not cleared, got %v", x, y, cell)
				}
			}
		}
	})

	t.Run("ClearArea fallback", func(t *testing.T) {
		scr := newMinimalMockScreen(5, 3)

		// Set cells across the screen
		testCell := &gamma.Cell{Content: "X", Width: 1}
		for y := 0; y < 3; y++ {
			for x := 0; x < 5; x++ {
				scr.SetCell(x, y, testCell)
			}
		}

		// ClearArea should use FillArea fallback
		area := gamma.Rect(1, 0, 3, 2)
		ClearArea(scr, area)

		// Verify only the area is cleared
		for y := 0; y < 3; y++ {
			for x := 0; x < 5; x++ {
				cell := scr.CellAt(x, y)
				if x >= 1 && x < 4 && y >= 0 && y < 2 {
					// Inside cleared area
					if cell != nil && cell.Content != " " {
						t.Errorf("Cell at (%d,%d) should be cleared, got %v", x, y, cell)
					}
				} else {
					// Outside cleared area
					if cell == nil || cell.Content != "X" {
						t.Errorf("Cell at (%d,%d) should not be cleared, got %v", x, y, cell)
					}
				}
			}
		}
	})

	t.Run("Fill fallback", func(t *testing.T) {
		scr := newMinimalMockScreen(5, 3)

		fillCell := &gamma.Cell{Content: "F", Width: 1}
		// Fill should use FillArea fallback
		Fill(scr, fillCell)

		// Verify all cells are filled
		for y := 0; y < 3; y++ {
			for x := 0; x < 5; x++ {
				cell := scr.CellAt(x, y)
				if cell == nil || cell.Content != "F" {
					t.Errorf("Cell at (%d,%d) not filled correctly, got %v", x, y, cell)
				}
			}
		}
	})

	t.Run("FillArea fallback loop", func(t *testing.T) {
		scr := newMinimalMockScreen(5, 3)

		fillCell := &gamma.Cell{Content: "A", Width: 1}
		area := gamma.Rect(1, 1, 2, 1)
		// FillArea should use the loop fallback
		FillArea(scr, fillCell, area)

		// Verify only the area is filled
		for y := 0; y < 3; y++ {
			for x := 0; x < 5; x++ {
				cell := scr.CellAt(x, y)
				if x >= 1 && x < 3 && y >= 1 && y < 2 {
					// Inside filled area
					if cell == nil || cell.Content != "A" {
						t.Errorf("Cell at (%d,%d) should be filled with 'A', got %v", x, y, cell)
					}
				} else {
					// Outside filled area (should be empty)
					if cell != nil && cell.Content != " " {
						t.Errorf("Cell at (%d,%d) should be empty, got %v", x, y, cell)
					}
				}
			}
		}
	})

	t.Run("Clone fallback", func(t *testing.T) {
		scr := newMinimalMockScreen(5, 3)

		// Set some cells
		scr.SetCell(0, 0, &gamma.Cell{Content: "A", Width: 1})
		scr.SetCell(2, 1, &gamma.Cell{Content: "B", Width: 1})
		scr.SetCell(4, 2, &gamma.Cell{Content: "C", Width: 1})

		// Clone should use CloneArea fallback
		cloned := Clone(scr)

		if cloned == nil {
			t.Fatal("Clone returned nil")
		}

		// Verify cloned buffer has same dimensions
		if cloned.Width() != 5 || cloned.Height() != 3 {
			t.Errorf("Cloned buffer has wrong dimensions: %dx%d", cloned.Width(), cloned.Height())
		}

		// Check specific cells
		if cell := cloned.CellAt(0, 0); cell == nil || cell.Content != "A" {
			t.Errorf("Cloned cell at (0,0) incorrect, got %v", cell)
		}
		if cell := cloned.CellAt(2, 1); cell == nil || cell.Content != "B" {
			t.Errorf("Cloned cell at (2,1) incorrect, got %v", cell)
		}
		if cell := cloned.CellAt(4, 2); cell == nil || cell.Content != "C" {
			t.Errorf("Cloned cell at (4,2) incorrect, got %v", cell)
		}
	})

	t.Run("CloneArea fallback loop", func(t *testing.T) {
		scr := newMinimalMockScreen(5, 3)

		// Set cells across the screen
		for y := 0; y < 3; y++ {
			for x := 0; x < 5; x++ {
				content := string(rune('A' + y*5 + x))
				scr.SetCell(x, y, &gamma.Cell{Content: content, Width: 1})
			}
		}

		// CloneArea should use the loop fallback
		area := gamma.Rect(1, 0, 3, 2)
		cloned := CloneArea(scr, area)

		if cloned == nil {
			t.Fatal("CloneArea returned nil")
		}

		// Verify cloned buffer has correct dimensions
		if cloned.Width() != 3 || cloned.Height() != 2 {
			t.Errorf("Cloned buffer has wrong dimensions: %dx%d, expected 3x2", cloned.Width(), cloned.Height())
		}

		// Verify content matches the cloned area
		for y := 0; y < 2; y++ {
			for x := 0; x < 3; x++ {
				expectedContent := string(rune('A' + y*5 + (x + 1)))
				cell := cloned.CellAt(x, y)
				if cell == nil || cell.Content != expectedContent {
					t.Errorf("Cloned cell at (%d,%d) incorrect, expected %s, got %v", x, y, expectedContent, cell)
				}
			}
		}
	})

	t.Run("CloneArea with nil and zero cells", func(t *testing.T) {
		scr := newMinimalMockScreen(5, 3)

		// Set some cells but leave some as zero/empty
		scr.SetCell(1, 0, &gamma.Cell{Content: "A", Width: 1})
		scr.SetCell(3, 1, &gamma.Cell{Content: "B", Width: 1})
		// Set a zero cell explicitly
		scr.SetCell(2, 0, &gamma.Cell{})
		// Leave other cells as empty

		area := gamma.Rect(0, 0, 5, 2)
		cloned := CloneArea(scr, area)

		if cloned == nil {
			t.Fatal("CloneArea returned nil")
		}

		// Check that non-zero cells are cloned
		if cell := cloned.CellAt(1, 0); cell == nil || cell.Content != "A" {
			t.Errorf("Cell at (1,0) should be 'A', got %v", cell)
		}
		if cell := cloned.CellAt(3, 1); cell == nil || cell.Content != "B" {
			t.Errorf("Cell at (3,1) should be 'B', got %v", cell)
		}

		// Check that empty cells remain empty
		if cell := cloned.CellAt(0, 0); cell != nil && !cell.IsZero() && cell.Content != " " {
			t.Errorf("Cell at (0,0) should be empty, got %v", cell)
		}
		if cell := cloned.CellAt(2, 0); cell != nil && !cell.IsZero() && cell.Content != " " {
			t.Errorf("Cell at (2,0) should be empty, got %v", cell)
		}
		if cell := cloned.CellAt(2, 1); cell != nil && !cell.IsZero() && cell.Content != " " {
			t.Errorf("Cell at (2,1) should be empty, got %v", cell)
		}
	})

	t.Run("CloneArea with nil cells", func(t *testing.T) {
		scr := newNilCellMockScreen(5, 3)

		// Set some cells
		scr.SetCell(1, 0, &gamma.Cell{Content: "A", Width: 1})
		scr.SetCell(3, 1, &gamma.Cell{Content: "B", Width: 1})

		// Make position (2, 1) return nil
		scr.SetNilAt(2, 1)

		area := gamma.Rect(0, 0, 5, 2)
		cloned := CloneArea(scr, area)

		if cloned == nil {
			t.Fatal("CloneArea returned nil")
		}

		// Check that non-nil cells are cloned
		if cell := cloned.CellAt(1, 0); cell == nil || cell.Content != "A" {
			t.Errorf("Cell at (1,0) should be 'A', got %v", cell)
		}
		if cell := cloned.CellAt(3, 1); cell == nil || cell.Content != "B" {
			t.Errorf("Cell at (3,1) should be 'B', got %v", cell)
		}

		// Check that nil cell position remains empty in cloned buffer
		if cell := cloned.CellAt(2, 1); cell != nil && !cell.IsZero() && cell.Content != " " {
			t.Errorf("Cell at (2,1) should be empty (was nil), got %v", cell)
		}
	})
}

func BenchmarkClear(b *testing.B) {
	scr := gamma.NewScreenBuffer(80, 24)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Clear(scr)
	}
}

func BenchmarkFill(b *testing.B) {
	scr := gamma.NewScreenBuffer(80, 24)
	cell := &gamma.Cell{Content: "X", Width: 1}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Fill(scr, cell)
	}
}

func BenchmarkClone(b *testing.B) {
	scr := gamma.NewScreenBuffer(80, 24)
	// Fill with some content
	for y := 0; y < 24; y++ {
		for x := 0; x < 80; x++ {
			scr.SetCell(x, y, &gamma.Cell{Content: "X", Width: 1})
		}
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Clone(scr)
	}
}

func BenchmarkCloneArea(b *testing.B) {
	scr := gamma.NewScreenBuffer(80, 24)
	// Fill with some content
	for y := 0; y < 24; y++ {
		for x := 0; x < 80; x++ {
			scr.SetCell(x, y, &gamma.Cell{Content: "X", Width: 1})
		}
	}
	area := gamma.Rect(10, 5, 20, 10)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		CloneArea(scr, area)
	}
}
