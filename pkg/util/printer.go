package util

import (
	"fmt"
	"io"
	"strings"
	"text/tabwriter"

	"github.com/fatih/color"
)

const (
	notype   = ""
	oktype   = "[OK]"
	errtype  = "[ERROR]"
	warntype = "[WARN]"
)

var green = color.New(color.FgGreen)
var red = color.New(color.FgRed)
var orange = color.New(color.FgRed, color.FgYellow)

// PrettyPrintOkf format output with tabs and color
func PrettyPrintOkf(out io.Writer, msg string, a ...interface{}) {
	print(out, msg, oktype, a...)
}

// PrettyPrintOk format output with tabs and color
func PrettyPrintOk(out io.Writer, msg string) {
	print(out, msg, oktype)
}

// PrettyPrintErrf format output with tabs and color
func PrettyPrintErrf(out io.Writer, msg string, a ...interface{}) {
	print(out, msg, errtype, a...)
}

// PrettyPrintErr format output with tabs and color
func PrettyPrintErr(out io.Writer, msg string) {
	print(out, msg, errtype)
}

// PrettyPrintWarnf format output with tabs and color
func PrettyPrintWarnf(out io.Writer, msg string, a ...interface{}) {
	print(out, msg, warntype, a...)
}

// PrettyPrintWarn format output with tabs and color
func PrettyPrintWarn(out io.Writer, msg string) {
	print(out, msg, warntype)
}

// PrettyPrintf format output with tabs and color
func PrettyPrintf(out io.Writer, msg string, a ...interface{}) {
	print(out, msg, notype, a...)
}

// PrettyPrint format output with tabs and color
func PrettyPrint(out io.Writer, msg string) {
	print(out, msg, notype)
}

// PrintErrorf print whole message
func PrintErrorf(out io.Writer, msg string, a ...interface{}) {
	printColor(out, msg, red, a...)
}

// PrintError print whole message
func PrintError(out io.Writer, msg string) {
	printColor(out, msg, red)
}

// PrintHeader will print header with = character
func PrintHeader(out io.Writer, msg string) {
	w := tabwriter.NewWriter(out, 84, 0, 0, '=', 0)
	fmt.Fprintln(w, "")
	format := msg + "\t\n"
	fmt.Fprintf(w, format)
	w.Flush()
}

func print(out io.Writer, msg, status string, a ...interface{}) {
	w := tabwriter.NewWriter(out, 80, 0, 0, ' ', 0)
	// print message
	format := msg + "\t"
	fmt.Fprintf(w, format, a...)
	// get correct color
	var clr *color.Color
	switch status {
	case oktype:
		clr = green
	case errtype:
		clr = red
	case warntype:
		clr = orange
	}
	// print status
	if status != notype {
		sformat := "%s\n"
		fmt.Fprintf(w, sformat, clr.SprintFunc()(status))
	}
	w.Flush()
}

func printColor(out io.Writer, msg string, clr *color.Color, a ...interface{}) {
	format := "%s"
	// Remove any newline, results in only one \n
	line := strings.Trim(fmt.Sprintf(format, clr.SprintfFunc()(msg, a...)), "\n") + "\n"
	fmt.Fprint(out, line)
}
