package ast

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"strings"
)

type Expr interface {
	fmt.Stringer
	PositionHolder
	exprMarker()
}

type ExprBase struct {
	Node
}

func (expr *ExprBase) exprMarker() {}

/* ConstExprs {{{ */

type ConstExpr interface {
	Expr
	constExprMarker()
}

type ConstExprBase struct {
	ExprBase
}

func (expr *ConstExprBase) constExprMarker() {}

type TrueExpr struct {
	ConstExprBase
}

func (t *TrueExpr) String() string { return "true" }

func (t *TrueExpr) MarshalJSON() ([]byte, error) {
	return sjson.SetBytes([]byte{}, DiscriminatorField, "true_expr")
}

type FalseExpr struct {
	ConstExprBase
}

func (f *FalseExpr) String() string { return "false" }

func (f *FalseExpr) MarshalJSON() ([]byte, error) {
	return sjson.SetBytes([]byte{}, DiscriminatorField, "false_expr")
}

type NilExpr struct {
	ConstExprBase
}

func (n *NilExpr) String() string { return "nil" }

func (n *NilExpr) MarshalJSON() ([]byte, error) {
	return sjson.SetBytes([]byte{}, DiscriminatorField, "nil_expr")
}

type NumberExpr struct {
	ConstExprBase

	Value string `json:"value"`
}

func (n *NumberExpr) String() string { return n.Value }

func (n *NumberExpr) MarshalJSON() ([]byte, error) {
	return marshalWithType(n, "number_expr")
}

type StringExpr struct {
	ConstExprBase

	Value string `json:"value"`
}

func (s *StringExpr) MarshalJSON() ([]byte, error) {
	return marshalWithType(s, "string_expr")
}

func (s *StringExpr) String() string {
	if strings.Contains(s.Value, "\n") {
		return fmt.Sprintf(`[[%s]]`, s.Value)
	}
	return fmt.Sprintf(`"%s"`, s.Value)
}

/* ConstExprs }}} */

type Comma3Expr struct {
	ExprBase
	AdjustRet bool `json:"adjust_ret"`
}

func (c *Comma3Expr) String() string { return "..." }

func (c *Comma3Expr) MarshalJSON() ([]byte, error) {
	return marshalWithType(c, "comma_3_expr")
}

type IdentExpr struct {
	ExprBase

	Value string `json:"value"`
}

func (i *IdentExpr) String() string { return i.Value }

func (i *IdentExpr) MarshalJSON() ([]byte, error) {
	return marshalWithType(i, "ident_expr")
}

type AttrGetExpr struct {
	ExprBase

	Object Expr `json:"object"`
	Key    Expr `json:"key"`
}

func (a *AttrGetExpr) String() string {
	switch a.Key.(type) {
	case *NumberExpr, *UnaryLenOpExpr, *UnaryMinusOpExpr:
		return fmt.Sprintf("%s[%s]", a.Object, a.Key)
	case *StringExpr:
		return fmt.Sprintf("%s.%s", a.Object, strings.Trim(a.Key.String(), "\""))
	default:
		return fmt.Sprintf("%s.%s", a.Object, a.Key)
	}
}

func (a *AttrGetExpr) UnmarshalJSON(bytes []byte) error {
	var temp struct {
		Object json.RawMessage `json:"object"`
		Key    json.RawMessage `json:"key"`
	}

	err := json.Unmarshal(bytes, &temp)
	if err != nil {
		return fmt.Errorf("attr_get_expr: failed to unmarshal: %w", err)
	}

	*a = AttrGetExpr{}

	a.Object, err = unmarshalExpr(temp.Object)
	if err != nil {
		return fmt.Errorf("attr_get_expr: failed to unmarshal object field: %w", err)
	}

	a.Key, err = unmarshalExpr(temp.Key)
	if err != nil {
		return fmt.Errorf("attr_get_expr: failed to unmarshal key field: %w", err)
	}

	return nil
}

func (a *AttrGetExpr) MarshalJSON() ([]byte, error) {
	return marshalWithType(a, "attr_get_expr")
}

type TableExpr struct {
	ExprBase

	Fields []*Field `json:"fields"`
}

func (t *TableExpr) String() string {
	var fields []string
	for _, field := range t.Fields {
		fields = append(fields, field.String())
	}
	return fmt.Sprintf(`{%s}`, strings.TrimRight(strings.Join(fields, ", "), ", "))
}

func (t *TableExpr) MarshalJSON() ([]byte, error) {
	return marshalWithType(t, "table_expr")
}

type FuncCallExpr struct {
	ExprBase

	Func      Expr   `json:"func"`
	Receiver  Expr   `json:"receiver"`
	Method    string `json:"method"`
	Args      []Expr `json:"args"`
	AdjustRet bool   `json:"adjust_ret"`
}

