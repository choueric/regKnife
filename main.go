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
)

// field range
type fRange struct {
	start int
	end   int
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

// setFieldRangeOfBinStr sets to 1 or clears to 0 with the filed range of binStr
// rStr is a filed range string
func setFieldOfBinStr(rStr string, set bool) {
	r, err := getRange(rStr)
	if err != nil {
		ui.Error(fmt.Sprintf("parse field start index failed, %v", err))
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

	binStr = string(binByte)
}

// updateBinStr converts decimal or heximal format s into register-length
// binary string and update the global variable binStr
// e.g. "0x11" -> "00010001" (regLen = 8)
func updateBinStr(s string) {
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

// writeFiledOfBinStr changes the value in filed rStr by vStr in binStr
func writeFiledOfBinStr(rStr, vStr string) {
	r, err := getRange(rStr)
	if err != nil {
		ui.Error(fmt.Sprintf("parse field start index failed, %v", err))
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

	binStr = string(binByte)
}

func showTruesOffset(target rune) {
	listOffsets(os.Stdout, binStr, target)
}

// showFieldOfBinStr shows the filed rStr of binStr in three formats.
func showFieldOfBinStr(rStr string) {
	subbin, err := getFieldStr(rStr, binStr)
	if err != nil {
		ui.Error(fmt.Sprintf("getFieldStr failed: %v", err))
		return
	}

	outputTriFormat(os.Stdout, subbin)
}

func handleInput(input string) (exit bool) {
	exit = false

	cmdline := strings.Fields(input)
	if len(cmdline) == 0 {
		printUsage()
		return
	}

	switch cmdline[0] {
	case "exit", "quit":
		exit = true
	case "help", "h":
		printUsage()
	case "print", "p":
		outputTriFormat(os.Stdout, binStr)
	case "value", "v":
		if len(cmdline) < 2 {
			ui.Error("Needs argument: <field>")
			return
		}
		updateBinStr(cmdline[1])
		outputTriFormat(os.Stdout, binStr)
	case "set", "s", "clear", "c":
		if len(cmdline) < 2 {
			ui.Error("Needs argument: <field>")
			return
		}
		set := true
		if cmdline[0] == "clear" || cmdline[0] == "c" {
			set = false
		}
		setFieldOfBinStr(cmdline[1], set)
		outputTriFormat(os.Stdout, binStr)
	case "write", "w":
		if len(cmdline) < 3 {
			ui.Error("Needs arguments: <field> <val>")
			return
		}
		writeFiledOfBinStr(cmdline[1], cmdline[2])
		outputTriFormat(os.Stdout, binStr)
	case "list", "l":
		target := '1'
		if len(cmdline) == 2 && cmdline[1] == "0" {
			target = '0'
		}
		showTruesOffset(target)
	default:
		showFieldOfBinStr(cmdline[0])
	}

	return
}

func printUsage() {
	ui.Output("Usage:")
	ui.Output("  [h]elp          : print this message.")
	ui.Output("  [p]rint         : show current value.")
	ui.Output("  [v]alue <v>     : change value to <v>.")
	ui.Output("  [s]et <f>       : set <f to 1.")
	ui.Output("  [c]lear <f>     : clear <f> to 0.")
	ui.Output("  [w]rite <f> <v> : write val <v> into field <f>.")
	ui.Output("  <f>             : read the value of field <f>.")
	ui.Output("  [l]ist [0]      : list all offsets of '1's or '0's.")
	ui.Output("  exit, quit      : exit this program.")
	ui.Output("  \nTwo formats to represent filed:")
	ui.Output("  single bit  : like 1, 3, 0")
	ui.Output("  field range : like 0:3, 3:1")
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

	updateBinStr("0")

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
