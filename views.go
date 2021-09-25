package id3

import (
	"encoding/csv"
	"io"
)

// View is the interface for ID3 to inspect CSV conformant data. It provides
// a cursor like mechanism for reading the data, through the First() and Next()
// functions.
//
type View interface {

	// Columns returns the column names in this view. Columns with names of ""
	// have been hidden from this view - see Drop method.
	//
	Columns() []string

	// First returns to before the first row in this view.
	//
	First()

	// Next returns the next row in the view, or nil if there are no more
	// rows.
	//
	Next() []string

	// Select returns a view that shows only rows having the given value in the
	// column.
	//
	Select(column, value string) View

	// Drop returns a view which 'hides' the named column.
	//
	Drop(column string) View
}

// Read CSV conformant data from the given reader and return a View on that.
//
func Read(reader io.Reader) (View, error) {
	r := csv.NewReader(reader)
	data, err := r.ReadAll()
	if err != nil {
		return nil, err
	}
	return &baseView{
		data: data,
		next: 1,
	}, nil
}

func find(slice []string, x string) int {
	for i, str := range slice {
		if str == x {
			return i
		}
	}
	panic("id3: '" + x + "'not in slice")
}

////////////////////////////////////////////////////////////////////////////////

type baseView struct {
	data [][]string // The original CSV conformant data.
	next int        // The index of the next row to return.
}

func (b *baseView) Columns() []string { return b.data[0] }

func (b *baseView) First() { b.next = 1 }

func (b *baseView) Next() []string {
	//
	// Skip the header row if that is next up.
	//
	if b.next == 0 {
		b.next = 1
	}
	if b.next == len(b.data) {
		return nil
	}
	row := b.data[b.next]
	b.next++
	return row
}

func (b *baseView) Select(column, value string) View {
	return &selectView{
		parent: b,
		col:    find(b.Columns(), column),
		val:    value,
	}
}

func (b *baseView) Drop(column string) View {
	return &dropView{
		parent: b,
		drop:   find(b.Columns(), column),
	}
}

////////////////////////////////////////////////////////////////////////////////

type selectView struct {
	parent View   // Inherit from the parent view.
	col    int    // Column index of the column to selct on.
	val    string // Value to select in that column.
}

func (s *selectView) Columns() []string { return s.parent.Columns() }

func (s *selectView) First() { s.parent.First() }

func (s *selectView) Next() []string {
	for {
		row := s.parent.Next()
		if row == nil {
			return nil
		}
		if row[s.col] == s.val {
			return row
		}
	}
}

func (s *selectView) Select(column, value string) View {
	return &selectView{
		parent: s,
		col:    find(s.Columns(), column),
		val:    value,
	}
}

func (s *selectView) Drop(column string) View {
	return &dropView{
		parent: s,
		drop:   find(s.Columns(), column),
	}
}

////////////////////////////////////////////////////////////////////////////////

type dropView struct {
	parent View // Inherit from the parent.
	drop   int  // Column index of the column to drop from the view.
}

func (d *dropView) Columns() []string {
	//
	// Approach is to replace the column name to be dropped with "", so that
	// subsequent use of 'locate' could panic. This avoids larger scale data
	// copyng.
	//
	c := make([]string, len(d.parent.Columns()))
	copy(c, d.parent.Columns())
	c[d.drop] = ""
	return c
}

func (d *dropView) First() { d.parent.First() }

func (d *dropView) Next() []string {
	row := d.parent.Next()
	if row == nil {
		return nil
	}
	return row
}

func (d *dropView) Select(column, value string) View {
	return &selectView{
		parent: d,
		col:    find(d.Columns(), column),
		val:    value,
	}
}

func (d *dropView) Drop(column string) View {
	return &dropView{
		parent: d,
		drop:   find(d.Columns(), column),
	}
}
