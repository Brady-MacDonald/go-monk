package repl

import (
	"bufio"
	"fmt"
	"monkey/evaluator"
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
	"os"
)

const PROMPT = ">>"

func File(name string) {
	content, err := os.ReadFile(name)
	if err != nil {
		panic(err)
	}

	l := lexer.New(string(content))
	p := parser.New(l)
	program := p.ParseProgram()
	if p.ParserErrors() {
		return
	}

	env := object.NewEnvironment()
	val := evaluator.Eval(program, env)
	if val != nil {
		fmt.Println(val.Inspect())
	} else {
		fmt.Println("Evaluated to nil")
	}
}

func Start() {
	scanner := bufio.NewScanner(os.Stdin)
	env := object.NewEnvironment()

	for {
		fmt.Printf("%s ", PROMPT)
		if ok := scanner.Scan(); !ok {
			return
		}

		input := scanner.Text()

		l := lexer.New(input)
		p := parser.New(l)

		program := p.ParseProgram()
		p.ParserErrors()

		eval := evaluator.Eval(program, env)
		if eval != nil {
			// fmt.Println(eval.Type())
			fmt.Println(eval.Inspect())
		}
	}
}
