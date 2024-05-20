package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/dominicgaliano/interpreter-demo/evaluator"
	"github.com/dominicgaliano/interpreter-demo/lexer"
	"github.com/dominicgaliano/interpreter-demo/parser"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Print(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		evaluated := evaluator.Eval(program)
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect() + "\n")
		}
	}
}

func printParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, " parser errors:\n")
	for _, error := range errors {
		io.WriteString(out, "\t"+error+"\n")
	}
}
