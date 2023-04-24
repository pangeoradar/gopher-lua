package ast

import (
	"fmt"
)

type Position struct {
	Source string `json:"source"`
	Line   int    `json:"line"`
	Column int    `json:"column"`
}

type Token struct {
	Type int      `json:"type"`
	Name string   `json:"name"`
	Str  string   `json:"str"`
	Pos  Position `json:"pos"`
}

func (self *Token) String() string {
	return fmt.Sprintf("<type:%v, str:%v>", self.Name, self.Str)
}
