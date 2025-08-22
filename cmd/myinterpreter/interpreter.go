package main

import (
	"fmt"
	"os"
	"reflect"
)

type Interpreter struct {
	state   *State
	globals *State
	locals  map[Expr]int
	parser  *Parser
}

func NewInterpreter(parser *Parser) *Interpreter {
	i := new(Interpreter)
	i.state = NewState(nil)
	i.addBuiltins()
	i.globals = i.state
	i.locals = make(map[Expr]int, 0)
	i.parser = parser
	return i
}

func (i *Interpreter) addBuiltins() {
	i.state.define("clock", &LoxTime{})
	i.state.define("floor", &Floor{})
	i.state.define("str", &Str{})
	i.state.define("len", &Len{})
	i.state.define("println", &PrintLine{})
}

func (i Interpreter) visitVarExpr(expr *VarExpr) any {
	return i.lookUpVariable(expr.name, expr)
}

func (i Interpreter) visitLiteralExpr(expr *LiteralExpr) any {
	return expr.value
}

func (i Interpreter) visitGroupingExpr(expr *GroupingExpr) any {
	return i.evaluate(expr.expr)
}

func (i Interpreter) visitUnaryExpr(expr *UnaryExpr) any {
	right := i.evaluate(expr.right)

	switch expr.operator.Token {
	case BANG:
		return !booleanCast(right)
	case MINUS:
		switch right.(type) {
		case float64:
			return -(right.(float64))
		default:
			i.loxRuntimePanicBinNumeric()
		}
	}

	return nil
}

func (i Interpreter) visitBinaryExpr(expr *BinaryExpr) any {
	left := i.evaluate(expr.left)
	right := i.evaluate(expr.right)

	switch expr.operator.Token {
	case STAR:
		if reflect.TypeOf(left).Kind() == reflect.Float64 && reflect.TypeOf(right).Kind() == reflect.Float64 {
			return left.(float64) * right.(float64)
		}
		i.loxRuntimePanicBinNumeric()
	case SLASH:
		if reflect.TypeOf(left).Kind() == reflect.Float64 && reflect.TypeOf(right).Kind() == reflect.Float64 {
			return left.(float64) / right.(float64)
		}
		i.loxRuntimePanicBinNumeric()
	case PERCENT:
		if reflect.TypeOf(left).Kind() == reflect.Float64 && reflect.TypeOf(right).Kind() == reflect.Float64 {
			leftIntegral := int64(left.(float64))
			rightIntegral := int64(right.(float64))
			if float64(leftIntegral) > left.(float64) || float64(rightIntegral) > right.(float64) {
				i.error(expr.operator, "Expect integral numbers")
			}
			return float64(leftIntegral % rightIntegral)
		}
		i.loxRuntimePanicBinNumeric()
	case PLUS:
		left_type := reflect.TypeOf(left)
		right_type := reflect.TypeOf(right)
		if left_type.Kind() == reflect.String && right_type.Kind() == reflect.String {
			return left.(string) + right.(string)
		} else if left_type.Kind() == reflect.Float64 && right_type.Kind() == reflect.Float64 {
			return left.(float64) + right.(float64)
		}
		i.error(expr.operator, "Operands must be two numbers or two strings")
	case MINUS:
		if reflect.TypeOf(left).Kind() == reflect.Float64 && reflect.TypeOf(right).Kind() == reflect.Float64 {
			return left.(float64) - right.(float64)
		}
		i.loxRuntimePanicBinNumeric()
	case GREATER:
		if reflect.TypeOf(left).Kind() == reflect.Float64 && reflect.TypeOf(right).Kind() == reflect.Float64 {
			return left.(float64) > right.(float64)
		}
		i.loxRuntimePanicBinNumeric()
	case GREATER_EQUAL:
		if reflect.TypeOf(left).Kind() == reflect.Float64 && reflect.TypeOf(right).Kind() == reflect.Float64 {
			return left.(float64) >= right.(float64)
		}
		i.loxRuntimePanicBinNumeric()
	case LESS:
		if reflect.TypeOf(left).Kind() == reflect.Float64 && reflect.TypeOf(right).Kind() == reflect.Float64 {
			return left.(float64) < right.(float64)
		}
		i.loxRuntimePanicBinNumeric()
	case LESS_EQUAL:
		if reflect.TypeOf(left).Kind() == reflect.Float64 && reflect.TypeOf(right).Kind() == reflect.Float64 {
			return left.(float64) <= right.(float64)
		}
		i.loxRuntimePanicBinNumeric()
	case EQUAL_EQUAL:
		left_type := reflect.TypeOf(left)
		right_type := reflect.TypeOf(right)
		if left_type.Kind() != right_type.Kind() {
			return false
		} else if left_type.Kind() == reflect.String {
			return left.(string) == right.(string)
		} else if left_type.Kind() == reflect.Float64 {
			return left.(float64) == right.(float64)
		} else if left_type.Kind() == reflect.Bool {
			return left.(bool) == right.(bool)
		}

	case BANG_EQUAL:
		left_type := reflect.TypeOf(left)
		right_type := reflect.TypeOf(right)
		if left_type.Kind() != right_type.Kind() {
			return true
		} else if left_type.Kind() == reflect.String {
			return left.(string) != right.(string)
		} else if left_type.Kind() == reflect.Float64 {
			return left.(float64) != right.(float64)
		} else if left_type.Kind() == reflect.Bool {
			return left.(bool) != right.(bool)
		}
	}

	return nil
}

