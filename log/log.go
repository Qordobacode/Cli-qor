package log

import (
	"fmt"
)

// Level data type
var (
	IsVerbose = false
)

func Infof(format string, v ...interface{}) {
	fmt.Printf(format, v...)
}

func Debugf(format string, v ...interface{}) {
	if IsVerbose {
		fmt.Printf(format, v...)
	}
}
