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

// index range
type iRange struct {
	start int
	end   int
}

// 1 -> (1,1)
// 2:3 -> (2,3)
// 3:2 -> (2,3)
func getRange(input string) (iRange, error) {
	var r iRange
	if strings.Contains(input, ":") {
		index := strings.Split(input, ":")
		if len(index) == 0 {
			return r, errors.New("invalid range pattern")
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
		return r, errors.New(fmt.Sprintf("range is invalid, [%d, %d].", 0, regLen-1))
	}
	return r, nil
}

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
	if binStr == "" {
		return
	}
	dec, err := parseBin(bin)
	if err != nil {
		ui.Error(fmt.Sprintf("convert subbin to decimal failed: %v", err))
		return
	}

	fmt.Println("bin:", formateBinStr(bin))
	fmt.Println("dec:", dec)
	fmt.Printf("hex: 0x%x\n", dec)
}

func updateBit(input string, set bool) {
	if binStr == "" {
		ui.Info("empty value. Use 'value' to update.")
		return
	}

	r, err := getRange(input)
	if err != nil {
		ui.Error(fmt.Sprintf("parse range start index failed, %v", err))
		return
	}

	binByte := []byte(binStr)
	c := byte('0')
	if set {
		c = '1'
	}
	for i := r.start; i <= r.end; i++ {
		binByte[regLen-1-i] = c
	}

	// update global variable
	binStr = string(binByte)
}

func showReg(input string) {
	if binStr == "" {
		ui.Info("empty value. Use 'value' to update.")
		return
	}

	r, err := getRange(input)
	if err != nil {
		ui.Error(fmt.Sprintf("parse range start index failed, %v", err))
		return
	}

	start_index := regLen - 1 - r.end
	end_index := regLen - 1 - r.start
	subbin := binStr[start_index : end_index+1]

	printAllFormat(subbin)
}

func updateValue(s string) {
	val, err := parseInt(s)
	if err != nil {
		ui.Error(fmt.Sprintf("convert to Int failed: %v", err))
		return
	}

	s = strconv.FormatInt(val, 2)
	l := len(s)
	bin := strings.Repeat("0", regLen-l)
	binStr = bin + s
}

func writeFiled(rStr, vStr string) {
	if binStr == "" {
		return
	}

	r, err := getRange(rStr)
	if err != nil {
		ui.Error(fmt.Sprintf("parse range start index failed, %v", err))
		return
	}

	val, err := parseInt(vStr)
	if err != nil {
		ui.Error(fmt.Sprintf("convert to Int failed: %v", err))
		return
	}

	max := (2 << uint(r.end-r.start)) - 1
	if val < 0 || int(val) > max {
		ui.Error(fmt.Sprint("val is out of range [%d, %d]", 0, max))
		return
	}

	s := strconv.FormatInt(val, 2)
	l := len(s)
	sub := strings.Repeat("0", r.end-r.start+1-l)
	sub = sub + s
	subByte := []byte(sub)
	fmt.Println(sub)

	binByte := []byte(binStr)
	j := r.end - r.start
	for i := r.start; i <= r.end; i++ {
		binByte[regLen-1-i] = subByte[j]
		j--
	}

	// update global variable
	binStr = string(binByte)
}

func handleInput(input string) (exit bool) {
	exit = false
	cmdline := strings.Fields(input)

	if len(cmdline) == 0 {
		printUsage()
		return
	}

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
			ui.Error("Needs argument: <range>")
			return
		}
		updateValue(cmdline[1])
		printAllFormat(binStr)
	case "set", "s", "clear", "c":
		if len(cmdline) < 2 {
			ui.Error("Needs argument: <range>")
			return
		}
		set := true
		if cmdline[0] == "clear" || cmdline[0] == "c" {
			set = false
		}
		updateBit(cmdline[1], set)
		printAllFormat(binStr)
	case "write", "w":
		if len(cmdline) < 3 {
			ui.Error("Needs arguments: <range> <val>")
		}
		writeFiled(cmdline[1], cmdline[2])
		printAllFormat(binStr)
	default:
		showReg(cmdline[0])
	}

	return
}

func printUsage() {
	ui.Output("Usage:")
	ui.Output("  [h]elp          : print this message.")
	ui.Output("  [p]rint         : print input value.")
	ui.Output("  [v]alue <val>   : input value.")
	ui.Output("  [s]et <bit>     : set <bit> to 1.")
	ui.Output("  [c]lear <bit>   : clear <bit> to 0.")
	ui.Output("  [w]rite <r> <v> : write val <v> into field range <r>.")
	ui.Output("  <range>         : read the value of field range <range>, like 1 or 2:3.")
	ui.Output("  exit            : exit this program.")
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
		input, err := ui.Ask("\n>>>")
		if err != nil {
			fmt.Println(err)
		}

		if handleInput(input) {
			break
		}
	}
}
