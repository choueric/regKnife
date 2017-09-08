package main

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// parseInt convert integer string into int64 type.
// format of s can be hex(0x), oct(0) otherwise decial
func parseInt(s string) (int64, error) {
	return strconv.ParseInt(s, 0, 64)
}

// parseBin convert binary string into int64 type.
// format of s is binary string. e.g. "11001" -> 0x19
func parseBin(s string) (int64, error) {
	return strconv.ParseInt(s, 2, 64)
}

// getRange convert range string into fRange structure.
// e.g. 1 -> (1,1), 2:3 -> (2,3), 3:2 -> (2,3)
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

// decorateBinStr convert binary string into decorated format
// e.g. "101010" -> "10,1010"
func decorateBinStr(bin string) string {
	var s string
	strlen := len(bin)
	count := 0

	for i := strlen - 1; i >= 0; i-- {
		if count%4 == 0 && count != 0 {
			s = "," + s
		}
		count++
		s = string(bin[i]) + s
	}
	return s
}

// outputTriFormat outputs binary string into three format: decorated binary,
// decimal and heximal.
func outputTriFormat(w io.Writer, bin string) {
	if binStr == "" {
		return
	}
	dec, err := parseBin(bin)
	if err != nil {
		inst.Error(fmt.Sprintf("parseBin failed: %v\n", err))
		return
	}

	fmt.Fprintln(w, "bin:", decorateBinStr(bin))
	fmt.Fprintln(w, "dec:", dec)
	fmt.Fprintf(w, "hex: 0x%x\n", dec)
}

// getFieldStr gets the field string from str.
func getFieldStr(rStr, str string) (string, error) {
	r, err := getRange(rStr)
	if err != nil {
		return "", err
	}

	l := len(str)
	return str[l-1-r.end : l-r.start], nil
}

func listOffsets(w io.Writer, bin string, target rune) {
	idx := 0
	for i, v := range bin {
		if v == target {
			fmt.Fprintf(w, "%d", regLen-1-i)
			idx = i
			break
		}
	}

	for i := idx + 1; i < regLen; i++ {
		if bin[i] == byte(target) {
			fmt.Fprintf(w, ",%d", regLen-1-i)
		}
	}
	fmt.Fprintf(w, "\n")
}
