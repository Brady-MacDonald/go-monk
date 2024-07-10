package evaluator

import (
	"fmt"
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
	"testing"
)

// Test for expressions which evaluate to Booleans
// BooleanLiterals and InfixExpression with bool left/right operands
func TestEvalBoolean(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"false != true", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
	}

	for _, tt := range tests {
		eval := testEval(tt.input)
		testBoolObject(t, eval, tt.expected)
	}
}

func TestEvalBangPrefix(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
	}

	for _, tt := range tests {
		eval := testEval(tt.input)
		testBoolObject(t, eval, tt.expected)
	}
}

// Test for expressions which produce Integer values
// Expression types: IntegerLiterals, PrefixExpression, InfoxExpression with all ints
func TestEvalInteger(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntObject(t, evaluated, tt.expected)
	}
}

// func TestErrorHandling(t *testing.T) {
//
// }
//
// func TestLetStatment(t *testing.T) {
//
// }

func TestFunctionObject(t *testing.T) {
	input := "fn(x) { x + 2; };"

	evaluated := testEval(input)
	fn, ok := evaluated.(*object.Function)
	if !ok {
		t.Fatalf("object is not Function. got=%T (%+v)", evaluated, evaluated)
	}

	if len(fn.Parameters) != 1 {
		t.Fatalf("function has wrong parameters. Parameters=%+v", fn.Parameters)
	}
	if fn.Parameters[0].String() != "x" {
		t.Fatalf("parameter is not 'x'. got=%q", fn.Parameters[0])
	}

	expectedBody := "(x + 2)"

	if fn.Body.String() != expectedBody {
		t.Fatalf("body is not %q. got=%q", expectedBody, fn.Body.String())
	}
}

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let identity = fn(x) { x; }; identity(5);", 5},
		{"let identity = fn(x) { return x; }; identity(5);", 5},
		{"let double = fn(x) { x * 2; }; double(5);", 10},
		{"let add = fn(x, y) { x + y; }; add(5, 5);", 10},
		{"let add = fn(x, y) { x + y; }; add(5 + 5, add(5, 5));", 20},
		{"fn(x) { x; }(5)", 5},
	}

	for _, tt := range tests {
		fmt.Println(tt)
		obj := testEval(tt.input)
		err, ok := obj.(*object.Error)
		if ok {
			t.Errorf("Error: %s", err.Message)
			continue
		}

		testIntObject(t, obj, tt.expected)
	}
}

func TestConditionals(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if (true) { 10 }", 10},
		{"if (false) { 10 }", nil},
		{"if (1) { 10 }", 10},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 > 2) { 10 }", nil},
		{"if (1 > 2) { 10 } else { 20 }", 20},
		{"if (1 < 2) { 10 } else { 10 }", 10},
	}

	for _, tt := range tests {
		eval := testEval(tt.input)
		exp, ok := tt.expected.(int)
		if ok {
			testIntObject(t, eval, int64(exp))
		} else {
			testNullObject(t, eval)
		}

	}
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"return 10;", 10},
		{"return 10; 9;", 10},
		{"return 2 * 5; 9;", 10},
		{"9; return 2 * 5; 9;", 10},
		{`
			if (10 > 1) {
				if (10 > 1) {
					return 10;
				}

				return 1;
			}
		`, 10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntObject(t, evaluated, tt.expected)
	}
}

/*** Helpers ***/

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	env := object.NewEnvironment()

	return Eval(program, env)
}

func testNullObject(t *testing.T, obj object.Object) bool {
	_, ok := obj.(*object.Null)
	if !ok {
		t.Errorf("Object is not a Null. Got=%t", obj)
		return false
	}

	return true
}

func testBoolObject(t *testing.T, obj object.Object, expected bool) bool {
	boolObj, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("Object is not a Boolean. Got=%t", obj)
		return false
	}

	if boolObj.Value != expected {
		t.Errorf("Boolean Object incorrect value, Got=%v, expected=%v", boolObj.Value, expected)
		return false
	}

	return true
}

func testIntObject(t *testing.T, obj object.Object, val int64) bool {
	intObj, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("Object id not an Integer. Got=%T", obj)
		return false
	}

	if intObj.Value != val {
		t.Errorf("Integer object value not match. Got=%d, expected=%d", intObj.Value, val)
		return false
	}

	return true
}
