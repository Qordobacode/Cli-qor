package log

import (
	"fmt"
	"os"
	"time"
)

// Level data type
type Level int

// Level data type
var (
	IsVerbose = false
)

// Error logs error message
func Error(v ...interface{}) {
	fmt.Fprintf(os.Stderr, fmt.Sprint(v...)+"\n")
}

// Errorf logs formatted error message
func Errorf(format string, v ...interface{}) {
	fmt.Fprintf(os.Stderr, fmt.Sprintf(format+"\n", v...))
}

// Info logs informational message
func Info(v ...interface{}) {
	out(fmt.Sprint(v...))
}

// Infof logs formatted informational message
func Infof(format string, v ...interface{}) {
	out(format, v...)
}

// Debug logs debug info
func Debug(v ...interface{}) {
	if IsVerbose {
		out(fmt.Sprint(v...))
	}
}

// Debugf logs formatted debug info
func Debugf(format string, v ...interface{}) {
	if IsVerbose {
		out(format, v...)
	}
}

func out(format string, v ...interface{}) {
	if len(v) > 0 {
		template := format + "\n"
		fmt.Printf(template, v...)
	} else {
		fmt.Printf("%s\n", format)
	}
}

// TimeTrack function for printing function time in log
func TimeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	Infof("---%s took %s---", name, elapsed)
}
