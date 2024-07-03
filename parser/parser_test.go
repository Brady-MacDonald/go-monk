package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"testing"
)

func TestIntLiteral(t *testing.T) {
	input := "5;"

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()

	if len(program.Statements) != 1 {
		t.Fatalf("Program does not have enough statements. Got=%d, Expected=1", len(program.Statements))
	}

	// Cast the slice of Statement interfaces to a specific type
	// The IntLiteral pointer implements the Statement interface (all methods have a pointer receiver, therefor must be cast as a pointer)
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not an *ExpressionStatement. Got=%T", program.Statements[0])
	}

	// Ensure the ExpressionStatement statement contains the expected Expression struct
	intLit, ok := stmt.Expr.(*ast.IntLiteral)
	if !ok {
		t.Fatalf("ExpressionStatement does not contain an *IntLiteral. Got=%T", stmt)
	}

	if intLit.Value != 5 {
		t.Fatalf("IntLiteral.Value does not contain correct value. Got=%d, Expected=5", intLit.Value)
	}

	if intLit.TokenLiteral() != "5" {
		t.Fatalf("IntLiteral.TokenLiteral() does not contain token value. Got=%s, Expected=5", intLit.TokenLiteral())
	}
}

func TestIdentifier(t *testing.T) {
	input := "foobar"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("Program does not have enough statements. Got=%d, Expected=1", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not an ExpressionStatement. Got=%T", program.Statements[0])
	}

	ident, ok := stmt.Expr.(*ast.Identifier)
	if !ok {
		t.Fatalf("ExpressionStatement does not contain an Identifier. Got=%T", stmt)
	}

	if ident.Name != "foobar" {
		t.Fatalf("Identifier.Name does not contain correct value. Got=%s, Expected=foobar", ident.Name)
	}

	if ident.TokenLiteral() != "foobar" {
		t.Fatalf("Identifier.TokenLiteral() does not contain token value. Got=%s, Expected=foobar", ident.TokenLiteral())
	}
}

func TestPrefixExpression(t *testing.T) {
	expected := []struct {
		input    string
		operator string
		value    int64
	}{
		{input: "!5", operator: "!", value: 5},
		{input: "-21", operator: "-", value: 21},
	}

	for idx, tt := range expected {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()

		if len(program.Statements) != 1 {
			t.Fatalf("test[%d]: Program does not have enough statements. Got=%d, Expected=1", idx, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("test[%d]: program.Statements[0] is not an ExpressionStatement. Got=%T", idx, program.Statements[0])
		}

		prefixExpr, ok := stmt.Expr.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("ExpressionStatement does not contain an PrefixExpression. Got=%T", stmt)
		}

		if prefixExpr.Operator != tt.operator {
			t.Fatalf("ExpressionStatement.Operator does not contain correct value. Got=%s, Expected=%s", prefixExpr.Operator, tt.operator)
		}

		testIntLiteral(t, prefixExpr.Right, tt.value)
	}

}

func TestLetStatment(t *testing.T) {
	input := `let x = 5`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()

	checkParserErrors(t, p)
	assertAstLength(t, program, 1)

	letStmt, ok := program.Statements[0].(*ast.LetStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *LetStatement. Got=%T", program.Statements[0])
	}

	testIdentifier(t, letStmt.Name, "x")
}

func TestBooleanExpression(t *testing.T) {
	input := "true"

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()

	checkParserErrors(t, p)
	assertAstLength(t, program, 1)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not an *ExpressionStatement. Got=%T", program.Statements[0])
	}

	boolExp, ok := stmt.Expr.(*ast.BoolLiteral)
	if !ok {
		t.Fatalf("ExpressionStatement does not contain an *Boolean. Got=%T", stmt)
	}

	if boolExp.Value != true {
		t.Fatalf("Boolean.Value does not contain correct value. Got=%v, Expected=5", boolExp.Value)
	}

	if boolExp.TokenLiteral() != "true" {
		t.Fatalf("Boolean.TokenLiteral() does not contain token value. Got=%s, Expected='true'", boolExp.TokenLiteral())
	}
}

func TestIfElseExpression(t *testing.T) {
	input := "if(x < y){x}else{ y }"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	checkParserErrors(t, p)
	assertAstLength(t, program, 1)
}

func TestIfExpression(t *testing.T) {
	input := "if(x < y){x}"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	checkParserErrors(t, p)
	assertAstLength(t, program, 1)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatal()
	}

	ifExpr, ok := stmt.Expr.(*ast.IfExpression)
	if !ok {
		t.Fatal()
	}

	if !testInfixExpression(t, ifExpr.Condition, "<") {
		t.Fatal()
	}

	if len(ifExpr.Consequence.Statements) != 1 {
		t.Fatal()
	}

	exprStmt, ok := ifExpr.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatal()
	}

	if !testIdentifier(t, exprStmt.Expr, "x") {
		t.Fatal()
	}

	if ifExpr.Alternitive != nil {
		t.Fatal()
	}
}

/*** Helpers ***/

func assertAstLength(t *testing.T, program *ast.Program, length int) {
	if len(program.Statements) != 1 {
		t.Fatalf("Program does not have enough statements. Got=%d, Expected=1", len(program.Statements))
	}
}

func checkParserErrors(t *testing.T, p *Parser) {
	if len(p.errors) != 0 {
		for _, err := range p.errors {
			t.Error(err)
		}
	}
}

// Generic helper function for testing an Expression
// Specific test func will be called based on type of Expression
func testExpression(t *testing.T, exp ast.Expression, expected any) bool {
	switch expression := exp.(type) {
	case *ast.BoolLiteral:
		val, _ := expected.(bool)
		return testBoolean(t, *expression, val)
	case *ast.Identifier:
		str, _ := expected.(string)
		return testIdentifier(t, expression, str)
	case *ast.IntLiteral:
		val, _ := expected.(int64)
		return testIntLiteral(t, expression, val)
	default:
		t.Fatalf("Unknown Expression type to test %T", expression)
		return false
	}
}

func testInfixExpression(t *testing.T, exp ast.Expression, operator string) bool {
	return true
}

func testBoolean(t *testing.T, boolExp ast.BoolLiteral, val bool) bool {
	if boolExp.Value != val {
		t.Errorf("Boolean.Value incorrect val, Got=%v, expected=%v", boolExp.Value, val)
		return false
	}

	if boolExp.TokenLiteral() != fmt.Sprintf("%v", val) {
		t.Errorf("Boolean.TokenLiteral() incorrect val, Got=%v, expected=%v", boolExp.TokenLiteral(), val)
		return false
	}

	return true
}

func testIdentifier(t *testing.T, identExp ast.Expression, val string) bool {
	ident, ok := identExp.(*ast.Identifier)
	if !ok {
		t.Errorf("Expression is not an *Identifier, Got=%T", identExp)
		return false
	}

	if ident.Name != val {
		t.Errorf("Identifier.Name does not match, Got=%s expected=%s", ident.Name, val)
		return false
	}

	if ident.TokenLiteral() != val {
		t.Errorf("Identifier.TokenLiteral() does not match, Got=%s expected=%s", ident.TokenLiteral(), val)
		return false
	}

	return true
}

func testIntLiteral(t *testing.T, intExp ast.Expression, val int64) bool {
	intLit, ok := intExp.(*ast.IntLiteral)
	if !ok {
		t.Fatalf("Expression not in IntLiteral, Got=%T", intExp)
		return false
	}

	if intLit.Value != val {
		t.Fatalf("IntLiteral.Value Got=%d, expected=%d", intLit.Value, val)
		return false
	}
	return true
}
