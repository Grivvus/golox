package main

import (

)

type LoxCallable interface {
    arity() int
    call(i Interpreter, args []any) any
}
