package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"testing"
)

func TestIntLiteral(t *testing.T) {
	input := "5;"
	program := createParseProgram(t, input)

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
	input := "foobar;"
	program := createParseProgram(t, input)

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

	if ident.Value != "foobar" {
		t.Fatalf("Identifier.Name does not contain correct value. Got=%s, Expected=foobar", ident.Value)
	}

	if ident.TokenLiteral() != "foobar" {
		t.Fatalf("Identifier.TokenLiteral() does not contain token value. Got=%s, Expected=foobar", ident.TokenLiteral())
	}
}

func TestPrefixExpression(t *testing.T) {
	expected := []struct {
		input    string
		operator string
		value    interface{}
	}{
		{"!5;", "!", 5},
		{"-15", "-", 15},
		{"!true", "!", true},
		{"!false", "!", false},
	}

	for idx, tt := range expected {
		program := createParseProgram(t, tt.input)

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

		testLiteralExpression(t, prefixExpr.Operand, tt.value)
	}
}

func TestInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
		{"true == true", true, "==", true},
		{"true != false", true, "!=", false},
		{"false == false", false, "==", false},
	}

	for _, tt := range infixTests {
		program := createParseProgram(t, tt.input)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d\n", 1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not %T. got=%T", &ast.ExpressionStatement{}, program.Statements[0])
		}

		if !testInfixExpression(t, stmt.Expr, tt.leftValue, tt.operator, tt.rightValue) {
			return
		}
	}
}

func TestLetStatement(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      interface{}
	}{
		{"let x = 5;", "x", 5},
		{"let y = true;", "y", true},
		{"let foobar = y;", "foobar", "y"},
	}

	for _, tt := range tests {
		program := createParseProgram(t, tt.input)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statements. got=%d", len(program.Statements))
		}

		stmt := program.Statements[0]
		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}

		val := stmt.(*ast.LetStatement).Value
		if !testLiteralExpression(t, val, tt.expectedValue) {
			return
		}
	}
}

func TestBooleanExpression(t *testing.T) {
	input := "true"
	program := createParseProgram(t, input)

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

func TestIfExpression(t *testing.T) {
	input := `if (x < y) { x }`
	program := createParseProgram(t, input)

	if len(program.Statements) != 1 {
		t.Fatalf("program.body does not contain %d statements. got=%d\n", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not %T. got=%T", &ast.ExpressionStatement{}, program.Statements[0])
	}

	exp, ok := stmt.Expr.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not %T. got=%T", &ast.IfExpression{}, stmt.Expr)
	}

	if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}

	if len(exp.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statements. got=%d\n", len(exp.Consequence.Statements))
	}

	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not %T. got=%T", &ast.ExpressionStatement{}, exp.Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Expr, "x") {
		return
	}

	if exp.Alternative != nil {
		t.Errorf("exp.Alternative.Statement was not nil. got=%+v", exp.Alternative)
	}
}

func TestCallExpressionParsing(t *testing.T) {
	input := "add(1, 2 * 3, 4 + 5);"
	program := createParseProgram(t, input)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt is not %T. got=%T", &ast.ExpressionStatement{}, program.Statements[0])
	}

	exp, ok := stmt.Expr.(*ast.CallExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not %T. got=%T", &ast.CallExpression{}, stmt.Expr)
	}

	if !testIdentifier(t, exp.Fn, "add") {
		return
	}

	if len(exp.Args) != 3 {
		t.Fatalf("wrong length of arguments. got=%d", len(exp.Args))
	}

	testLiteralExpression(t, exp.Args[0], 1)
	testInfixExpression(t, exp.Args[1], 2, "*", 3)
	testInfixExpression(t, exp.Args[2], 4, "+", 5)
}

/*** Helpers ***/

func createParseProgram(t *testing.T, input string) *ast.Program {
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	return program
}

func assertAstLength(t *testing.T, program *ast.Program, length int) {
	if len(program.Statements) != 1 {
		t.Fatalf("Program does not have enough statements. Got=%d, Expected=%d", len(program.Statements), length)
	}
}

func checkParserErrors(t *testing.T, p *Parser) {
	if len(p.errors) != 0 {
		for _, err := range p.errors {
			t.Error(err)
		}
	}
}

func testLiteralExpression(t *testing.T, exp ast.Expression, expected interface{}) bool {
	switch v := expected.(type) {
	case int:
		return testIntLiteral(t, exp, int64(v))
	case int64:
		return testIntLiteral(t, exp, v)
	case string:
		return testIdentifier(t, exp, v)
	case bool:
		return testBoolLiteral(t, exp, v)
	}
	t.Errorf("type of exp not handled. got=%T", exp)
	return false
}

func testLetStatement(t *testing.T, stmt ast.Statement, name string) bool {
	if stmt.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral not 'let'. got=%q", stmt.TokenLiteral())
		return false
	}

	letStmt, ok := stmt.(*ast.LetStatement)
	if !ok {
		t.Errorf("Statement not %T. got=%T", &ast.LetStatement{}, stmt)
		return false
	}

	if letStmt.Name.Value != name {
		t.Errorf("letStmt.Name.Value not '%s'. got=%s", name, letStmt.Name.Value)
		return false
	}

	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("s.Name not '%s'. got=%s", name, letStmt.Name)
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

	if ident.Value != val {
		t.Errorf("Identifier.Name does not match, Got=%s expected=%s", ident.Value, val)
		return false
	}

	if ident.TokenLiteral() != val {
		t.Errorf("Identifier.TokenLiteral() does not match, Got=%s expected=%s", ident.TokenLiteral(), val)
		return false
	}

	return true
}

func testBoolLiteral(t *testing.T, expr ast.Expression, val bool) bool {
	boolLit, ok := expr.(*ast.BoolLiteral)
	if !ok {
		t.Errorf("Expression is not a BoolLiteral, Got=%t", expr)
		return false
	}

	if boolLit.Value != val {
		t.Errorf("Boolean.Value incorrect val, Got=%v, expected=%v", boolLit.Value, val)
		return false
	}

	if boolLit.TokenLiteral() != fmt.Sprintf("%v", val) {
		t.Errorf("Boolean.TokenLiteral() incorrect val, Got=%v, expected=%v", expr.TokenLiteral(), val)
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

func testInfixExpression(t *testing.T, exp ast.Expression, left interface{}, operator string, right interface{}) bool {
	infixExp, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("exp is not %T. got=%T('%s')", &ast.InfixExpression{}, exp, exp)
		return false
	}

	if !testLiteralExpression(t, infixExp.Left, left) {
		return false
	}

	if infixExp.Operator != operator {
		t.Errorf("exp.Operator is not '%s'. got=%q", operator, infixExp.Operator)
		return false
	}

	if !testLiteralExpression(t, infixExp.Right, right) {
		return false
	}

	return true
}
