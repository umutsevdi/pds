package hs

import (
	"fmt"
	"os"
	"smtp_phishing_detection/model"

	"github.com/flier/gohs/hyperscan"
)

var hsVar model.HSrelated

func Init() {
	fmt.Println("hello from init")
	var err error
	patternStr := "test"
	hsVar.Patterns = hyperscan.NewPattern(patternStr, hyperscan.DotAll|hyperscan.SomLeftMost)
	hsVar.HsDB, err = hyperscan.NewBlockDatabase(hsVar.Patterns)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Unable to compile pattern \"%s\": %s\n", hsVar.Patterns.String(), err.Error())
		os.Exit(-1)
	}

	defer hsVar.HsDB.Close()

	hsVar.Scratch, err = hyperscan.NewScratch(hsVar.HsDB)

	defer hsVar.Scratch.Free()

}

func CheckForPossiblePhishing(buff string) {

	if err := hsVar.HsDB.Scan([]byte(buff), hsVar.Scratch, nil, buff); err != nil {
		fmt.Fprint(os.Stderr, "ERROR: Unable to scan input buffer. Exiting.\n")
		os.Exit(-1)
	}

}