func (i Interpreter) visitLogicalExpr(expr *LogicalExpr) any {
	left := i.evaluate(expr.left)
	if expr.operator.Token == OR {
		if booleanCast(left) == true {
			return left
		}
	} else {
		if booleanCast(left) == false {
			return left
		}
	}
	right := i.evaluate(expr.right)
	return right
}

func (i Interpreter) visitCallExpr(expr *CallExpr) any {
	callee := i.evaluate(expr.callee)
	switch callee.(type) {
	case LoxCallable:
		goto FINE
	default:
		i.error(expr.caleeToken, "Can only call functions and classes")
	}
FINE:
	arguments := make([]any, 0)
	for _, arg := range expr.args {
		arguments = append(arguments, i.evaluate(arg))
	}
	var function LoxCallable
	function = callee.(LoxCallable)
	if function.arity() != len(arguments) {
		i.error(expr.caleeToken, fmt.Sprintf("expected %v arguments but got %v\n", function.arity(), len(arguments)))
	}
	return function.call(i, arguments)
}

func (i Interpreter) visitArrayDeclExpr(expr *ArrayDeclExpr) any {
	eval_elements := make([]any, len(expr.elements))
	for idx, element := range expr.elements {
		eval_elements[idx] = i.evaluate(element)
	}
	return eval_elements
}

func (i Interpreter) visitSubscriptExpr(expr *SubscriptExpr) any {
	array := i.evaluate(expr.object)
	index := i.evaluate(expr.index)
	switch array := array.(type) {
	case []any:
		switch index := index.(type) {
		case float64:
			intIndex := int64(index)
			if float64(intIndex) < index {
				i.error(expr.indexToken, "Expected integral number")
			}
			if intIndex >= int64(len(array)) {
				i.error(expr.indexToken, "Out of range")
			}
			return array[intIndex]
		default:
			i.error(expr.indexToken, "Expect number")
		}
	default:
		i.error(expr.objectToken, "Only arrays can be subscripted")
	}
	panic("unreachable")
}

func (i Interpreter) visitGetExpr(expr *GetExpr) any {
	object := i.evaluate(expr.object)
	switch object.(type) {
	case *LoxInstance:
		return object.(*LoxInstance).Get(expr.name)
	default:
		i.error(expr.name, "Only instance have properties")
	}
	return nil
}

func (i Interpreter) visitSetExpr(expr *SetExpr) any {
	object := i.evaluate(expr.object)
	var exprRes any
	switch object.(type) {
	case *LoxInstance:
		exprRes = i.evaluate(expr.value)
		object.(*LoxInstance).Set(expr.name, exprRes)
	default:
		i.error(expr.name, "Only instance have fields")
	}
	return exprRes
}

func (i Interpreter) visitSuperExpr(expr *SuperExpr) any {
	distance := i.locals[expr]
	superclass := i.state.accessAt(distance, "super").(*LoxClass)
	instance := i.state.accessAt(distance-1, "this").(*LoxInstance)
	method := superclass.findMethod(expr.method.Lexeme)
	if method == nil {
		i.error(expr.keyword, fmt.Sprintf("Undefined property %v'", expr.method.Lexeme))
	}
	return method.bind(instance)
}

func (i Interpreter) visitThisExpr(expr *ThisExpr) any {
	return i.lookUpVariable(expr.keyword, expr)
}

