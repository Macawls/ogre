package layout

import (
	"math"
	"testing"

	"github.com/macawls/ogre/v3/parse"
	"github.com/macawls/ogre/v3/style"
)

// layoutTable builds the layout tree for an HTML fragment and returns the first
// node whose display is table.
func layoutTable(t *testing.T, html string, width, height float64) (*LayoutTree, *Node) {
	t.Helper()

	root, err := parse.Parse(html)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	styles := style.Resolve(root, width, height, style.TailwindV4)
	tree := ComputeLayout(root, styles, width, height, func(_ *parse.Node, text string, cs *style.ComputedStyle, _ float64) (float64, float64) {
		return float64(len([]rune(text))) * 10, 20
	})

	var found *Node
	var walk func(n *Node)
	walk = func(n *Node) {
		if found != nil {
			return
		}
		if n.Style.Display == DisplayTable {
			found = n
			return
		}
		for _, c := range n.Children {
			walk(c)
		}
	}
	walk(tree.Root)
	if found == nil {
		t.Fatal("no display:table node in layout tree")
	}
	return tree, found
}

func tableRows(t *testing.T, table *Node) []*Node {
	t.Helper()

	var rows []*Node
	var walk func(n *Node)
	walk = func(n *Node) {
		for _, c := range n.Children {
			switch c.Style.Display {
			case DisplayTableRow:
				rows = append(rows, c)
			case DisplayTableRowGroup:
				walk(c)
			}
		}
	}
	walk(table)
	if len(rows) == 0 {
		t.Fatal("table has no rows")
	}
	return rows
}

func approx(got, want, tol float64) bool {
	return math.Abs(got-want) <= tol
}

// Cells in a row must sit side by side, not stack. Before display:table support
// every cell was a block child, so all three columns landed at x=0.
func TestTableCellsSitSideBySide(t *testing.T) {
	_, table := layoutTable(t, `<table style="width:600px;font-size:16px"><tr>`+
		`<td style="padding:10px;width:200px">А</td>`+
		`<td style="padding:10px;width:200px">Б</td>`+
		`<td style="padding:10px;width:200px">В</td>`+
		`</tr></table>`, 1280, 720)

	rows := tableRows(t, table)
	cells := rows[0].Children
	if len(cells) != 3 {
		t.Fatalf("got %d cells, want 3", len(cells))
	}

	// Three cells of width:200px + 20px padding want 660px inside a 600px table,
	// so the browser scales the columns down to 200px each.
	wantX := []float64{0, 200, 400}
	for i, cell := range cells {
		if !approx(cell.Layout.X, wantX[i], 0.5) {
			t.Errorf("cell %d x = %.2f, want %.2f", i, cell.Layout.X, wantX[i])
		}
		if !approx(cell.Layout.Width, 200, 0.5) {
			t.Errorf("cell %d width = %.2f, want 200", i, cell.Layout.Width)
		}
		if !approx(cell.Layout.Y, 0, 0.01) {
			t.Errorf("cell %d y = %.2f, want 0 (offset is carried by the row)", i, cell.Layout.Y)
		}
	}
}

// Columns with no declared width share the table width equally.
func TestTableAutoColumnsShareSlack(t *testing.T) {
	_, table := layoutTable(t, `<table style="width:600px;font-size:18px"><tr>`+
		`<td style="padding:10px">А</td><td style="padding:10px">Б</td>`+
		`</tr></table>`, 1280, 720)

	cells := tableRows(t, table)[0].Children
	if len(cells) != 2 {
		t.Fatalf("got %d cells, want 2", len(cells))
	}
	for i, cell := range cells {
		if !approx(cell.Layout.Width, 300, 0.5) {
			t.Errorf("cell %d width = %.2f, want 300", i, cell.Layout.Width)
		}
	}
	if !approx(cells[1].Layout.X, 300, 0.5) {
		t.Errorf("second cell x = %.2f, want 300", cells[1].Layout.X)
	}
}

// Rows stack vertically and every row is as tall as its tallest cell.
func TestTableRowsStackAndSizeToTallestCell(t *testing.T) {
	_, table := layoutTable(t, `<table style="width:400px;font-size:16px">`+
		`<tr><td style="padding:10px">1</td><td style="padding:30px">2</td></tr>`+
		`<tr><td style="padding:10px">3</td><td style="padding:10px">4</td></tr>`+
		`</table>`, 1280, 720)

	rows := tableRows(t, table)
	if len(rows) != 2 {
		t.Fatalf("got %d rows, want 2", len(rows))
	}

	// Row 1 is sized by the 30px-padded cell: 30 + 20 + 30 = 80.
	if !approx(rows[0].Layout.Height, 80, 0.5) {
		t.Errorf("row 0 height = %.2f, want 80", rows[0].Layout.Height)
	}
	if !approx(rows[1].Layout.Y, 80, 0.5) {
		t.Errorf("row 1 y = %.2f, want 80", rows[1].Layout.Y)
	}
	// Row 2 is sized by a 10px-padded cell: 10 + 20 + 10 = 40.
	if !approx(rows[1].Layout.Height, 40, 0.5) {
		t.Errorf("row 1 height = %.2f, want 40", rows[1].Layout.Height)
	}
}

