package main

import (
	"github.com/qordobacode/test/cmd"
	"log"
)

func main() {
	// if we crash the go code, we get the file name and line number
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	cmd.Execute()
}
