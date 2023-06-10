package main

import (
	"log"
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

	program, parseErr := parser.ParseProgram(lexer)
	if parseErr != nil {
		lexer.Position = parseErr.Position
		log.Fatal(parseErr.ToError(lexer.CurrentLine()))
	}

	compiled, compileErr := compiler.Compile(*program)
	if compileErr != nil {
		lexer.Position = compileErr.Position
		log.Fatal(compileErr.ToError(lexer.CurrentLine()))
	}

	output := compiler.Render(compiled)
	outputBytes := []byte(output)

	os.WriteFile(fOut, outputBytes, 0644)

}