// A table never stretches to a definite container height: its height is the sum of
// its rows. The flex pass upstream hands down a definite height, which used to be
// adopted verbatim and pushed everything below the table out of place.
func TestTableHeightIgnoresDefiniteContainerHeight(t *testing.T) {
	_, table := layoutTable(t, `<div style="height:600px;display:flex;flex-direction:column">`+
		`<table style="width:400px;font-size:16px">`+
		`<tr><td style="padding:10px">1</td></tr>`+
		`<tr><td style="padding:10px">2</td></tr>`+
		`</table></div>`, 1280, 720)

	if !approx(table.Layout.Height, 80, 0.5) {
		t.Errorf("table height = %.2f, want 80 (two 40px rows), not the container's 600",
			table.Layout.Height)
	}
}

// An explicit height on the table is still honoured.
func TestTableExplicitHeightWins(t *testing.T) {
	_, table := layoutTable(t, `<table style="width:400px;height:300px;font-size:16px">`+
		`<tr><td style="padding:10px">1</td></tr></table>`, 1280, 720)

	if !approx(table.Layout.Height, 300, 0.5) {
		t.Errorf("table height = %.2f, want 300", table.Layout.Height)
	}
}

// Rows wrapped in thead/tbody are laid out exactly like bare rows.
func TestTableRowGroupsAreTransparent(t *testing.T) {
	_, table := layoutTable(t, `<table style="width:400px;font-size:16px">`+
		`<thead><tr><td style="padding:10px">А</td></tr></thead>`+
		`<tbody><tr><td style="padding:10px">1</td></tr>`+
		`<tr><td style="padding:10px">2</td></tr></tbody></table>`, 1280, 720)

	rows := tableRows(t, table)
	if len(rows) != 3 {
		t.Fatalf("got %d rows, want 3", len(rows))
	}
	wantY := []float64{0, 40, 80}
	for i, row := range rows {
		if !approx(row.Layout.Y, wantY[i], 0.5) {
			t.Errorf("row %d y = %.2f, want %.2f", i, row.Layout.Y, wantY[i])
		}
	}
}

// Cell offsets are relative to their row. Storing the absolute y on both the row
// and the cell makes the renderer add the row offset twice, so every row after the
// first lands one row-height too low.
func TestTableCellYIsRelativeToRow(t *testing.T) {
	_, table := layoutTable(t, `<table style="width:400px;font-size:16px">`+
		`<tr><td style="padding:10px">1</td><td style="padding:10px">2</td></tr>`+
		`<tr><td style="padding:10px">3</td><td style="padding:10px">4</td></tr>`+
		`</table>`, 1280, 720)

	rows := tableRows(t, table)
	for r, row := range rows {
		for c, cell := range row.Children {
			if !approx(cell.Layout.Y, 0, 0.01) {
				t.Errorf("row %d cell %d y = %.2f, want 0 (the row carries the offset)",
					r, c, cell.Layout.Y)
			}
		}
	}
	if !approx(rows[1].Layout.Y, 40, 0.5) {
		t.Errorf("row 1 y = %.2f, want 40", rows[1].Layout.Y)
	}
}

// A wrapper whose height is auto must size itself from the table's rows. When a row
// is treated as a column by the intrinsic-size estimator, its cells' heights are
// summed instead of maxed and the wrapper ends up twice as tall as the table.
func TestTableAutoHeightWrapperFitsRows(t *testing.T) {
	tree, table := layoutTable(t, `<!DOCTYPE html><html><head><meta charset="utf-8"></head>`+
		`<body style="margin:0;padding:0;width:1280px;height:720px">`+
		`<div style="padding:56px">`+
		`<table style="width:400px;font-size:16px">`+
		`<tr><td style="padding:10px">1</td><td style="padding:10px">2</td></tr>`+
		`<tr><td style="padding:10px">3</td><td style="padding:10px">4</td></tr>`+
		`</table></div></body></html>`, 1280, 720)

	var wrapper *Node
	var find func(n *Node)
	find = func(n *Node) {
		for _, c := range n.Children {
			if c == table {
				wrapper = n
				return
			}
			find(c)
		}
	}
	find(tree.Root)
	if wrapper == nil {
		t.Fatal("table has no parent node")
	}
	if !approx(wrapper.Layout.Height, table.Layout.Height+112, 0.5) {
		t.Errorf("wrapper height = %.2f, want %.2f (table %.2f + 2x56 padding)",
			wrapper.Layout.Height, table.Layout.Height+112, table.Layout.Height)
	}
}

// A cell with display:none is skipped and does not consume a column.
func TestTableSkipsHiddenCells(t *testing.T) {
	_, table := layoutTable(t, `<table style="width:400px;font-size:16px"><tr>`+
		`<td style="padding:10px">1</td>`+
		`<td style="padding:10px;display:none">2</td>`+
		`<td style="padding:10px">3</td>`+
		`</tr></table>`, 1280, 720)

	cells := tableRows(t, table)[0].Children
	if len(cells) != 2 {
		t.Fatalf("got %d laid-out cells, want 2", len(cells))
	}
	for i, cell := range cells {
		if !approx(cell.Layout.Width, 200, 0.5) {
			t.Errorf("cell %d width = %.2f, want 200", i, cell.Layout.Width)
		}
	}
}
