package repl

import (
	"bufio"
	"fmt"
	"monkey/lexer"
	"monkey/parser"
	"os"
)

const PROMPT = ">> "

func Start() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Fprintf(os.Stdout, "%s", PROMPT)

		if ok := scanner.Scan(); !ok {
			fmt.Println("idk")
			return
		}

		input := scanner.Text()

		l := lexer.New(input)
		p := parser.New(l)
		program := p.ParseProgram()
		p.CheckErrors()

		for _, stmt := range program.Statements {
			fmt.Println(stmt)
		}
	}
}
