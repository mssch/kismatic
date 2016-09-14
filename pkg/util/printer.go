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

// PrettyPrintOkf [OK](Green) with formatted string
func PrettyPrintOkf(out io.Writer, msg string, a ...interface{}) {
	print(out, msg, oktype, a...)
}

// PrettyPrintOk [OK](Green)
func PrettyPrintOk(out io.Writer, msg string) {
	print(out, msg, oktype)
}

// PrettyPrintErrf [ERROR](Red) with formatted string
func PrettyPrintErrf(out io.Writer, msg string, a ...interface{}) {
	print(out, msg, errtype, a...)
}

// PrettyPrintErr [ERROR](Red)
func PrettyPrintErr(out io.Writer, msg string) {
	print(out, msg, errtype)
}

// PrettyPrintWarnf [WARNING](Orange) with formatted string
func PrettyPrintWarnf(out io.Writer, msg string, a ...interface{}) {
	print(out, msg, warntype, a...)
}

// PrettyPrintWarn [WARNING](Orange)
func PrettyPrintWarn(out io.Writer, msg string) {
	print(out, msg, warntype)
}

// PrettyPrintf no type will be displayed, used for just single line printing
func PrettyPrintf(out io.Writer, msg string, a ...interface{}) {
	print(out, msg, notype, a...)
}

// PrettyPrint no type will be displayed, used for just single line printing
func PrettyPrint(out io.Writer, msg string) {
	print(out, msg, notype)
}

// PrintErrorf print whole message in error(Red) format
func PrintErrorf(out io.Writer, msg string, a ...interface{}) {
	printColor(out, msg, red, a...)
}

// PrintError print whole message in error(Red) format
func PrintError(out io.Writer, msg string) {
	printColor(out, msg, red)
}

// PrintHeader will print header with predifined width
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

	// print status
	if status != notype {
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
