package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/choueric/readline"
)

var (
	inst   *readline.Instance
	regLen = 32
	binStr string
)

// field range
type fRange struct {
	start int
	end   int
}

// setFieldRangeOfBinStr sets to 1 or clears to 0 with the filed range of binStr
// rStr is a filed range string
func setFieldOfBinStr(rStr string, set bool) {
	r, err := getRange(rStr)
	if err != nil {
		inst.Error(fmt.Sprintf("parse field start index failed, %v\n", err))
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
		inst.Error(fmt.Sprintf("convert to Int failed: %v\n", err))
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
		inst.Error(fmt.Sprintf("parse field start index failed, %v\n", err))
		return
	}

	val, err := parseInt(vStr)
	if err != nil {
		inst.Error(fmt.Sprintf("convert to Int failed: %v\n", err))
		return
	}

	max := (2 << uint(r.end-r.start)) - 1
	if val < 0 || int(val) > max {
		inst.Error(fmt.Sprint("val is out of range [%d, %d]\n", 0, max))
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
		inst.Error(fmt.Sprintf("getFieldStr failed: %v\n", err))
		return
	}

	outputTriFormat(os.Stdout, subbin)
}

func executeCmdline(input string, data interface{}) (exit bool) {
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
			inst.Error("Needs argument: <field>\n")
			return
		}
		updateBinStr(cmdline[1])
		outputTriFormat(os.Stdout, binStr)
	case "set", "s", "clear", "c":
		if len(cmdline) < 2 {
			inst.Error("Needs argument: <field>\n")
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
			inst.Error("Needs arguments: <field> <val>\n")
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
	inst.Print("Usage:\n")
	inst.Print("  [h]elp          : print this message.\n")
	inst.Print("  [p]rint         : show current value.\n")
	inst.Print("  [v]alue <v>     : change value to <v>.\n")
	inst.Print("  [s]et <f>       : set <f to 1.\n")
	inst.Print("  [c]lear <f>     : clear <f> to 0.\n")
	inst.Print("  [w]rite <f> <v> : write val <v> into field <f>.\n")
	inst.Print("  <f>             : read the value of field <f>.\n")
	inst.Print("  [l]ist [0]      : list all offsets of '1's or '0's.\n")
	inst.Print("  exit, quit      : exit this program.\n")
	inst.Print("  \nTwo formats to represent filed:\n")
	inst.Print("  single bit  : like 1, 3, 0\n")
	inst.Print("  field range : like 0:3, 3:1\n")
}

func main() {
	var debug bool
	flag.BoolVar(&debug, "d", false, "enable debug")
	flag.IntVar(&regLen, "l", 32, "register length.")
	flag.Parse()
	if len(os.Args) == 2 && (os.Args[1] == "help" || os.Args[1] == "-h") {
		flag.Usage()
		return
	}

	_inst, err := readline.New("\033[32m>>\033[0m ")
	if err != nil {
		fmt.Println(err)
		return
	}
	inst = _inst
	inst.Debug = debug
	defer inst.Destroy()

	inst.SetExecute(executeCmdline, nil)
	inst.SetCompleter(
		readline.Cmd("print"),
		readline.Cmd("value"),
		readline.Cmd("set"),
		readline.Cmd("clear"),
		readline.Cmd("write"),
		readline.Cmd("list"),
		readline.Cmd("help"),
		readline.Cmd("exit"),
		readline.Cmd("quit"),
	)

	updateBinStr("0")

	readline.InputLoop(inst)
}
