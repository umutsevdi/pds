package yara

import (
	"fmt"
	"log"
	"os"
	"smtp_phishing_detection/model"

	"github.com/hillu/go-yara/v4"
)

var yara_vars model.YaraVars

func init() {
	var err error
	yara_vars.Compiler, err = yara.NewCompiler()
	if err != nil {
		log.Fatal(err)
	}
	LoadYaraRules()
}

func LoadYaraRules() {
	path := "./yara/phishing.yara"
	fd, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	err = yara_vars.Compiler.AddFile(fd, "scam_rules")
	if err != nil {
		log.Fatal(err)
	}
	fd.Close()
	YaraGetRules()
}

func YaraGetRules() {
	var err error
	yara_vars.Rules, err = yara_vars.Compiler.GetRules()
	if err != nil {
		log.Fatal(err)
	}
	// for _, v := range yara_vars.Rules.GetRules() {
	//	fmt.Println(v)
	// }
}

func YaraScanMemory(message string) {
	var err error
	err = yara_vars.Rules.ScanMem([]byte(message), 0, 0, &yara_vars.Matches)
	if err != nil {
		log.Fatal(err)
	}
}

func PrintMatches() {
	if yara_vars.Matches != nil {
		fmt.Println("None Phishing Keyword Detected")
	}
	for _, v := range yara_vars.Matches {
		fmt.Println(v.Rule)
		for _, i := range v.Strings {
			fmt.Println("Possible phishing keyword", i.Name)
		}
	}
	ClearMatches()
}

func ClearMatches() {
	yara_vars.Matches = nil
}
