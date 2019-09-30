package domain

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// Cron can take different input types
// 5 : a fixed value
// 5-10 : a range
// */5 : a frequency
// 5,7 : a list of values

type ElementType int

const (
	InvalidType ElementType = iota
	Fixed
	Range
	Frequency
	List
)

var format = []Field{
	{
		name:     "minutes",
	},
	{
		name:     "hour",
	},
	{
		name:     "day of month",
	},
	{
		name:     "month",
	},
	{
		name:     "day of week",
	},
}

var ranges = [][2]int{

	{
		0, 59, // minutes
	},
	{
		0, 23, // hours
	},
	{
		1, 31, // day of month
	},
	{
		1, 12, // month
	},
	{
		1, 7, // day of week
	},
}

var CronElementsSize = len(format)

const (
	FrequencyDelimiter = "/"
	ListDelimiter      = ","
	RangeDelimiter     = "-"
	FixedDelimiter     = ""
	Wildcard           = "*"
)

type Field struct {
	position int
	name     string
	Values   []int
}

type ParsedExpression struct {
	Fields []Field
}

type CronParser interface {
	ParseExpression() (*ParsedExpression, error)
}

func newParsedExpression(elements [][]int) *ParsedExpression {
	fields := make([]Field, len(format))
	copy(fields, format)

	for i,e := range elements {
		fields[i].Values = e
	}

	return &ParsedExpression{Fields: fields}
}

func ParseExpression(elements []string) (*ParsedExpression, error) {

	if len(elements) != CronElementsSize {
		return nil, errors.New(
			fmt.Sprintf("malformed expression expecting %d elements but got %d",
				CronElementsSize,
				len(elements)))
	}

	res := make([][]int, len(elements))
	var err error

	for i, e := range elements {
		// parse each element and save result
		res[i], err = parse(e, i)
		if err != nil {
			return nil, err
		}
	}
	return newParsedExpression(res), nil
}

func parse(input string, position int) ([]int, error) {
	// replace wildcard with values we can use
	input = replaceWildcard(input, ranges[position])

	// get type of input for differentiated parsing
	inputType, err := getType(input)
	if err != nil {
		return nil, err
	}

	// each input type has a different parsing method
	switch inputType {
	case Fixed:
		return parseFixed(input)
	case Range:
		return parseRange(input)
	case Frequency:
		return parseFrequency(input)
	case List:
		return parseList(input)
	default:
		return nil, errors.New("unrecognised type")
	}
}

// replaceWildcard replaces the wildcard by its equivalent range
func replaceWildcard(input string, replaceWith [2]int) string {
	r := fmt.Sprintf("%d-%d", replaceWith[0], replaceWith[1])
	res := strings.ReplaceAll(input, Wildcard, r)
	return res
}

// getType returns the ElementType of input which will allow
// to determine which function to call
func getType(input string) (ElementType, error) {
	var element string
	nonAlpha := getNonAlphaNumerical(input)
	lna := len(nonAlpha)

	// if more than 1 non-alphanumerical character and contain '/'
	// then it is type Frequency
	if lna > 1 {
		// the slash will always be at index 1
		if nonAlpha[1] == FrequencyDelimiter {
			return Frequency, nil
		}
	}

	if lna >= 1 {
		if nonAlpha[0] == ListDelimiter {
			return List, nil
		}
		element = nonAlpha[0]
	}

	switch element {
	case RangeDelimiter:
		return Range, nil
	case ListDelimiter:
		return List, nil
	case FixedDelimiter:
		return Fixed, nil
	default:
		return InvalidType, errors.New("could not determine type")
	}
}

// getNonAlphaNumerical returns all non-alphanumerical chars in input
func getNonAlphaNumerical(input string) []string {
	nonAlphaNumerical := []string{}
	re := regexp.MustCompile(`[^a-zA-Z0-9]+`)

	submatchall := re.FindAllString(input, -1)
	for _, element := range submatchall {
		nonAlphaNumerical = append(nonAlphaNumerical, element)
	}
	return nonAlphaNumerical
}
