package id3

import (
	"encoding/json"
)

// Decision represents a decision within the decision tree for a single column.
// Each distinct value in that column is a case. The cases are in decreasing
// probability sequence.
//
type Decision struct {
	Column string  // The name of the data column.
	Cases  []*Case // The cases for that column.
}

// A Case is a distinct value and its associated action; either a decided class
// value or a subsequent decision.
//
type Case struct {
	Value  string    // The distinct column value.
	Class  string    // The decided class value, or "" if further decision(s) are needed.
	Decide *Decision // The subsequent decision, or nil.
}

// ToJSON returns this decision as a JSON formatted bytes slice.
//
func (d *Decision) ToJSON(indent bool) ([]byte, error) {
	switch indent {
	case false:
		return json.Marshal(d)
	default:
		return json.MarshalIndent(d, "", "    ")
	}
}

// FromJSON translates the given JSON formatted byte slice into a decision.
//
func FromJSON(b []byte) (*Decision, error) {
	d := new(Decision)
	err := json.Unmarshal(b, d)
	if err != nil {
		return nil, err
	}
	return d, nil
}

// Decide on the given CSV conformant data. The first row must be the column
// headings.
//
func (d *Decision) Decide(data [][]string) (result []string) {
	for i := range data {
		if i == 0 {
			continue
		}
		result = append(result, d.decide(data, i))
	}
	return
}

func (d *Decision) decide(data [][]string, at int) string {
	i := find(data[0], d.Column)
	value := data[at][i]
	for _, c := range d.Cases {
		if value == c.Value {
			if c.Class != "" {
				return c.Class
			}
			return c.Decide.decide(data, at)
		}
	}
	panic("id3: no rule for column " + d.Column)
}
