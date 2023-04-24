package ast

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
)

type Stmt interface {
	PositionHolder
	stmtMarker()
}

type StmtBase struct {
	Node
}

func (stmt *StmtBase) stmtMarker() {}

type AssignStmt struct {
	StmtBase

	Lhs []Expr `json:"lhs"`
	Rhs []Expr `json:"rhs"`
}

func (a *AssignStmt) UnmarshalJSON(bytes []byte) error {
	var temp struct {
		Lhs []json.RawMessage `json:"lhs"`
		Rhs []json.RawMessage `json:"rhs"`
	}

	err := json.Unmarshal(bytes, &temp)
	if err != nil {
		return fmt.Errorf("assign_stmt: failed to unmarshal: %w", err)
	}

	*a = AssignStmt{}

	for i, e := range temp.Lhs {
		expr, err := unmarshalExpr(e)
		if err != nil {
			return fmt.Errorf("assign_stmt: failed to unmarshal %d lhs expr: %w", i, err)
		}
		a.Lhs = append(a.Lhs, expr)
	}

	for i, e := range temp.Rhs {
		expr, err := unmarshalExpr(e)
		if err != nil {
			return fmt.Errorf("assign_stmt: failed to unmarshal %d rhs expr: %w", i, err)
		}
		a.Rhs = append(a.Rhs, expr)
	}

	return nil
}

func (a *AssignStmt) MarshalJSON() ([]byte, error) {
	return marshalWithType(a, "assign_stmt")
}

type LocalAssignStmt struct {
	StmtBase

	Names []string `json:"names"`
	Exprs []Expr   `json:"exprs"`
}

func (l *LocalAssignStmt) UnmarshalJSON(bytes []byte) error {
	var temp struct {
		Names []string          `json:"names"`
		Exprs []json.RawMessage `json:"exprs"`
	}

	err := json.Unmarshal(bytes, &temp)
	if err != nil {
		return fmt.Errorf("local_assign_stmt: failed to unmarshal: %w", err)
	}

	*l = LocalAssignStmt{Names: temp.Names}

	for i, e := range temp.Exprs {
		expr, err := unmarshalExpr(e)
		if err != nil {
			return fmt.Errorf("local_assign_stmt: failed to unmarshal %d expr: %w", i, err)
		}
		l.Exprs = append(l.Exprs, expr)
	}

	return nil
}

func (l *LocalAssignStmt) MarshalJSON() ([]byte, error) {
	return marshalWithType(l, "local_assign_stmt")
}

type FuncCallStmt struct {
	StmtBase

	Expr Expr `json:"expr"`
}

func (f *FuncCallStmt) UnmarshalJSON(bytes []byte) error {
	var temp struct {
		Expr json.RawMessage `json:"expr"`
	}

	err := json.Unmarshal(bytes, &temp)
	if err != nil {
		return fmt.Errorf("func_call_stmt: failed to unmarshal: %w", err)
	}

	*f = FuncCallStmt{}

	f.Expr, err = unmarshalExpr(temp.Expr)
	if err != nil {
		return fmt.Errorf("func_call_stmt: failed to unmarshal expr: %w", err)
	}

	return nil
}

func (f *FuncCallStmt) MarshalJSON() ([]byte, error) {
	return marshalWithType(f, "func_call_stmt")
}

type DoBlockStmt struct {
	StmtBase

	Stmts []Stmt `json:"stmts"`
}

func (d *DoBlockStmt) UnmarshalJSON(bytes []byte) error {
	var temp struct {
		Stmts []json.RawMessage `json:"stmts"`
	}

	err := json.Unmarshal(bytes, &temp)
	if err != nil {
		return fmt.Errorf("do_block_stmt: failed to unmarshal: %w", err)
	}

	*d = DoBlockStmt{}

	for k, s := range temp.Stmts {
		stmt, err := unmarshalStmt(s)
		if err != nil {
			return fmt.Errorf("do_block_stmt: failed to unmarshal %d statement: %w", k, err)
		}
		d.Stmts = append(d.Stmts, stmt)
	}

	return nil
}

func (d *DoBlockStmt) MarshalJSON() ([]byte, error) {
	return marshalWithType(d, "do_block_stmt")
}

type WhileStmt struct {
	StmtBase

	Condition Expr   `json:"condition"`
	Stmts     []Stmt `json:"stmts"`
}

func (w *WhileStmt) UnmarshalJSON(bytes []byte) error {
	var temp struct {
		Condition json.RawMessage   `json:"condition"`
		Stmts     []json.RawMessage `json:"stmts"`
	}

	err := json.Unmarshal(bytes, &temp)
	if err != nil {
		return fmt.Errorf("while_stmt: failed to unmarshal: %w", err)
	}

	*w = WhileStmt{}

	for k, s := range temp.Stmts {
		stmt, err := unmarshalStmt(s)
		if err != nil {
			return fmt.Errorf("while_stmt: failed to unmarshal %d statement: %w", k, err)
		}
		w.Stmts = append(w.Stmts, stmt)
	}

	w.Condition, err = unmarshalExpr(temp.Condition)
	if err != nil {
		return fmt.Errorf("while_stmt: failed to unmarshal expr: %w", err)
	}

	return nil
}

