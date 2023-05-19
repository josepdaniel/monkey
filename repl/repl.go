package repl

import (
	"bufio"
	"fmt"
	"io"
	"monkey/lexer"
	"monkey/token"
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

		for l, tok := l.Next(); tok.Type != token.EOF; l, tok = l.Next() {
			fmt.Printf("%+v\n", tok)
		}
	}
}
