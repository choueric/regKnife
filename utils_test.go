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

func Test_decorateBinStr(t *testing.T) {
	input := "10"
	expect := "10"
	output := decorateBinStr(input)
	if output != expect {
		t.Error("error output:", output)
	}

	input = "101010"
	expect = "10,1010"
	output = decorateBinStr(input)
	if output != expect {
		t.Error("error output:", output)
	}

	input = "110101010"
	expect = "1,1010,1010"
	output = decorateBinStr(input)
	if output != expect {
		t.Error("error output:", output)
	}

	input = ""
	expect = ""
	output = decorateBinStr(input)
	if output != expect {
		t.Error("error output:", output)
	}
}

func Test_getFieldStr(t *testing.T) {
	str := "11110000"
	input := "1:3"
	expect := "000"
	output, err := getFieldStr(input, str)
	if err != nil {
		t.Error(err)
	}
	if output != expect {
		t.Error("error output:", output)
	}

	input = "5:3"
	expect = "110"
	output, err = getFieldStr(input, str)
	if err != nil {
		t.Error(err)
	}
	if output != expect {
		t.Error("error output:", output)
	}
}