func (w *WhileStmt) MarshalJSON() ([]byte, error) {
	return marshalWithType(w, "while_stmt")
}

type RepeatStmt struct {
	StmtBase

	Condition Expr   `json:"condition"`
	Stmts     []Stmt `json:"stmts"`
}

func (r *RepeatStmt) UnmarshalJSON(bytes []byte) error {
	var temp struct {
		Condition json.RawMessage   `json:"condition"`
		Stmts     []json.RawMessage `json:"stmts"`
	}

	err := json.Unmarshal(bytes, &temp)
	if err != nil {
		return fmt.Errorf("repeat_stmt: failed to unmarshal: %w", err)
	}

	*r = RepeatStmt{}

	for k, s := range temp.Stmts {
		stmt, err := unmarshalStmt(s)
		if err != nil {
			return fmt.Errorf("repeat_stmt: failed to unmarshal %d statement: %w", k, err)
		}
		r.Stmts = append(r.Stmts, stmt)
	}

	r.Condition, err = unmarshalExpr(temp.Condition)
	if err != nil {
		return fmt.Errorf("repeat_stmt: failed to unmarshal expr: %w", err)
	}

	return nil
}

func (r *RepeatStmt) MarshalJSON() ([]byte, error) {
	return marshalWithType(r, "repeat_stmt")
}

type IfStmt struct {
	StmtBase

	Condition Expr   `json:"condition"`
	Then      []Stmt `json:"then"`
	Else      []Stmt `json:"else"`
}

func (i *IfStmt) UnmarshalJSON(bytes []byte) error {
	var temp struct {
		Condition json.RawMessage   `json:"condition"`
		Then      []json.RawMessage `json:"then"`
		Else      []json.RawMessage `json:"else"`
	}

	err := json.Unmarshal(bytes, &temp)
	if err != nil {
		return fmt.Errorf("if_stmt: failed to unmarshal: %w", err)
	}

	*i = IfStmt{}

	for k, s := range temp.Then {
		stmt, err := unmarshalStmt(s)
		if err != nil {
			return fmt.Errorf("if_stmt: failed to unmarshal %d then statement: %w", k, err)
		}
		i.Then = append(i.Then, stmt)
	}

	for k, s := range temp.Else {
		stmt, err := unmarshalStmt(s)
		if err != nil {
			return fmt.Errorf("if_stmt: failed to unmarshal %d else statement: %w", k, err)
		}
		i.Else = append(i.Else, stmt)
	}

	i.Condition, err = unmarshalExpr(temp.Condition)
	if err != nil {
		return fmt.Errorf("if_stmt: failed to unmarshal expr: %w", err)
	}

	return nil
}

func (i *IfStmt) MarshalJSON() ([]byte, error) {
	return marshalWithType(i, "if_stmt")
}

type NumberForStmt struct {
	StmtBase

	Name  string `json:"name"`
	Init  Expr   `json:"init"`
	Limit Expr   `json:"limit"`
	Step  Expr   `json:"step"`
	Stmts []Stmt `json:"stmts"`
}

func (n *NumberForStmt) MarshalJSON() ([]byte, error) {
	return marshalWithType(n, "number_for_stmt")
}

func (n *NumberForStmt) UnmarshalJSON(bytes []byte) error {
	var temp struct {
		Name  string            `json:"name"`
		Init  json.RawMessage   `json:"init"`
		Limit json.RawMessage   `json:"limit"`
		Step  json.RawMessage   `json:"step"`
		Stmts []json.RawMessage `json:"stmts"`
	}

	err := json.Unmarshal(bytes, &temp)
	if err != nil {
		return fmt.Errorf("number_for_stmt: failed to unmarshal: %w", err)
	}

	*n = NumberForStmt{Name: temp.Name}

	for i, s := range temp.Stmts {
		stmt, err := unmarshalStmt(s)
		if err != nil {
			return fmt.Errorf("number_for_stmt: failed to unmarshal %d statement: %w", i, err)
		}
		n.Stmts = append(n.Stmts, stmt)
	}

	n.Limit, err = unmarshalExpr(temp.Limit)
	if err != nil {
		return fmt.Errorf("number_for_stmt: failed to unmarshal limit field: %w", err)
	}

	n.Init, err = unmarshalExpr(temp.Init)
	if err != nil {
		return fmt.Errorf("number_for_stmt: failed to unmarshal init field: %w", err)
	}

	n.Step, err = unmarshalExpr(temp.Step)
	if err != nil {
		return fmt.Errorf("number_for_stmt: failed to unmarshal step field: %w", err)
	}

	return nil
}

