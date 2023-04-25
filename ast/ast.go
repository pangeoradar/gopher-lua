package ast

import (
	"encoding/json"
	"fmt"
	"github.com/tidwall/sjson"
	"reflect"
	"strings"
)

type Printer interface {
	Print(ident int) string
}

type PositionHolder interface {
	Line() int
	SetLine(int)
	LastLine() int
	SetLastLine(int)
}

type Node struct {
	line     int
	lastline int
}

func (self *Node) Line() int {
	return self.line
}

func (self *Node) SetLine(line int) {
	self.line = line
}

func (self *Node) LastLine() int {
	return self.lastline
}

func (self *Node) SetLastLine(line int) {
	self.lastline = line
}

const DiscriminatorField = "_type"

func marshalWithType(o interface{}, t string) ([]byte, error) {
	val := reflect.ValueOf(o)
	if val.Kind() != reflect.Pointer {
		return nil, fmt.Errorf("marshal_with_type: object should be a pointer, got: %T", o)
	}

	data, err := json.Marshal(val.Elem().Interface())
	if err != nil {
		return nil, fmt.Errorf("marshal_with_type: marshal error: %w", err)
	}
	return sjson.SetBytes(data, DiscriminatorField, t)
}

func ParseRule(bytes []byte) ([]Stmt, error) {
	var lines []json.RawMessage
	if err := json.Unmarshal(bytes, &lines); err != nil {
		return nil, fmt.Errorf("parse rule: incorrect format: %w", err)
	}

	var stmts []Stmt
	for i, s := range lines {
		stmt, err := unmarshalStmt(s)
		if err != nil {
			return nil, fmt.Errorf("parse rule: %d statement error: %w", i, err)
		}
		stmts = append(stmts, stmt)
	}

	return stmts, nil
}

func PrintRule(chunks []Stmt) string {
	builder := strings.Builder{}
	for _, s := range chunks {
		builder.WriteString(s.StringIndent(0))
	}
	return builder.String()
}