func (f *FuncCallExpr) String() string {
	var body string
	if f.Receiver != nil {
		receiver := "%s"
		if _, ok := f.Receiver.(*StringExpr); ok {
			receiver = "(%s)"
		}
		body = fmt.Sprintf("%s:%s", fmt.Sprintf(receiver, f.Receiver), f.Method)
	} else {
		body = f.Func.String()
	}
	var args []string
	for _, arg := range f.Args {
		args = append(args, arg.String())
	}
	return fmt.Sprintf("%s(%s)", body, strings.TrimRight(strings.Join(args, ", "), ", "))
}

func (f *FuncCallExpr) UnmarshalJSON(bytes []byte) error {
	var temp struct {
		Func      json.RawMessage   `json:"func"`
		Receiver  json.RawMessage   `json:"receiver"`
		Method    string            `json:"method"`
		Args      []json.RawMessage `json:"args"`
		AdjustRet bool              `json:"adjust_ret"`
	}

	err := json.Unmarshal(bytes, &temp)
	if err != nil {
		return fmt.Errorf("func_call_expr: failed to unmarshal: %w", err)
	}

	*f = FuncCallExpr{Method: temp.Method, AdjustRet: temp.AdjustRet}

	f.Func, err = unmarshalExpr(temp.Func)
	if err != nil {
		return fmt.Errorf("func_call_expr: failed to unmarshal field func: %w", err)
	}

	f.Receiver, err = unmarshalExpr(temp.Receiver)
	if err != nil {
		return fmt.Errorf("func_call_expr: failed to unmarshal field receiver: %w", err)
	}

	for i, arg := range temp.Args {
		e, err := unmarshalExpr(arg)
		if err != nil {
			return fmt.Errorf("func_call_expr: failed to unmarshal %d argument: %w", i, err)
		}
		f.Args = append(f.Args, e)
	}

	return nil
}

func (f *FuncCallExpr) MarshalJSON() ([]byte, error) {
	return marshalWithType(f, "func_call_expr")
}

type LogicalOpExpr struct {
	ExprBase

	Operator string `json:"operator"`
	Lhs      Expr   `json:"lhs"`
	Rhs      Expr   `json:"rhs"`
}

func (l *LogicalOpExpr) String() string {
	var lhs, rhs string
	if v, ok := l.Lhs.(*LogicalOpExpr); ok && v.Operator != l.Operator {
		lhs = fmt.Sprintf("(%s)", l.Lhs)
	} else {
		lhs = l.Lhs.String()
	}
	if v, ok := l.Rhs.(*LogicalOpExpr); ok && v.Operator != l.Operator {
		rhs = fmt.Sprintf("(%s)", l.Rhs)
	} else {
		rhs = l.Rhs.String()
	}
	return fmt.Sprintf("%s %s %s", lhs, l.Operator, rhs)
}

func (l *LogicalOpExpr) UnmarshalJSON(bytes []byte) error {
	var temp struct {
		Operator string          `json:"operator"`
		Lhs      json.RawMessage `json:"lhs"`
		Rhs      json.RawMessage `json:"rhs"`
	}

	err := json.Unmarshal(bytes, &temp)
	if err != nil {
		return fmt.Errorf("logical_op_expr: failed to unmarshal: %w", err)
	}

	// todo: make constant operators
	*l = LogicalOpExpr{Operator: temp.Operator}

	l.Lhs, err = unmarshalExpr(temp.Lhs)
	if err != nil {
		return fmt.Errorf("logical_op_expr: failed to unmarshal field lhs: %w", err)
	}

	l.Rhs, err = unmarshalExpr(temp.Rhs)
	if err != nil {
		return fmt.Errorf("logical_op_expr: failed to unmarshal field rhs: %w", err)
	}

	return nil
}

func (l *LogicalOpExpr) MarshalJSON() ([]byte, error) {
	return marshalWithType(l, "logical_op_expr")
}

type RelationalOpExpr struct {
	ExprBase

	Operator string `json:"operator"`
	Lhs      Expr   `json:"lhs"`
	Rhs      Expr   `json:"rhs"`
}

func (r *RelationalOpExpr) String() string { return fmt.Sprintf("%s %s %s", r.Lhs, r.Operator, r.Rhs) }

func (r *RelationalOpExpr) UnmarshalJSON(bytes []byte) error {
	var temp struct {
		Operator string          `json:"operator"`
		Lhs      json.RawMessage `json:"lhs"`
		Rhs      json.RawMessage `json:"rhs"`
	}

	err := json.Unmarshal(bytes, &temp)
	if err != nil {
		return fmt.Errorf("relational_op_expr: failed to unmarshal: %w", err)
	}

	// todo: make constant operators
	*r = RelationalOpExpr{Operator: temp.Operator}

	r.Lhs, err = unmarshalExpr(temp.Lhs)
	if err != nil {
		return fmt.Errorf("relational_op_expr: failed to unmarshal field lhs: %w", err)
	}

	r.Rhs, err = unmarshalExpr(temp.Rhs)
	if err != nil {
		return fmt.Errorf("relational_op_expr: failed to unmarshal field rhs: %w", err)
	}

	return nil
}

