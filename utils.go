package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// s is hex(0x), oct(0) otherwise decial
func parseInt(s string) (int64, error) {
	return strconv.ParseInt(s, 0, 64)
}

// s is binary string "11001" -> 0x19
func parseBin(s string) (int64, error) {
	return strconv.ParseInt(s, 2, 64)
}

// 1 -> (1,1)
// 2:3 -> (2,3)
// 3:2 -> (2,3)
func getRange(input string) (fRange, error) {
	var r fRange
	if input == "" {
		return r, errors.New("getRange empty input")
	}
	if strings.Contains(input, ":") {
		index := strings.Split(input, ":")
		if len(index) == 0 {
			return r, errors.New("getRange invalid pattern")
		}

		v, err := parseInt(index[0])
		if err != nil {
			return r, err
		}
		r.start = int(v)
		v, err = parseInt(index[1])
		if err != nil {
			return r, err
		}
		r.end = int(v)
		if r.start > r.end {
			r.start, r.end = r.end, r.start
		}
	} else {
		v, err := parseInt(input)
		if err != nil {
			return r, err
		}
		r.start = int(v)
		r.end = int(v)
	}

	if r.start < 0 || r.end >= regLen {
		return r, errors.New(fmt.Sprintf("getRange invalid range [%d, %d].", 0, regLen-1))
	}
	return r, nil
}
