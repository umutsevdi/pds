package yara

import (
	"fmt"
	"log"
	"os"
	"smtp_phishing_detection/model"

	"github.com/hillu/go-yara"
)

var yara_vars model.YaraVars

func YaraInit() {
	var err error
	fmt.Println("hello from init")

	yara_vars.Compiler, err = yara.NewCompiler()
	if err != nil {
		log.Fatal(err)
	}
	LoadYaraRules()
}

func LoadYaraRules() {
	path := "./yara/phishing.yara"
	fd, err := os.Open(path)
	log.Fatal(err)
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
}

func YaraScanMemory(message string) {
	var err error
	yara_vars.Matches, err = yara_vars.Rules.ScanMem([]byte(message), 0, 0)
	if err != nil {
		log.Fatal(err)
	}
}

func PrintMatches() {
	for _, v := range yara_vars.Matches {
		fmt.Println(v.Rule)
		for _, i := range v.Strings {
			fmt.Println("Possible phishing keyword", i.Name)
		}
	}
}
