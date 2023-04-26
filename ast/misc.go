package ast

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Field struct {
	Key   Expr `json:"key"`
	Value Expr `json:"value"`
}

// todo: добавить экранирование для зарезервированных слов
func (f *Field) String() string {
	rhs := ""
	if f.Key != nil {
		rhs = fmt.Sprintf("%s = ", strings.Trim(f.Key.String(), "\""))
	}
	return fmt.Sprintf("%s%s", rhs, f.Value)
}

func (f *Field) UnmarshalJSON(bytes []byte) error {
	var temp struct {
		Key   json.RawMessage `json:"key"`
		Value json.RawMessage `json:"value"`
	}

	err := json.Unmarshal(bytes, &temp)
	if err != nil {
		return fmt.Errorf("field: failed to unmarshal: %w", err)
	}

	*f = Field{}

	f.Value, err = unmarshalExpr(temp.Value)
	if err != nil {
		return fmt.Errorf("field: failed to unmarshal value field: %w", err)
	}

	f.Key, err = unmarshalExpr(temp.Key)
	if err != nil {
		return fmt.Errorf("field: failed to unmarshal key field: %w", err)
	}

	return nil
}

type ParList struct {
	HasVargs bool     `json:"has_vargs"`
	Names    []string `json:"names"`
}

func (p *ParList) String() string {
	names := make([]string, len(p.Names))
	copy(names, p.Names)
	if p.HasVargs {
		names = append(names, "...")
	}
	return strings.TrimRight(strings.Join(names, ", "), ", ")
}

type FuncName struct {
	Func     Expr   `json:"func"`
	Receiver Expr   `json:"receiver"`
	Method   string `json:"method"`
}

func (f *FuncName) String() string {
	var body string
	if f.Receiver != nil {
		body = fmt.Sprintf("%s:%s", f.Receiver, f.Method)
	} else {
		body = f.Func.String()
	}
	return body
}

func (f *FuncName) UnmarshalJSON(bytes []byte) error {
	var temp struct {
		Func     json.RawMessage `json:"func"`
		Receiver json.RawMessage `json:"receiver"`
		Method   string          `json:"method"`
	}

	err := json.Unmarshal(bytes, &temp)
	if err != nil {
		return fmt.Errorf("func_name: failed to unmarshal: %w", err)
	}

	*f = FuncName{Method: temp.Method}

	f.Func, err = unmarshalExpr(temp.Func)
	if err != nil {
		return fmt.Errorf("func_name: failed to unmarshal func field: %w", err)
	}

	f.Receiver, err = unmarshalExpr(temp.Receiver)
	if err != nil {
		return fmt.Errorf("func_name: failed to unmarshal func receiver: %w", err)
	}

	return nil
}
