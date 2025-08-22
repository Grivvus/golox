
## Tree-walk interpreter of lox programming language written in go

[crafting interpreters book](https://www.craftinginterpreters.com/a-tree-walk-interpreter.html)

## build/run

`sh your_program.sh run filename.lox # if you want to run code`

`sh your_program.sh tokenize filename.lox # if you want to get all tokens`

`sh your_program.sh parse filename.lox # if you want to parse expression`

`sh your_program.sh evaluate filename.lox # if you want to evaluate expression`

## supports
- [x] expressions
- [x] statements
- [x] control flow (branches, loops)
- [x] functions
- [x] classes, methods
- [x] inheritance

- [x] arrays (partly, they're immutable, no helpfull builtins, only declaration and subscription) 


## some lox code

```
class Base {
  method() {
    print "Base.method()";
  }
}

class Parent < Base {
  method() {
    super.method();
  }
}

class Child < Parent {
  method() {
    super.method();
  }
}

var parent = Parent();
parent.method();
var child = Child();
child.method();
```
### output

```
Base.method()
Base.method()
```
