package repl

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"monkey/lexer"
	"monkey/parser"
)

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	for {
		fmt.Printf(">> ")
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(&line)

		_, statement, err := parser.ParseStatement(l)
		if err != nil {
			fmt.Println(err)
		} else {
			b, _ := json.Marshal(statement)
			fmt.Println(string(b))
		}

	}
}