func (r *RelationalOpExpr) MarshalJSON() ([]byte, error) {
	return marshalWithType(r, "relational_op_expr")
}

type StringConcatOpExpr struct {
	ExprBase

	Lhs Expr `json:"lhs"`
	Rhs Expr `json:"rhs"`
}

func (s *StringConcatOpExpr) String() string { return fmt.Sprintf("%s..%s", s.Lhs, s.Rhs) }

func (s *StringConcatOpExpr) UnmarshalJSON(bytes []byte) error {
	var temp struct {
		Lhs json.RawMessage `json:"lhs"`
		Rhs json.RawMessage `json:"rhs"`
	}

	err := json.Unmarshal(bytes, &temp)
	if err != nil {
		return fmt.Errorf("string_concat_op_expr: failed to unmarshal: %w", err)
	}

	*s = StringConcatOpExpr{}

	s.Lhs, err = unmarshalExpr(temp.Lhs)
	if err != nil {
		return fmt.Errorf("string_concat_op_expr: failed to unmarshal field lhs: %w", err)
	}

	s.Rhs, err = unmarshalExpr(temp.Rhs)
	if err != nil {
		return fmt.Errorf("string_concat_op_expr: failed to unmarshal field rhs: %w", err)
	}

	return nil
}

func (s *StringConcatOpExpr) MarshalJSON() ([]byte, error) {
	return marshalWithType(s, "string_concat_op_expr")
}

type ArithmeticOpExpr struct {
	ExprBase

	Operator string `json:"operator"`
	Lhs      Expr   `json:"lhs"`
	Rhs      Expr   `json:"rhs"`
}

func (a *ArithmeticOpExpr) String() string { return fmt.Sprintf("%s %s %s", a.Lhs, a.Operator, a.Rhs) }

func (a *ArithmeticOpExpr) UnmarshalJSON(bytes []byte) error {
	var temp struct {
		Operator string          `json:"operator"`
		Lhs      json.RawMessage `json:"lhs"`
		Rhs      json.RawMessage `json:"rhs"`
	}

	err := json.Unmarshal(bytes, &temp)
	if err != nil {
		return fmt.Errorf("arithmetic_op_expr: failed to unmarshal: %w", err)
	}

	// todo: make constant operators
	*a = ArithmeticOpExpr{Operator: temp.Operator}

	a.Lhs, err = unmarshalExpr(temp.Lhs)
	if err != nil {
		return fmt.Errorf("arithmetic_op_expr: failed to unmarshal field lhs: %w", err)
	}

	a.Rhs, err = unmarshalExpr(temp.Rhs)
	if err != nil {
		return fmt.Errorf("arithmetic_op_expr: failed to unmarshal field rhs: %w", err)
	}

	return nil
}

func (a *ArithmeticOpExpr) MarshalJSON() ([]byte, error) {
	return marshalWithType(a, "arithmetic_op_expr")
}

type UnaryMinusOpExpr struct {
	ExprBase
	Expr Expr `json:"expr"`
}

func (u *UnaryMinusOpExpr) String() string { return fmt.Sprintf("-%s", u.Expr) }

func (u *UnaryMinusOpExpr) UnmarshalJSON(bytes []byte) error {
	var temp struct {
		Expr json.RawMessage `json:"expr"`
	}

	err := json.Unmarshal(bytes, &temp)
	if err != nil {
		return fmt.Errorf("unary_minus_op_expr: failed to unmarshal: %w", err)
	}

	*u = UnaryMinusOpExpr{}

	u.Expr, err = unmarshalExpr(temp.Expr)
	if err != nil {
		return fmt.Errorf("unary_minus_op_expr: failed to unmarshal field expr: %w", err)
	}

	return nil
}

func (u *UnaryMinusOpExpr) MarshalJSON() ([]byte, error) {
	return marshalWithType(u, "unary_minus_op_expr")
}

type UnaryNotOpExpr struct {
	ExprBase
	Expr Expr `json:"expr"`
}

func (u *UnaryNotOpExpr) String() string {
	switch u.Expr.(type) {
	case *LogicalOpExpr:
		return fmt.Sprintf("not (%s)", u.Expr)
	default:
		return fmt.Sprintf("not %s", u.Expr)
	}
}

