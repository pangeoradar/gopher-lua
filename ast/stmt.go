package ast

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
	"strings"
)

type StringerIndent interface {
	StringIndent(int) string
}

type Stmt interface {
	PositionHolder
	StringerIndent
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

func (a *AssignStmt) StringIndent(indent int) string {
	var lhs, rhs []string
	for i := 0; i < len(a.Lhs); i++ {
		lhs = append(lhs, a.Lhs[i].String())
		rhs = append(rhs, a.Rhs[i].String())
	}
	return fmt.Sprintf("%s%s = %s\n", strings.Repeat("\t", indent),
		strings.TrimRight(strings.Join(lhs, ", "), ", "),
		strings.TrimRight(strings.Join(rhs, ", "), ", "),
	)
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

func (l *LocalAssignStmt) StringIndent(indent int) string {
	// todo: добавить поддержку локальных функций
	lhs := strings.TrimRight(strings.Join(l.Names, ", "), ", ")
	rhs := ""
	if len(l.Exprs) > 0 {
		var exprs []string
		for _, expr := range l.Exprs {
			exprs = append(exprs, expr.String())
		}
		rhs = fmt.Sprintf(" = %s", strings.TrimRight(strings.Join(exprs, ", "), ", "))
	}
	return fmt.Sprintf("%s%s%s\n", strings.Repeat("\t", indent), lhs, rhs)
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

func (f *FuncCallStmt) StringIndent(indent int) string {
	// todo: с первого взгляда этот тип создан для того чтобы описывать
	// todo: единичный вызов функции на строке. но в парсере логика сложнее.
	// todo: хорошо бы разобраться @a.odintsov
	s := fmt.Sprintf("%s%s\n", strings.Repeat("\t", indent), f.Expr)
	return s
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

func (d *DoBlockStmt) StringIndent(indent int) string {
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("%sdo\n", strings.Repeat("\t", indent)))
	for _, stmt := range d.Stmts {
		builder.WriteString(stmt.StringIndent(indent + 1))
	}
	builder.WriteString(fmt.Sprintf("%send\n", strings.Repeat("\t", indent)))
	return builder.String()
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

func (w *WhileStmt) StringIndent(indent int) string {
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("%swhile %s do\n", strings.Repeat("\t", indent), w.Condition))
	for _, stmt := range w.Stmts {
		builder.WriteString(stmt.StringIndent(indent + 1))
	}
	builder.WriteString(fmt.Sprintf("%send\n", strings.Repeat("\t", indent)))
	return builder.String()
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

func (r *RepeatStmt) StringIndent(indent int) string {
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("%srepeat\n", strings.Repeat("\t", indent)))
	for _, stmt := range r.Stmts {
		builder.WriteString(stmt.StringIndent(indent + 1))
	}
	builder.WriteString(fmt.Sprintf("%suntil %s\n", strings.Repeat("\t", indent), r.Condition))
	return builder.String()
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

func (i *IfStmt) StringIndent(indent int) string {
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("%sif %s then\n", strings.Repeat("\t", indent), i.Condition))
	for _, stmt := range i.Then {
		builder.WriteString(stmt.StringIndent(indent + 1))
	}
	if len(i.Else) > 0 {
		builder.WriteString(fmt.Sprintf("%selse\n", strings.Repeat("\t", indent)))
		for _, stmt := range i.Else {
			builder.WriteString(stmt.StringIndent(indent + 1))
		}
	}
	builder.WriteString(fmt.Sprintf("%send\n", strings.Repeat("\t", indent)))
	return builder.String()
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

func (n *NumberForStmt) StringIndent(indent int) string {
	exprs := []string{n.Init.String(), n.Limit.String()}
	if n.Step != nil {
		exprs = append(exprs, n.Step.String())
	}
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("%sfor %s = %s do\n", strings.Repeat("\t", indent),
		n.Name, strings.TrimRight(strings.Join(exprs, ", "), ", ")))
	for _, stmt := range n.Stmts {
		builder.WriteString(stmt.StringIndent(indent + 1))
	}
	builder.WriteString(fmt.Sprintf("%send\n", strings.Repeat("\t", indent)))
	return builder.String()
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

func (g *GenericForStmt) StringIndent(indent int) string {
	var exprs []string
	for _, expr := range g.Exprs {
		exprs = append(exprs, expr.String())
	}
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("%sfor %s in %s do\n",
		strings.Repeat("\t", indent),
		strings.TrimRight(strings.Join(g.Names, ", "), ", "),
		strings.TrimRight(strings.Join(exprs, ", "), ", "),
	))
	for _, stmt := range g.Stmts {
		builder.WriteString(stmt.StringIndent(indent + 1))
	}
	builder.WriteString(fmt.Sprintf("%send\n", strings.Repeat("\t", indent)))
	return builder.String()
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

func (f *FuncDefStmt) StringIndent(indent int) string {
	return fmt.Sprintf("%sfunction %s%s",
		strings.Repeat("\t", indent),
		f.Name,
		f.Func.StringIndent(indent),
	)
}

func (f *FuncDefStmt) MarshalJSON() ([]byte, error) {
	return marshalWithType(f, "func_def_stmt")
}

type ReturnStmt struct {
	StmtBase

	Exprs []Expr `json:"exprs"`
}

func (r *ReturnStmt) StringIndent(indent int) string {
	var exprs []string
	for _, expr := range r.Exprs {
		exprs = append(exprs, expr.String())
	}
	return fmt.Sprintf("%sreturn %s\n", strings.Repeat("\t", indent),
		strings.TrimRight(strings.Join(exprs, ", "), ", "))
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

func (b *BreakStmt) StringIndent(indent int) string {
	return fmt.Sprintf("%sbreak\n", strings.Repeat("\t", indent))
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

func (l *LabelStmt) StringIndent(indent int) string {
	return fmt.Sprintf("%s::%s::\n", strings.Repeat("\t", indent), l.Name)
}

type GotoStmt struct {
	StmtBase

	Label string `json:"label"`
}

func (g *GotoStmt) StringIndent(indent int) string {
	return fmt.Sprintf("%sgoto %s\n", strings.Repeat("\t", indent), g.Label)
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
