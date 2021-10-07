package id3

import (
	"fmt"
	"strings"
	"testing"
)

// From https://iq.opengenus.org/id3-algorithm/
//
const example = `outlook,temperature,humidity,wind,play
sunny,hot,high,weak,no
sunny,hot,high,strong,no
overcast,hot,high,weak,yes
rain,mild,high,weak,yes
rain,cool,normal,weak,yes
rain,cool,normal,strong,no
overcast,cool,normal,strong,yes
sunny,mild,high,weak,no
sunny,cool,normal,weak,yes
rain,mild,normal,weak,yes
sunny,mild,normal,strong,yes
overcast,mild,high,strong,yes
overcast,hot,normal,weak,yes
rain,mild,high,strong,no
`

func TestViews(t *testing.T) {
	view, err := Read(strings.NewReader(example))
	if err != nil {
		t.Error()
	}
	row := view.Next()
	if row == nil {
		t.Error()
	}
	if row[0] != "sunny" {
		t.Error()
	}
	if len(row) != 5 {
		t.Error()
	}
	//
	view = view.Select("outlook", "overcast")
	view.First()
	row = view.Next()
	row = view.Next()
	row = view.Next()
	row = view.Next()
	if row == nil {
		t.Error()
	}
	row = view.Next()
	if row != nil {
		t.Error()
	}
}

func TestEntropy(t *testing.T) {
	cases := []struct {
		probability float64
		expected    float64
	}{
		{0.5, 0.5},
		{0.0, 0.0},
		{1.0, 0.0},
	}
	for _, c := range cases {
		if Entropy(c.probability) != c.expected {
			t.Error()
		}
	}
}

func TestTotalEntropy(t *testing.T) {
	view, _ := Read(strings.NewReader(example))
	//
	// Results from https://iq.opengenus.org/id3-algorithm/
	//
	h := TotalEntropy(view, "play")
	if fmt.Sprintf("%.2f", h) != "0.94" {
		t.Error()
	}
	view = view.Select("outlook", "sunny")
	h = TotalEntropy(view, "play")
	if fmt.Sprintf("%.2f", h) != "0.97" {
		t.Error()
	}
	view = view.Drop("outlook")
	h = TotalEntropy(view, "play")
	if fmt.Sprintf("%.2f", h) != "0.97" {
		t.Error()
	}
}

func TestAverageEntropy(t *testing.T) {
	view, _ := Read(strings.NewReader(example))
	//
	// Results from https://iq.opengenus.org/id3-algorithm/
	//
	h := AverageEntropy(view, "outlook", "play")
	if fmt.Sprintf("%.2f", h) != "0.69" {
		t.Error()
	}
	//
	// Dropping an irrelevant column should have no impact.
	//
	h = AverageEntropy(view.Drop("temperature"), "outlook", "play")
	if fmt.Sprintf("%.2f", h) != "0.69" {
		t.Error()
	}
}

func TestLearningOutput(t *testing.T) {
	view, _ := Read(strings.NewReader(example))
	decision := Learn(view, "play")
	//
	//
	//
	b, err := decision.ToJSON(true)
	if err != nil {
		t.Error()
	}
	fmt.Println(string(b))
	//
	//
	//
	d, err := FromJSON(b)
	if err != nil {
		t.Error()
	}
	b2, _ := d.ToJSON(true)
	if len(b2) != len(b2) {
		t.Error()
	}
}

func TestDecide(t *testing.T) {
	// t.Skip()
	//
	// Set up.
	//
	rule := &Decision{
		Column: "outlook",
		Cases: []*Case{
			{
				Value: "sunny",
				Class: "",
				Decide: &Decision{
					Column: "humidity",
					Cases: []*Case{
						{
							Value:  "high",
							Class:  "no",
							Decide: nil,
						},
						{
							Value:  "normal",
							Class:  "yes",
							Decide: nil,
						},
					},
				},
			},
			{
				Value:  "overcast",
				Class:  "yes",
				Decide: nil,
			},
			{
				Value: "rain",
				Class: "",
				Decide: &Decision{
					Column: "wind",
					Cases: []*Case{
						{
							Value:  "weak",
							Class:  "yes",
							Decide: nil,
						},
						{
							Value:  "strong",
							Class:  "no",
							Decide: nil,
						},
					},
				},
			},
		},
	}
	data := [][]string{
		{"outlook", "temperature", "humidity", "wind", "play"},
		{"sunny", "hot", "high", "weak", "no"},
	}
	//
	// Test.
	//
	answer := rule.Decide(data)
	if answer[0] != data[1][4] {
		t.Error()
	}
}
