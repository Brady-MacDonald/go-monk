package evaluator

import (
	"monkey/ast"
	"monkey/object"
)

// Single instance values
var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

// Evaluate given ast.Node based on its type
func Eval(astNode ast.Node, env *object.Environment) object.Object {
	switch node := astNode.(type) {

	case *ast.Program:
		return evalProgram(node.Statements, env)

	case *ast.BlockStatement:
		return evalBlockStatment(node.Statements, env)

	case *ast.ReturnStatement:
		return evalReturnStatement(node, env)

	case *ast.LetStatement:
		return evalLetStatement(node, env)

	case *ast.ExpressionStatement:
		return Eval(node.Expr, env)

	case *ast.IntLiteral:
		return &object.Integer{Value: node.Value}

	case *ast.StringLiteral:
		return &object.String{Value: node.Value}

	case *ast.BoolLiteral:
		return getBoolObj(node.Value)

	case *ast.ArrayLiteral:
		return evalArrayLiteral(node, env)

	case *ast.HashLiteral:
		return evalHashLiteral(node, env)

	case *ast.IndexExpression:
		return evalIndexExpression(node, env)

	case *ast.Identifier:
		return evalIdentifier(node, env)

	case *ast.FnLiteral:
		return evalFnLiteral(node, env)

	case *ast.CallExpression:
		return evalCallExpression(node, env)

	case *ast.IfExpression:
		return evalIfExpression(node, env)

	case *ast.PrefixExpression:
		return evalPrefixExpression(node, env)

	case *ast.InfixExpression:
		return evalInfixExpression(node, env)

	default:
		return nil
	}
}

// Evaluate the list of Statements in the Program node returning the last value
func evalProgram(statements []ast.Statement, env *object.Environment) object.Object {
	var obj object.Object

	for _, stmt := range statements {
		obj = Eval(stmt, env)

		// Do not evaluate remaining statements after Return/Error objects
		switch result := obj.(type) {
		case *object.Return:
			return result.Value // Unwrap the Return object, returning its value
		case *object.Error:
			return result
		}
	}

	return obj
}

// Evaluate a BlockStatement's slice of statements
func evalBlockStatment(statements []ast.Statement, env *object.Environment) object.Object {
	var obj object.Object

	for _, stmt := range statements {
		obj = Eval(stmt, env)

		// Check for Return object in BlockStatement
		// This should return from not just block scope but entire program
		if obj != nil && obj.Type() == object.RETURN_OBJ {
			// Do not unwrap, but return the actual ReturnObject.
			// Any nested BlockStatements will now return this obj and return from the top level Program
			return obj
		}
	}

	return obj
}

// Eval each expression provided in slice of ast.Expression
func evalExpressions(expressions []ast.Expression, env *object.Environment) []object.Object {
	var results []object.Object

	for _, exp := range expressions {
		expRes := Eval(exp, env)
		if isError(expRes) {
			return []object.Object{expRes}
		}

		results = append(results, expRes)
	}

	return results
}

// Bind the identifier of LetStatement with the value produces by its expression in the given Environment
func evalLetStatement(stmt *ast.LetStatement, env *object.Environment) object.Object {
	expVal := Eval(stmt.Value, env)
	if isError(expVal) {
		return expVal
	}

	env.Set(stmt.Name.Value, expVal)
	return nil
}

func evalReturnStatement(ret *ast.ReturnStatement, env *object.Environment) object.Object {
	val := Eval(ret.Value, env)
	if isError(val) {
		return val
	}

	return &object.Return{Value: val}
}

// Evaluate the expressions in an ast.ArrayLiteral
func evalArrayLiteral(arrAst *ast.ArrayLiteral, env *object.Environment) object.Object {
	evaledExpr := evalExpressions(arrAst.Elements, env)
	if len(evaledExpr) == 1 && isError(evaledExpr[0]) {
		return evaledExpr[0]
	}

	return &object.Array{Value: evaledExpr}
}

func evalHashLiteral(hashAst *ast.HashLiteral, env *object.Environment) object.Object {
	hash := &object.Hash{
		Pairs: map[object.HashKey]object.HashPair{},
	}

	for key, val := range hashAst.Pairs {
		evalKey := Eval(key, env)
		if isError(evalKey) {
			return evalKey
		}

		hashableKey, ok := evalKey.(object.Hashable)
		if !ok {
			return newError("Key is not HashAble. Got=%s", evalKey.Type())
		}

		evalVal := Eval(val, env)
		if isError(evalVal) {
			return evalVal
		}

		hashKey := hashableKey.HashKey()
		hash.Pairs[hashKey] = object.HashPair{
			Key: evalKey,
			Val: evalVal,
		}
	}

	return hash
}

// Get the element of the Array or Hash specified by the index
func evalIndexExpression(idxExp *ast.IndexExpression, env *object.Environment) object.Object {
	arrObj := Eval(idxExp.Left, env)
	if isError(arrObj) {
		return arrObj
	}

	idxObj := Eval(idxExp.Index, env)
	if isError(idxObj) {
		return idxObj
	}

	switch obj := arrObj.(type) {
	case *object.Array:
		return evalArrayIndex(obj, idxObj, env)

	case *object.Hash:
		return evalHashIndex(obj, idxObj, env)

	default:
		return newError("Only Array/Hash is index-able, Got=%s. Line %d Column %d", arrObj.Type(), idxExp.Token.Position.Line, idxExp.Token.Position.Column)
	}
}

