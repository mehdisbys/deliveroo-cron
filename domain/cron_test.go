package domain

import (
	"testing"

	"github.com/go-test/deep"
)

func TestParseExpression(t *testing.T) {
	tests := []struct {
		name       string
		expression []string
		expected   *ParsedExpression
	}{
		{
			expression: []string{"0-59/15", "0", "1,15", "*", "1-5"},
			expected: &ParsedExpression{
				Fields: []Field{
					Field{name: "minutes", Values: []int{0, 15, 30, 45}},
					Field{name: "hour", Values: []int{0}},
					Field{name: "day of month", Values: []int{1, 15}},
					Field{name: "month", Values: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}},
					Field{name: "day of week", Values: []int{1, 2, 3, 4, 5}},
				},
			},
		},
		{
			expression: []string{"0-59/15", "0", "1,15", "*", "*"},
			expected: &ParsedExpression{
				Fields: []Field{
					Field{name: "minutes", Values: []int{0, 15, 30, 45}},
					Field{name: "hour", Values: []int{0}},
					Field{name: "day of month", Values: []int{1, 15}},
					Field{name: "month", Values: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}},
					Field{name: "day of week", Values: []int{1, 2, 3, 4, 5, 6, 7}},
				},
			},
		},
		{
			expression: []string{"*/15", "*/10", "1,13,15", "*", "*"},
			expected: &ParsedExpression{
				Fields: []Field{
					Field{name: "minutes", Values: []int{0, 15, 30, 45}},
					Field{name: "hour", Values: []int{0, 10, 20}},
					Field{name: "day of month", Values: []int{1, 13, 15}},
					Field{name: "month", Values: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}},
					Field{name: "day of week", Values: []int{1, 2, 3, 4, 5, 6, 7}},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, _ := ParseExpression(test.expression)

			if diff := deep.Equal(test.expected, res); diff != nil {
				t.Error(diff)
			}
		})
	}
}

func TestParse(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []int
	}{
		{
			name:     "parse range",
			input:    "5-10",
			expected: []int{5, 6, 7, 8, 9, 10},
		},
		{
			name:     "parse fixed",
			input:    "5",
			expected: []int{5},
		},

		{
			name:     "parse list",
			input:    "5,6,7",
			expected: []int{5, 6, 7},
		},
		{
			name:     "parse list",
			input:    "10,9,3",
			expected: []int{3, 9, 10},
		},
		{
			name:     "parse frequency",
			input:    "1-7/2",
			expected: []int{1, 3, 5, 7},
		},
		{
			name:     "parse frequency",
			input:    "1-7/4",
			expected: []int{1, 5},
		},
		{
			name:     "parse frequency",
			input:    "1-7/10",
			expected: []int{1},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			res, _ := parse(test.input, 0)

			if diff := deep.Equal(test.expected, res); diff != nil {
				t.Error(diff)
			}
		})
	}
}

func TestParseRange(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []int
	}{
		{
			input:    "5-10",
			expected: []int{5, 6, 7, 8, 9, 10},
		},
		{
			input:    "2-2",
			expected: []int{2},
		},
		{
			input:    "0-0",
			expected: []int{0},
		},
		{
			input:    "-3-5",
			expected: nil,
		},
		{
			input:    "10-5",
			expected: nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			res, _ := parseRange(test.input)

			if diff := deep.Equal(test.expected, res); diff != nil {
				t.Error(diff)
			}
		})
	}
}
