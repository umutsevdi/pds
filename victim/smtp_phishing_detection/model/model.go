package model

import (
	"github.com/hillu/go-yara/v4"
)

type YaraVars struct {
	Compiler *yara.Compiler
	Rules    *yara.Rules
	Matches  yara.MatchRules
}
