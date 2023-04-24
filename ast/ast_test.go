package ast

import (
	"encoding/json"
	"testing"
)

func TestNumberExprJSON(t *testing.T) {
	var expr Expr
	expr = &NumberExpr{Value: "10"}

	data, err := json.Marshal(expr)
	if err != nil {
		t.Error(err)
	}

	expr1, err := unmarshalExpr(data)
	if err != nil {
		t.Fatal(err)
	}
	if ne, ok := expr1.(*NumberExpr); !ok || ne.Value != "10" {
		t.Log("should be number expr with value equal 10")
		t.Fail()
	}
}
