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
	i.state.define("clock", newLoxTime())
	i.globals = i.state
	i.locals = make(map[Expr]int, 0)
	i.parser = parser
	return i
}

func (i Interpreter) interpret(statements []Stmt) {
	for _, stmt := range statements {
		i.execute(stmt)
	}
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
			fmt.Fprintln(os.Stderr, "Operand must be a number")
			os.Exit(70)
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
		loxRuntimePanicBinNumeric()
	case SLASH:
		if reflect.TypeOf(left).Kind() == reflect.Float64 && reflect.TypeOf(right).Kind() == reflect.Float64 {
			return left.(float64) / right.(float64)
		}
		loxRuntimePanicBinNumeric()
	case PLUS:
		left_type := reflect.TypeOf(left)
		right_type := reflect.TypeOf(right)
		if left_type.Kind() == reflect.String && right_type.Kind() == reflect.String {
			return left.(string) + right.(string)
		} else if left_type.Kind() == reflect.Float64 && right_type.Kind() == reflect.Float64 {
			return left.(float64) + right.(float64)
		}
		fmt.Fprintln(os.Stderr, "Operands must be two numbers or two strings")
		os.Exit(70)
	case MINUS:
		if reflect.TypeOf(left).Kind() == reflect.Float64 && reflect.TypeOf(right).Kind() == reflect.Float64 {
			return left.(float64) - right.(float64)
		}
		loxRuntimePanicBinNumeric()
	case GREATER:
		if reflect.TypeOf(left).Kind() == reflect.Float64 && reflect.TypeOf(right).Kind() == reflect.Float64 {
			return left.(float64) > right.(float64)
		}
		loxRuntimePanicBinNumeric()
	case GREATER_EQUAL:
		if reflect.TypeOf(left).Kind() == reflect.Float64 && reflect.TypeOf(right).Kind() == reflect.Float64 {
			return left.(float64) >= right.(float64)
		}
		loxRuntimePanicBinNumeric()
	case LESS:
		if reflect.TypeOf(left).Kind() == reflect.Float64 && reflect.TypeOf(right).Kind() == reflect.Float64 {
			return left.(float64) < right.(float64)
		}
		loxRuntimePanicBinNumeric()
	case LESS_EQUAL:
		if reflect.TypeOf(left).Kind() == reflect.Float64 && reflect.TypeOf(right).Kind() == reflect.Float64 {
			return left.(float64) <= right.(float64)
		}
		loxRuntimePanicBinNumeric()
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

func (i Interpreter) visitDefinedVar(expr *VarExpr) any {
	return i.state.access(expr.name.Lexeme)
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
		fmt.Fprintf(os.Stderr, "Can only call functions and classes.\n[line %v]", i.parser.getPrev().line)
		os.Exit(70)
	}
FINE:
	arguments := make([]any, 0)
	for _, arg := range expr.args {
		arguments = append(arguments, i.evaluate(arg))
	}
	var function LoxCallable
	function = callee.(LoxCallable)
	if function.arity() != len(arguments) {
		fmt.Fprintf(os.Stderr, "expected %v arguments but got %v\n", function.arity(), len(arguments))
		os.Exit(70)
	}
	return function.call(i, arguments)
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
		// fmt.Printf("Lookup: %v at distance %v\n", name.Lexeme, distance)
		return i.state.accessAt(distance, name.Lexeme)
	} else {
		// fmt.Printf("Lookup global: %v\n", name.Lexeme)
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

func loxRuntimePanicBinNumeric() {
	fmt.Fprintln(os.Stderr, "Operands must be a numbers")
	os.Exit(70)
}
