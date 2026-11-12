package interpreter

import (
	"fmt"
	"golox/internal/parser"
	"golox/internal/scanner"
	"reflect"
)

type Interpreter struct {
	state   *State
	globals *State
	locals  map[parser.Expr]int
	parser  *parser.Parser
}

func New(par *parser.Parser) *Interpreter {
	i := new(Interpreter)
	i.state = NewState(nil)
	i.addBuiltins()
	i.globals = i.state
	i.locals = make(map[parser.Expr]int, 0)
	i.parser = par
	return i
}

func (i *Interpreter) addBuiltins() {
	i.state.define("clock", &LoxTime{})
	i.state.define("floor", &Floor{})
	i.state.define("str", &Str{})
	i.state.define("len", &Len{})
	i.state.define("println", &PrintLine{})
}

func (i Interpreter) VisitVarExpr(expr *parser.VarExpr) any {
	return i.lookUpVariable(expr.Name, expr)
}

func (i Interpreter) VisitLiteralExpr(expr *parser.LiteralExpr) any {
	return expr.Value
}

func (i Interpreter) VisitGroupingExpr(expr *parser.GroupingExpr) any {
	return i.Evaluate(expr.Expr)
}

func (i Interpreter) VisitUnaryExpr(expr *parser.UnaryExpr) any {
	right := i.Evaluate(expr.Right)

	switch expr.Operator.Token {
	case scanner.BANG:
		return !booleanCast(right)
	case scanner.MINUS:
		switch right := right.(type) {
		case float64:
			return -right
		default:
			i.loxRuntimePanicBinNumeric()
		}
	}

	return nil
}

