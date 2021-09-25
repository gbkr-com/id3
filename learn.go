package id3

import (
	"math"
	"sort"
)

// Distinct is a distinct column value and its associated probability.
//
type Distinct struct {
	Value       string
	Probability float64
}

// Likelihood returns the probability of each distinct value in the named column
// of the view. The slice is sorted in decreasing probability.
//
func Likelihood(view View, column string) []Distinct {
	//
	// Find the distinct values and count the frequency.
	//
	i := find(view.Columns(), column)
	distinct := make(map[string]float64)
	total := 0.0
	view.First()
	for {
		row := view.Next()
		if row == nil {
			break
		}
		v := row[i]
		if _, ok := distinct[v]; !ok {
			distinct[v] = 0
		}
		distinct[v]++
		total++
	}
	//
	// Convert the map to a slice, then sort.
	//
	var sorted []Distinct
	for k, v := range distinct {
		sorted = append(sorted, Distinct{Value: k, Probability: v / total})
	}
	sort.Slice(
		sorted,
		func(i, j int) bool {
			return sorted[i].Probability > sorted[j].Probability
		},
	)
	return sorted
}

// Entropy returns the Shannon entropy for the given probability. It converts
// the edge cases of probability zero and one to a zero entropy value.
//
func Entropy(p float64) float64 {
	switch {
	case p == 0 || p == 1:
		return 0
	default:
		return -p * math.Log2(p)
	}
}

// TotalEntropy returns the total entropy of the class column in the view.
//
func TotalEntropy(view View, class string) (h float64) {
	for _, v := range Likelihood(view, class) {
		h += Entropy(v.Probability)
	}
	return
}

// AverageEntropy returns the average entropy of the class column over each
// distinct value of the attribute column.
//
func AverageEntropy(view View, attribute, class string) (h float64) {
	//
	// Use find to confirm the class column exists.
	//
	find(view.Columns(), class)
	//
	// Calculate the probability weighted class entropy for each of the
	// distinct values.
	//
	for _, v := range Likelihood(view, attribute) {
		h += v.Probability * TotalEntropy(view.Select(attribute, v.Value), class)
	}
	return
}

// Learn runs the ID3 algorithm on the given view using the named class column.
//
func Learn(view View, class string) *Decision {
	//
	// Calculate the total entropy of this view and the information gain from
	// each column (ignoring the class column).
	//
	h := TotalEntropy(view, class)
	cols := view.Columns()
	gain := make([]float64, len(cols))
	maxGain := -1.0
	maxColumn := ""
	for i, v := range cols {
		if v == class || v == "" {
			continue
		}
		gain[i] = h - AverageEntropy(view, v, class)
		if gain[i] > maxGain {
			maxGain = gain[i]
			maxColumn = cols[i]
		}
	}
	//
	// The column with the maximum gain is the basis for the decision.
	//
	decision := &Decision{Column: maxColumn}
	//
	// For each distinct value in the maximum gain column, in decreasing
	// probability, check if the value is terminal or whether to recurse.
	//
	for _, v := range Likelihood(view, maxColumn) {
		c := &Case{Value: v.Value}
		decision.Cases = append(decision.Cases, c)
		//
		// The case is terminal if there is a single class for all rows, in
		// which case the total entropy would be zero.
		//
		subview := view.Select(maxColumn, v.Value)
		subh := TotalEntropy(subview, class)
		if subh == 0.0 {
			//
			// Get the first class value in this subview.
			//
			subview.First()
			row := subview.Next()
			c.Class = row[find(subview.Columns(), class)]
		} else {
			//
			// Recurse on this view dropping the just decided column.
			//
			c.Decide = Learn(subview.Drop(maxColumn), class)
		}
	}
	return decision
}