func evalHashIndex(hash *object.Hash, idxObj object.Object, env *object.Environment) object.Object {
	idx, ok := idxObj.(object.Hashable)
	if !ok {
		return newError("Key is not HashAble, Got=%s", idxObj.Type())
	}

	key := idx.HashKey()
	if pair, ok := hash.Pairs[key]; ok {
		return pair.Val
	}

	return NULL

}

func evalArrayIndex(arr *object.Array, idxObj object.Object, env *object.Environment) object.Object {
	idx, ok := idxObj.(*object.Integer)
	if !ok {
		return newError("Index is not an Integer, Got=%s", idxObj.Type())
	}

	if idx.Value < 0 || int(idx.Value) >= len(arr.Value) {
		return NULL
	}

	return arr.Value[idx.Value]
}

// Get the Object bound to the given Identifier
func evalIdentifier(ident *ast.Identifier, env *object.Environment) object.Object {
	val := env.Get(ident.Value)
	if val != nil {
		return val
	}

	//Check for identifier as a builtin function
	if builtin, ok := builtins[ident.Value]; ok {
		return builtin
	}

	return newError("Unknown Identifier %s", ident.Value)
}

func evalCallExpression(callExp *ast.CallExpression, env *object.Environment) object.Object {
	// Identifier or FnLiteral should produce a Function object
	// Builtin function also possible
	fnObj := Eval(callExp.Fn, env)
	if isError(fnObj) {
		return fnObj
	}

	evalFnArgs := evalExpressions(callExp.Args, env)
	if len(evalFnArgs) == 1 && evalFnArgs[0].Type() == object.ERROR_OBJ {
		return evalFnArgs[0]
	}

	return applyFunction(fnObj, evalFnArgs)
}

// Create the new Environment to be use inside function execution.
// Bind args of function call to function parameters in new env
func applyFunction(obj object.Object, args []object.Object) object.Object {

	switch fn := obj.(type) {
	case *object.Function:
		// Create new Environment to be used when evaluating function
		// Environment which the function was declared in (closure) used as enclosing env
		fnEnv := object.NewEnclosingEnvironment(fn.Env)

		if len(fn.Parameters) != len(args) {
			return newError("Call expression does not match number of Function paramters: args=%d, params=%d", len(args), len(fn.Parameters))
		}

		for idx, param := range fn.Parameters {
			fnEnv.Set(param.Value, args[idx])
		}

		evalFn := Eval(fn.Body, fnEnv)
		returnVal, ok := evalFn.(*object.Return)
		if ok {
			return returnVal.Value
		}

		return evalFn

	case *object.Builtin:
		val := fn.Fn(args...)
		return val

	default:
		return newError("Not a function %s", obj.Type())
	}
}

func evalFnLiteral(fn *ast.FnLiteral, env *object.Environment) object.Object {
	fnObj := &object.Function{
		Parameters: fn.Parameters,
		Body:       fn.Body,
		// Lexical env the function was declared in, creating a closure
		// Includes the identifier this fnObj is bound allowing recursive calls
		Env: env,
	}

	return fnObj
}

// Evaluate ast.IfExpression
func evalIfExpression(ie *ast.IfExpression, env *object.Environment) object.Object {
	cond := Eval(ie.Condition, env)
	if isError(cond) {
		return cond
	}

	if isTruthy(cond) {
		return Eval(ie.Consequence, env)
	}

	if !isTruthy(cond) && ie.Alternative != nil {
		return Eval(ie.Alternative, env)
	}

	// Condition is false and no else clause provided
	return NULL
}

// Evaluate the prefix expression (! -)
// If operator is not valid an Error is returned
func evalPrefixExpression(prefix *ast.PrefixExpression, env *object.Environment) object.Object {
	operand := Eval(prefix.Operand, env)
	if isError(operand) {
		return operand
	}

	switch prefix.Operator {
	case "!":
		return evalBangPrefix(operand)
	case "-":
		return evalMinusPrefix(operand)
	default:
		return newError("Unknown prefix operator: %s%s", prefix.Operator, operand.Type())
	}

}

// Evaluate the given infix expression
func evalInfixExpression(infix *ast.InfixExpression, env *object.Environment) object.Object {
	right := Eval(infix.Right, env)
	if isError(right) {
		return right
	}

	left := Eval(infix.Left, env)
	if isError(left) {
		return left
	}

	switch {
	case left.Type() != right.Type():
		return newError("Infix expression type mismatch: %s %s %s", left.Type(), infix.Operator, right.Type())

	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfix(infix.Operator, left, right)

	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return evalStringInfix(infix.Operator, left, right)
	}

	switch infix.Operator {
	case "==":
		return getBoolObj(left == right)
	case "!=":
		return getBoolObj(left != right)

	default:
		return newError("Infix expression type mismatch: %s %s %s", left.Type(), infix.Operator, right.Type())
	}
}
