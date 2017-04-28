package main

import "testing"

func Test_getRange(t *testing.T) {
	input := "3"
	expect := fRange{3, 3}

	output, err := getRange(input)
	if err != nil {
		t.Error(err)
	}
	if expect != output {
		t.Error("error output:", output)
	}

	input = "3:0"
	expect = fRange{0, 3}
	output, err = getRange(input)
	if err != nil {
		t.Error(err)
	}
	if expect != output {
		t.Error("error output:", output)
	}

	input = "1:4"
	expect = fRange{1, 4}
	output, err = getRange(input)
	if err != nil {
		t.Error(err)
	}
	if expect != output {
		t.Error("error output:", output)
	}

	input = "test"
	output, err = getRange(input)
	if err == nil {
		t.Error("expect error")
	}

	input = "-1:0"
	output, err = getRange(input)
	if err == nil {
		t.Error("expect error")
	}

	regLen = 32
	input = "1:33"
	output, err = getRange(input)
	if err == nil {
		t.Error("expect error")
	}
}