func (i Interpreter) VisitBinaryExpr(expr *parser.BinaryExpr) any {
	left := i.Evaluate(expr.Left)
	right := i.Evaluate(expr.Right)

	switch expr.Operator.Token {
	case scanner.STAR:
		if reflect.TypeOf(left).Kind() == reflect.Float64 && reflect.TypeOf(right).Kind() == reflect.Float64 {
			return left.(float64) * right.(float64)
		}
		i.loxRuntimePanicBinNumeric()
	case scanner.SLASH:
		if reflect.TypeOf(left).Kind() == reflect.Float64 && reflect.TypeOf(right).Kind() == reflect.Float64 {
			return left.(float64) / right.(float64)
		}
		i.loxRuntimePanicBinNumeric()
	case scanner.PERCENT:
		if reflect.TypeOf(left).Kind() == reflect.Float64 && reflect.TypeOf(right).Kind() == reflect.Float64 {
			leftIntegral := int64(left.(float64))
			rightIntegral := int64(right.(float64))
			if float64(leftIntegral) > left.(float64) || float64(rightIntegral) > right.(float64) {
				i.error(expr.Operator, "Expect integral numbers")
			}
			return float64(leftIntegral % rightIntegral)
		}
		i.loxRuntimePanicBinNumeric()
	case scanner.PLUS:
		left_type := reflect.TypeOf(left)
		right_type := reflect.TypeOf(right)
		if left_type.Kind() == reflect.String && right_type.Kind() == reflect.String {
			return left.(string) + right.(string)
		} else if left_type.Kind() == reflect.Float64 && right_type.Kind() == reflect.Float64 {
			return left.(float64) + right.(float64)
		}
		i.error(expr.Operator, "Operands must be two numbers or two strings")
	case scanner.MINUS:
		if reflect.TypeOf(left).Kind() == reflect.Float64 && reflect.TypeOf(right).Kind() == reflect.Float64 {
			return left.(float64) - right.(float64)
		}
		i.loxRuntimePanicBinNumeric()
	case scanner.GREATER:
		if reflect.TypeOf(left).Kind() == reflect.Float64 && reflect.TypeOf(right).Kind() == reflect.Float64 {
			return left.(float64) > right.(float64)
		}
		i.loxRuntimePanicBinNumeric()
	case scanner.GREATER_EQUAL:
		if reflect.TypeOf(left).Kind() == reflect.Float64 && reflect.TypeOf(right).Kind() == reflect.Float64 {
			return left.(float64) >= right.(float64)
		}
		i.loxRuntimePanicBinNumeric()
	case scanner.LESS:
		if reflect.TypeOf(left).Kind() == reflect.Float64 && reflect.TypeOf(right).Kind() == reflect.Float64 {
			return left.(float64) < right.(float64)
		}
		i.loxRuntimePanicBinNumeric()
	case scanner.LESS_EQUAL:
		if reflect.TypeOf(left).Kind() == reflect.Float64 && reflect.TypeOf(right).Kind() == reflect.Float64 {
			return left.(float64) <= right.(float64)
		}
		i.loxRuntimePanicBinNumeric()
	case scanner.EQUAL_EQUAL:
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

	case scanner.BANG_EQUAL:
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

func (i Interpreter) VisitLogicalExpr(expr *parser.LogicalExpr) any {
	left := i.Evaluate(expr.Left)
	if expr.Operator.Token == scanner.OR {
		if booleanCast(left) == true {
			return left
		}
	} else {
		if booleanCast(left) == false {
			return left
		}
	}
	right := i.Evaluate(expr.Right)
	return right
}

func (i Interpreter) VisitCallExpr(expr *parser.CallExpr) any {
	callee := i.Evaluate(expr.Callee)
	switch callee.(type) {
	case LoxCallable:
		goto FINE
	default:
		i.error(expr.CaleeToken, "Can only call functions and classes")
	}
FINE:
	arguments := make([]any, 0)
	for _, arg := range expr.Args {
		arguments = append(arguments, i.Evaluate(arg))
	}
	var function LoxCallable
	function = callee.(LoxCallable)
	if function.arity() != len(arguments) {
		i.error(expr.CaleeToken, fmt.Sprintf("expected %v arguments but got %v\n", function.arity(), len(arguments)))
	}
	return function.call(i, arguments)
}

func (i Interpreter) VisitArrayDeclExpr(expr *parser.ArrayDeclExpr) any {
	eval_elements := make([]any, len(expr.Elements))
	for idx, element := range expr.Elements {
		eval_elements[idx] = i.Evaluate(element)
	}
	return eval_elements
}

func (i Interpreter) VisitSubscriptExpr(expr *parser.SubscriptExpr) any {
	array := i.Evaluate(expr.Object)
	index := i.Evaluate(expr.Index)
	switch array := array.(type) {
	case []any:
		switch index := index.(type) {
		case float64:
			intIndex := int64(index)
			if float64(intIndex) < index {
				i.error(expr.IndexToken, "Expected integral number")
			}
			if intIndex >= int64(len(array)) {
				i.error(expr.IndexToken, "Out of range")
			}
			return array[intIndex]
		default:
			i.error(expr.IndexToken, "Expect number")
		}
	default:
		i.error(expr.ObjectToken, "Only arrays can be subscripted")
	}
	panic("unreachable")
}

func (i Interpreter) VisitGetExpr(expr *parser.GetExpr) any {
	object := i.Evaluate(expr.Object)
	switch object.(type) {
	case *LoxInstance:
		return object.(*LoxInstance).Get(expr.Name)
	default:
		i.error(expr.Name, "Only instance have properties")
	}
	return nil
}

func (i Interpreter) VisitSetExpr(expr *parser.SetExpr) any {
	object := i.Evaluate(expr.Object)
	var exprRes any
	switch object.(type) {
	case *LoxInstance:
		exprRes = i.Evaluate(expr.Value)
		object.(*LoxInstance).Set(expr.Name, exprRes)
	default:
		i.error(expr.Name, "Only instance have fields")
	}
	return exprRes
}

func (i Interpreter) VisitSuperExpr(expr *parser.SuperExpr) any {
	distance := i.locals[expr]
	superclass := i.state.accessAt(distance, "super").(*LoxClass)
	instance := i.state.accessAt(distance-1, "this").(*LoxInstance)
	method := superclass.findMethod(expr.Method.Lexeme)
	if method == nil {
		i.error(expr.Keyword, fmt.Sprintf("Undefined property %v'", expr.Method.Lexeme))
	}
	return method.bind(instance)
}

func (i Interpreter) VisitThisExpr(expr *parser.ThisExpr) any {
	return i.lookUpVariable(expr.Keyword, expr)
}

func (i Interpreter) VisitAssignExpr(expr *parser.AssignExpr) any {
	value := i.Evaluate(expr.Value)
	distance, ok := i.locals[expr]
	if ok {
		i.state.assignAt(distance, expr.Name.Lexeme, value)
	} else {
		i.globals.assign(expr.Name.Lexeme, value)
	}
	return value
}

func (i Interpreter) VisitExpressionStmt(stmt *parser.Expression) {
	i.Evaluate(stmt.Expr)
}

func (i Interpreter) VisitPrintStmt(stmt *parser.Print) {
	value := i.Evaluate(stmt.Expr)
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

func (i Interpreter) VisitVarStmt(stmt *parser.Var) {
	var value any = nil
	if stmt.VarValue != nil {
		value = i.Evaluate(stmt.VarValue)
	}
	i.state.define(stmt.VarName.Lexeme, value)
}

func (i Interpreter) VisitBlockStmt(stmt *parser.Block) {
	i.executeBlock(stmt, NewState(i.state))
}

func (i Interpreter) VisitIfStmt(stmt *parser.If) {
	if booleanCast(i.Evaluate(stmt.Condition)) {
		i.Execute(stmt.ThenBranch)
	} else if stmt.ElseBranch != nil {
		i.Execute(stmt.ElseBranch)
	}
}

func (i Interpreter) VisitWhileStmt(stmt *parser.While) {
	for booleanCast(i.Evaluate(stmt.Condition)) == true {
		i.Execute(stmt.Body)
	}
}

func (i Interpreter) VisitClassStmt(stmt *parser.Class) {
	var superclass *LoxClass = nil
	if stmt.Superclass != nil {
		superclassEval := i.Evaluate(stmt.Superclass)
		if _, ok := superclassEval.(*LoxClass); !ok {
			i.error(stmt.Name, "Can't inherit not from class")
		}
		superclass = superclassEval.(*LoxClass)
	}
	i.state.define(stmt.Name.Lexeme, nil)

	if stmt.Superclass != nil {
		i.state = NewState(i.state)
		i.state.define("super", superclass)
	}

	methods := make(map[string]*LoxFunction, 0)
	for _, method := range stmt.Methods {
		function := NewLoxFunction(method, i.state, method.Name.Lexeme == "init")
		methods[method.Name.Lexeme] = function
	}
	cls := NewLoxClass(stmt.Name.Lexeme, superclass, methods)
	if superclass != nil {
		i.state = i.state.enclosing
	}
	i.state.assign(stmt.Name.Lexeme, cls)
}

func (i Interpreter) VisitFunctionStmt(stmt *parser.Function) {
	closure := i.state
	fn := NewLoxFunction(stmt, closure, false)
	i.state.define(stmt.Name.Lexeme, fn)
}

func (i Interpreter) VisitReturnStmt(stmt *parser.Return) {
	var result any = nil
	if stmt.Value != nil {
		result = i.Evaluate(stmt.Value)
	}
	if result == nil {
		result = "nil"
	}
	panic(result)
}

func (i Interpreter) executeBlock(block *parser.Block, state *State) {
	prevState := i.state
	i.state = state
	defer func() { i.state = prevState }()
	for _, stmt := range block.Stmts {
		i.Execute(stmt)
	}
}

func (i Interpreter) Evaluate(expr parser.Expr) any {
	return expr.Accept(i)
}

func (i Interpreter) Execute(stmt parser.Stmt) {
	stmt.Accept(i)
}

func (i Interpreter) resolve(expr parser.Expr, depth int) {
	i.locals[expr] = depth
}

func (i Interpreter) lookUpVariable(name scanner.Token, expr parser.Expr) any {
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
	panic(fmt.Sprintf("Operands must be a numbers"))
}

func (i Interpreter) error(token scanner.Token, msg string) {
	panic(fmt.Sprintf("[line %v] at '%v' %v.\n", token.Line, token.Lexeme, msg))
}
