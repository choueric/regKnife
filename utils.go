package main

import "strconv"

// s is hex(0x), oct(0) otherwise decial
func parseInt(s string) (int64, error) {
	return strconv.ParseInt(s, 0, 64)
}

// s is binary string "11001" -> 0x19
func parseBin(s string) (int64, error) {
	return strconv.ParseInt(s, 2, 64)
}
