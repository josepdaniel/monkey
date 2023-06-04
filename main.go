package main

import (
	"monkey/compiler"
	"monkey/lexer"
	"monkey/parser"
	"os"
)

func main() {
	fIn := os.Args[1]
	fOut := os.Args[2]

	bytes, err := os.ReadFile(fIn)

	if err != nil {
		panic(err)
	}

	input := string(bytes)
	lexer := lexer.New(&input)

	program, err := parser.ParseProgram(lexer)
	if err != nil {
		panic(err)
	}

	compiled, err := compiler.Compile(*program)
	if err != nil {
		panic(err)
	}

	output := compiler.Render(compiled)
	outputBytes := []byte(output)

	os.WriteFile(fOut, outputBytes, 0644)

}