func (i Interpreter) visitAssignExpr(expr *AssignExpr) any {
	value := i.evaluate(expr.value)
	distance, ok := i.locals[expr]
	if ok {
		i.state.assignAt(distance, expr.name.Lexeme, value)
	} else {
		i.globals.assign(expr.name.Lexeme, value)
	}
	return value
}

func (i Interpreter) visitExpressionStmt(stmt *Expression) {
	i.evaluate(stmt.expr)
}

func (i Interpreter) visitPrintStmt(stmt *Print) {
	value := i.evaluate(stmt.expr)
	if value != nil {
		switch v := value.(type) {
		case float64:
			if v == float64(int64(v)) {
				fmt.Println(int64(v))
			} else {
				fmt.Println(v)
			}
		default:
			fmt.Println(value)
		}
	} else {
		fmt.Println("nil")
	}
}

func (i Interpreter) visitVarStmt(stmt *Var) {
	var value any = nil
	if stmt.varValue != nil {
		value = i.evaluate(stmt.varValue)
	}
	i.state.define(stmt.varName.Lexeme, value)
}

func (i Interpreter) visitBlockStmt(stmt *Block) {
	i.executeBlock(stmt, NewState(i.state))
}

func (i Interpreter) visitIfStmt(stmt *If) {
	if booleanCast(i.evaluate(stmt.condition)) {
		i.execute(stmt.thenBranch)
	} else if stmt.elseBranch != nil {
		i.execute(stmt.elseBranch)
	}
}

func (i Interpreter) visitWhileStmt(stmt *While) {
	for booleanCast(i.evaluate(stmt.condition)) == true {
		i.execute(stmt.body)
	}
}

func (i Interpreter) visitClassStmt(stmt *Class) {
	var superclass *LoxClass = nil
	if stmt.superclass != nil {
		superclassEval := i.evaluate(stmt.superclass)
		if _, ok := superclassEval.(*LoxClass); !ok {
			i.error(stmt.name, "Can't inherit not from class")
		}
		superclass = superclassEval.(*LoxClass)
	}
	i.state.define(stmt.name.Lexeme, nil)

	if stmt.superclass != nil {
		i.state = NewState(i.state)
		i.state.define("super", superclass)
	}

	methods := make(map[string]*LoxFunction, 0)
	for _, method := range stmt.methods {
		function := NewLoxFunction(method, i.state, method.name.Lexeme == "init")
		methods[method.name.Lexeme] = function
	}
	cls := NewLoxClass(stmt.name.Lexeme, superclass, methods)
	if superclass != nil {
		i.state = i.state.enclosing
	}
	i.state.assign(stmt.name.Lexeme, cls)
}

func (i Interpreter) visitFunctionStmt(stmt *Function) {
	closure := i.state
	fn := NewLoxFunction(stmt, closure, false)
	i.state.define(stmt.name.Lexeme, fn)
}

func (i Interpreter) visitReturnStmt(stmt *Return) {
	var result any = nil
	if stmt.value != nil {
		result = i.evaluate(stmt.value)
	}
	if result == nil {
		result = "nil"
	}
	panic(result)
}

func (i Interpreter) executeBlock(block *Block, state *State) {
	prevState := i.state
	i.state = state
	defer func() { i.state = prevState }()
	for _, stmt := range block.stmts {
		i.execute(stmt)
	}
}

func (i Interpreter) evaluate(expr Expr) any {
	return expr.accept(i)
}

func (i Interpreter) execute(stmt Stmt) {
	stmt.accept(i)
}

func (i Interpreter) resolve(expr Expr, depth int) {
	i.locals[expr] = depth
}

func (i Interpreter) lookUpVariable(name Token, expr Expr) any {
	distance, ok := i.locals[expr]
	if ok {
		return i.state.accessAt(distance, name.Lexeme)
	} else {
		return i.globals.access(name.Lexeme)
	}
}

func booleanCast(expr any) bool {
	if expr == nil {
		return false
	}
	switch expr.(type) {
	case bool:
		return expr.(bool)
	default:
		return true
	}
}

func (i Interpreter) loxRuntimePanicBinNumeric() {
	fmt.Fprintln(os.Stderr, "Operands must be a numbers")
	os.Exit(70)
}

func (i Interpreter) error(token Token, msg string) {
	fmt.Fprintf(os.Stderr, "[line %v] at '%v' %v.\n", token.Line, token.Lexeme, msg)
	os.Exit(70)
}
