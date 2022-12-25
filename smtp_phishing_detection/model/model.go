package model

import "github.com/flier/gohs/hyperscan"

type HSrelated struct {
	Patterns *hyperscan.Pattern
	HsDB     hyperscan.BlockDatabase
	Scratch  *hyperscan.Scratch
}
