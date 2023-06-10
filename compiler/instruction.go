package compiler

import (
	"fmt"
	"strings"
)

type Instruction struct {
	Opcode   string
	Args     []string
	IsIndent bool
}

var labelCounter = 0

func genLabel() string {
	labelCounter++
	return fmt.Sprintf("label_%d", labelCounter)
}

func Render(instructions []Instruction) string {
	var output strings.Builder

	for _, instruction := range instructions {
		args := strings.Join(instruction.Args, ", ")

		if instruction.IsIndent {
			output.WriteString("\t")
		}

		output.WriteString(instruction.Opcode)
		output.WriteString(" ")
		output.WriteString(args)
		output.WriteString("\n")
	}

	return output.String()
}

func SECTION(section string) Instruction {
	return Instruction{
		Opcode:   "section",
		Args:     []string{section},
		IsIndent: false,
	}
}

func GLOBAL(symbol string) Instruction {
	return Instruction{
		Opcode:   "global",
		Args:     []string{symbol},
		IsIndent: false,
	}
}

func LABEL(symbol string) Instruction {
	return Instruction{
		Opcode:   symbol + ":",
		Args:     []string{},
		IsIndent: false,
	}
}

func PUSH(address string) Instruction {
	return Instruction{
		Opcode:   "push",
		Args:     []string{address},
		IsIndent: true,
	}
}

func MOV(destination string, source string) Instruction {
	return Instruction{
		Opcode:   "mov",
		Args:     []string{destination, source},
		IsIndent: true,
	}
}

func ADD(destination string, source string) Instruction {
	return Instruction{
		Opcode:   "add",
		Args:     []string{destination, source},
		IsIndent: true,
	}
}

func SUB(destination string, source string) Instruction {
	return Instruction{
		Opcode:   "sub",
		Args:     []string{destination, source},
		IsIndent: true,
	}
}

func CMP(destination string, source string) Instruction {
	return Instruction{
		Opcode:   "cmp",
		Args:     []string{destination, source},
		IsIndent: true,
	}
}

func JL(label string) Instruction {
	return Instruction{
		Opcode:   "jl",
		Args:     []string{label},
		IsIndent: true,
	}
}

func JMP(label string) Instruction {
	return Instruction{
		Opcode:   "jmp",
		Args:     []string{label},
		IsIndent: true,
	}
}

func SYSCALL() Instruction {
	return Instruction{
		Opcode:   "syscall",
		Args:     []string{},
		IsIndent: true,
	}
}