type GenericForStmt struct {
	StmtBase

	Names []string `json:"names"`
	Exprs []Expr   `json:"exprs"`
	Stmts []Stmt   `json:"stmts"`
}

func (g *GenericForStmt) MarshalJSON() ([]byte, error) {
	return marshalWithType(g, "generic_for_stmt")
}

func (g *GenericForStmt) UnmarshalJSON(bytes []byte) error {
	var temp struct {
		Names []string          `json:"names"`
		Exprs []json.RawMessage `json:"exprs"`
		Stmts []json.RawMessage `json:"stmts"`
	}

	err := json.Unmarshal(bytes, &temp)
	if err != nil {
		return fmt.Errorf("generic_for_stmt: failed to unmarshal: %w", err)
	}

	*g = GenericForStmt{Names: temp.Names}

	for i, s := range temp.Stmts {
		stmt, err := unmarshalStmt(s)
		if err != nil {
			return fmt.Errorf("generic_for_stmt: failed to unmarshal %d statement: %w", i, err)
		}
		g.Stmts = append(g.Stmts, stmt)
	}

	for i, e := range temp.Exprs {
		expr, err := unmarshalExpr(e)
		if err != nil {
			return fmt.Errorf("generic_for_stmt: failed to unmarshal %d expr: %w", i, err)
		}
		g.Exprs = append(g.Exprs, expr)
	}

	return nil
}

type FuncDefStmt struct {
	StmtBase

	Name *FuncName     `json:"name"`
	Func *FunctionExpr `json:"func"`
}

func (f *FuncDefStmt) MarshalJSON() ([]byte, error) {
	return marshalWithType(f, "func_def_stmt")
}

type ReturnStmt struct {
	StmtBase

	Exprs []Expr `json:"exprs"`
}

func (r *ReturnStmt) MarshalJSON() ([]byte, error) {
	return marshalWithType(r, "return_stmt")
}

func (r *ReturnStmt) UnmarshalJSON(bytes []byte) error {
	var temp struct {
		Exprs []json.RawMessage `json:"exprs"`
	}

	err := json.Unmarshal(bytes, &temp)
	if err != nil {
		return fmt.Errorf("return_stmt: failed to unmarshal: %w", err)
	}

	*r = ReturnStmt{}

	for i, e := range temp.Exprs {
		expr, err := unmarshalExpr(e)
		if err != nil {
			return fmt.Errorf("return_stmt: failed to unmarshal %d expr: %w", i, err)
		}
		r.Exprs = append(r.Exprs, expr)
	}

	return nil
}

type BreakStmt struct {
	StmtBase
}

func (b *BreakStmt) MarshalJSON() ([]byte, error) {
	return marshalWithType(b, "break_stmt")
}

type LabelStmt struct {
	StmtBase

	Name string `json:"name"`
}

func (l *LabelStmt) MarshalJSON() ([]byte, error) {
	return marshalWithType(l, "label_stmt")
}

type GotoStmt struct {
	StmtBase

	Label string `json:"label"`
}

func (g *GotoStmt) MarshalJSON() ([]byte, error) {
	return marshalWithType(g, "goto_stmt")
}

func unmarshalStmt(data []byte) (Stmt, error) {
	if bytes.Equal(data, []byte(`null`)) {
		return nil, nil
	}
	t := gjson.GetBytes(data, DiscriminatorField)
	if !t.Exists() {
		return nil, fmt.Errorf("stmt unmarshal error: type discriminator not found")
	}
	var s Stmt
	switch t.String() {
	case "break_stmt":
		return &BreakStmt{}, nil
	case "goto_stmt":
		s = &GotoStmt{}
	case "label_stmt":
		s = &LabelStmt{}
	case "return_stmt":
		s = &ReturnStmt{}
	case "generic_for_stmt":
		s = &GenericForStmt{}
	case "number_for_stmt":
		s = &NumberForStmt{}
	case "if_stmt":
		s = &IfStmt{}
	case "repeat_stmt":
		s = &RepeatStmt{}
	case "while_stmt":
		s = &WhileStmt{}
	case "do_block_stmt":
		s = &DoBlockStmt{}
	case "func_call_stmt":
		s = &FuncCallStmt{}
	case "local_assign_stmt":
		s = &LocalAssignStmt{}
	case "assign_stmt":
		s = &AssignStmt{}
	case "func_def_stmt":
		s = &FuncDefStmt{}
	default:
		return nil, fmt.Errorf("stmt unmarshal error: unknown %s", t.String())
	}

	if err := json.Unmarshal(data, s); err != nil {
		return nil, fmt.Errorf("stmt unmarshal error: %w", err)
	}

	return s, nil
}
