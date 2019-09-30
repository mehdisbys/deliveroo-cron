package domain

import (
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strconv"
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

type ParsedExpression struct {
	Minutes    []int
	Hour       []int
	DayOfMonth []int
	Month      []int
	DayOfWeek  []int
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

type CronParser interface {
	ParseExpression() (*ParsedExpression, error)
}

func NewParsedExpression(elements [][]int) *ParsedExpression {
	return &ParsedExpression{
		Minutes:    elements[0],
		Hour:       elements[1],
		DayOfMonth: elements[2],
		Month:      elements[3],
		DayOfWeek:  elements[4],
	}
}

func ParseExpression(elements []string) (*ParsedExpression, error) {

	if len(elements) != 5 {
		return nil, errors.New(fmt.Sprintf("malformed expression expecting 5 elements but got %d", len(elements)))
	}

	res := make([][]int, len(elements))
	var err error

	for i, e := range elements {
		res[i], err = Parse(e, i)
		if err != nil {
			return nil, err
		}
	}

	return NewParsedExpression(res), nil
}

func Parse(input string, position int) ([]int, error) {
	input = replaceWildcard(input, ranges[position])

	inputType, err := getType(input)
	if err != nil {
		return nil, err
	}

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

func replaceWildcard(input string, replaceWith [2]int) string {
	r := fmt.Sprintf("%d-%d", replaceWith[0], replaceWith[1])
	res := strings.ReplaceAll(input, "*", r)
	return res
}

func parseFixed(input string) ([]int, error) {
	value, err := strconv.Atoi(input)
	return []int{value}, err
}

func parseRange(input string) ([]int, error) {
	rangeVals := strings.Split(input, "-")

	if len(rangeVals) != 2 {
		return nil, errors.New(fmt.Sprintf("was expecting two elements in range but got %d", len(rangeVals)))
	}

	left, err := strconv.Atoi(rangeVals[0])
	if err != nil {
		return nil, err
	}

	right, err := strconv.Atoi(rangeVals[1])
	if err != nil {
		return nil, err
	}

	if left < 0 {
		return nil, errors.New("left value in range is negative")
	}

	if right < 0 {
		return nil, errors.New("right value in range is negative")
	}

	if left > right {
		return nil, errors.New("left value should be inferior or equal to right value")
	}

	res := []int{}
	for i := left; i <= right; i++ {
		res = append(res, i)
	}
	return res, nil
}

func parseList(input string) ([]int, error) {
	rangeVals := strings.Split(input, ",")

	if len(rangeVals) < 2 {
		return nil, errors.New(fmt.Sprintf("was expecting at least two elements in list but got %d", len(rangeVals)))
	}

	res := []int{}

	// convert each element in string to an int
	for _, r := range rangeVals {
		val, err := strconv.Atoi(r)
		if err != nil {
			return nil, err
		}

		if val < 0 {
			return nil, errors.New("value in list is negative")
		}

		res = append(res, val)
	}

	sort.Ints(res)
	return res, nil
}

func parseFrequency(input string) ([]int, error) {
	rangeVals := strings.Split(input, "/")

	if len(rangeVals) != 2 {
		return nil, errors.New(fmt.Sprintf("was expecting two elements in list but got %d", len(rangeVals)))
	}

	left := strings.Split(rangeVals[0], "-")

	leftRange := []int{}
	var err error
	var val int

	switch len(left) {
	case 1: // fixed value
		val, err = strconv.Atoi(left[0])
		leftRange = []int{val}
	case 2: // range
		leftRange, err = parseRange(rangeVals[0])
	default: // error
		return nil, errors.New("left side of expression is malformed")
	}

	if err != nil {
		return nil, err
	}

	right, err := strconv.Atoi(rangeVals[1])
	if err != nil {
		return nil, err
	}

	res := []int{leftRange[0]}
	for i := leftRange[0] + right; i <= leftRange[len(leftRange)-1]; i += right {
		res = append(res, i)
	}
	return res, nil
}

func getType(input string) (ElementType, error) {
	var element string
	nonAlpha := getNonAlphaNumerical(input)
	lna := len(nonAlpha)

	// if more than 1 non-alphanumerical character and contain slash
	// then it is type Frequency
	if lna > 1 {
		// the slash will always be at index 1
		if nonAlpha[1] == "/" {
			return Frequency, nil
		}
	}

	if lna >= 1 {
		if nonAlpha[0] == "," {
			return List, nil
		}
		element = nonAlpha[0]
	}

	switch element {
	case "-":
		return Range, nil
	case ",":
		return List, nil
	case "":
		return Fixed, nil
	default:
		return InvalidType, errors.New("could not determine type")
	}
}

func getNonAlphaNumerical(input string) []string {
	nonAlphaNumerical := []string{}
	re := regexp.MustCompile(`[^a-zA-Z0-9]+`)

	submatchall := re.FindAllString(input, -1)
	for _, element := range submatchall {
		nonAlphaNumerical = append(nonAlphaNumerical, element)
	}
	return nonAlphaNumerical
}
