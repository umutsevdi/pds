package model

import (
	"github.com/hillu/go-yara"
)

type YaraVars struct {
	Compiler *yara.Compiler
	Rules    *yara.Rules
	Matches  yara.MatchRules
}
