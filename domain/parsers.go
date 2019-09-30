package domain

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

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