func (u *UnaryNotOpExpr) UnmarshalJSON(bytes []byte) error {
	var temp struct {
		Expr json.RawMessage `json:"expr"`
	}

	err := json.Unmarshal(bytes, &temp)
	if err != nil {
		return fmt.Errorf("unary_not_op_expr: failed to unmarshal: %w", err)
	}

	*u = UnaryNotOpExpr{}

	u.Expr, err = unmarshalExpr(temp.Expr)
	if err != nil {
		return fmt.Errorf("unary_not_op_expr: failed to unmarshal field expr: %w", err)
	}

	return nil
}

func (u *UnaryNotOpExpr) MarshalJSON() ([]byte, error) {
	return marshalWithType(u, "unary_not_op_expr")
}

type UnaryLenOpExpr struct {
	ExprBase
	Expr Expr `json:"expr"`
}

func (u *UnaryLenOpExpr) String() string { return fmt.Sprintf("#%s", u.Expr) }

func (u *UnaryLenOpExpr) UnmarshalJSON(bytes []byte) error {
	var temp struct {
		Expr json.RawMessage `json:"expr"`
	}

	err := json.Unmarshal(bytes, &temp)
	if err != nil {
		return fmt.Errorf("unary_len_op_expr: failed to unmarshal: %w", err)
	}

	*u = UnaryLenOpExpr{}

	u.Expr, err = unmarshalExpr(temp.Expr)
	if err != nil {
		return fmt.Errorf("unary_len_op_expr: failed to unmarshal field expr: %w", err)
	}

	return nil
}

func (u *UnaryLenOpExpr) MarshalJSON() ([]byte, error) {
	return marshalWithType(u, "unary_len_op_expr")
}

type FunctionExpr struct {
	ExprBase

	ParList *ParList `json:"par_list"`
	Stmts   []Stmt   `json:"stmts"`
}

func (f *FunctionExpr) String() string { return f.StringIndent(0) }

func (f *FunctionExpr) StringIndent(indent int) string {
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("(%s)\n", f.ParList))
	for _, stmt := range f.Stmts {
		builder.WriteString(stmt.StringIndent(indent + 1))
	}
	builder.WriteString(fmt.Sprintf("%send\n", strings.Repeat("\t", indent)))
	return builder.String()
}

func (f *FunctionExpr) UnmarshalJSON(bytes []byte) error {
	var temp struct {
		ParList *ParList          `json:"par_list"`
		Stmts   []json.RawMessage `json:"stmts"`
	}

	err := json.Unmarshal(bytes, &temp)
	if err != nil {
		return fmt.Errorf("function_expr: failed to unmarshal: %w", err)
	}

	*f = FunctionExpr{ParList: temp.ParList}

	for i, s := range temp.Stmts {
		stmt, err := unmarshalStmt(s)
		if err != nil {
			return fmt.Errorf("function_expr: failed to unmarshal %d statement: %w", i, err)
		}
		f.Stmts = append(f.Stmts, stmt)
	}

	return nil
}

func (f *FunctionExpr) MarshalJSON() ([]byte, error) {
	return marshalWithType(f, "function_expr")
}

func unmarshalExpr(data []byte) (Expr, error) {
	if bytes.Equal(data, []byte(`null`)) {
		return nil, nil
	}
	t := gjson.GetBytes(data, DiscriminatorField)
	if !t.Exists() {
		return nil, fmt.Errorf("expr unmarshal error: type discriminator not found")
	}
	var e Expr
	switch t.String() {
	case "true_expr":
		return &TrueExpr{}, nil
	case "false_expr":
		return &FalseExpr{}, nil
	case "nil_expr":
		return &NilExpr{}, nil
	case "number_expr":
		e = &NumberExpr{}
	case "string_expr":
		e = &StringExpr{}
	case "comma_3_expr":
		e = &Comma3Expr{}
	case "ident_expr":
		e = &IdentExpr{}
	case "attr_get_expr":
		e = &AttrGetExpr{}
	case "table_expr":
		e = &TableExpr{}
	case "func_call_expr":
		e = &FuncCallExpr{}
	case "logical_op_expr":
		e = &LogicalOpExpr{}
	case "relational_op_expr":
		e = &RelationalOpExpr{}
	case "string_concat_op_expr":
		e = &StringConcatOpExpr{}
	case "arithmetic_op_expr":
		e = &ArithmeticOpExpr{}
	case "unary_minus_op_expr":
		e = &UnaryMinusOpExpr{}
	case "unary_not_op_expr":
		e = &UnaryNotOpExpr{}
	case "unary_len_op_expr":
		e = &UnaryLenOpExpr{}
	case "function_expr":
		e = &FunctionExpr{}
	default:
		return nil, fmt.Errorf("unsupported expr: %s", t.String())
	}

	if err := json.Unmarshal(data, &e); err != nil {
		return nil, fmt.Errorf("expr unmarshal error: %w", err)
	}
	return e, nil
}
