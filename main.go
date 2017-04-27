package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/mitchellh/cli"
)

var (
	ui     *cli.ColoredUi
	regLen = 32
	binStr string
	value  int64
)

func initUi() error {
	ui = new(cli.ColoredUi)
	if ui == nil {
		fmt.Printf("error of ui\n")
		return errors.New("failed to new cli")
	}

	bui := new(cli.BasicUi)
	bui.Reader = os.Stdin
	bui.Writer = os.Stdout
	bui.ErrorWriter = os.Stderr

	ui.Ui = bui
	ui.OutputColor = cli.UiColorNone
	ui.InfoColor = cli.UiColorGreen
	ui.ErrorColor = cli.UiColorRed
	ui.WarnColor = cli.UiColorYellow

	return nil
}

func updateBit(bit int, set bool) {
	bit = regLen - 1 - bit
	binByte := []byte(binStr)

	if set {
		binByte[bit] = '1'
	} else {
		binByte[bit] = '0'
	}

	binStr = string(binByte)
}

func updateBinStr(val int64) string {
	s := strconv.FormatInt(val, 2)
	l := len(s)
	bin := strings.Repeat("0", regLen-l)
	bin = bin + s
	return bin
}

func formateBinStr(bin string) string {
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

func printAllFormat(bin string) {
	dec, err := strconv.ParseInt(bin, 2, 64)
	if err != nil {
		ui.Error(fmt.Sprintf("convert subbin to decimal failed: %v", err))
		return
	}

	fmt.Println("bin:", formateBinStr(bin))
	fmt.Println("dec:", dec)
	fmt.Printf("hex: 0x%x\n", dec)
}

func showBits(start, end int) {
	start_index := regLen - 1 - end
	end_index := regLen - 1 - start
	subbin := binStr[start_index : end_index+1]

	printAllFormat(subbin)
}

func handleViewReg(input string) {
	var start, end int
	if strings.Contains(input, ":") {
		index := strings.Split(input, ":")
		if len(index) == 0 {
			ui.Error("invalid range pattern")
			return
		}

		v, err := strconv.ParseInt(index[0], 0, 64)
		if err != nil {
			ui.Error(fmt.Sprintf("parse range start index failed, %v", err))
			return
		}
		start = int(v)
		v, err = strconv.ParseInt(index[1], 0, 64)
		if err != nil {
			ui.Error(fmt.Sprintf("parse range end index failed, %v", err))
			return
		}
		end = int(v)
		if start > end {
			start, end = end, start
		}
	} else {
		v, err := strconv.ParseInt(input, 0, 64)
		if err != nil {
			ui.Error(fmt.Sprintf("parse single index failed, %v", err))
			return
		}
		start = int(v)
		end = int(v)
	}

	if start < 0 || end >= regLen {
		ui.Error(fmt.Sprintf("range is invalid, [%d, %d].", 0, regLen-1))
		return
	}
	showBits(start, end)
}

func changeValue(s string) string {
	val, err := strconv.ParseInt(s, 0, 64)
	if err != nil {
		ui.Error(fmt.Sprintf("convert to Int failed: %v", err))
		return ""
	}
	return updateBinStr(val)
}

func handleInput(input string) (exit bool) {
	exit = false
	cmdline := strings.Fields(input)

	switch cmdline[0] {
	case "exit":
		exit = true
		return
	case "help", "h":
		printUsage()
	case "print", "p":
		printAllFormat(binStr)
	case "value", "v":
		if len(cmdline) < 2 {
			ui.Error("Needs an argument")
			return
		}
		binStr = changeValue(cmdline[1])
		printAllFormat(binStr)
	case "set", "s", "clear", "c":
		if len(cmdline) < 2 {
			ui.Error("Needs an argument")
			return
		}
		bit, err := strconv.Atoi(cmdline[1])
		if err != nil {
			ui.Error(fmt.Sprintf("%v", err))
			return
		}
		set := true
		if cmdline[0] == "clear" || cmdline[0] == "c" {
			set = false
		}
		updateBit(bit, set)
		printAllFormat(binStr)
	default:
		handleViewReg(input)
	}

	return
}

func printUsage() {
	ui.Output("Usage:")
	ui.Output("  [h]elp         : print this message.")
	ui.Output("  [p]rint        : print input value.")
	ui.Output("  [v]alue <val>  : input value.")
	ui.Output("  [s]et <bit>    : set <bit> to 1.")
	ui.Output("  [c]lear <bit>  : clear <bit> to 0.")
	ui.Output("  <patter>       : bit number or range of bits, like 21, 12:14 or 14:12.")
	ui.Output("  exit           : exit this program.")
}

func main() {
	flag.IntVar(&regLen, "l", 32, "register length.")
	flag.Parse()

	if len(os.Args) == 2 && (os.Args[1] == "help" || os.Args[1] == "-h") {
		flag.Usage()
		os.Exit(0)
	}

	if err := initUi(); err != nil {
		fmt.Println(err)
		os.Exit(3)
	}

	for {
		fmt.Println()
		input, err := ui.Ask(">>>")
		if err != nil {
			fmt.Println(err)
		}

		if handleInput(input) {
			break
		}
	}
}
