package domain

import (
	"fmt"
)

type Displayer interface {
	Display()
}

type CommandLineDisplay struct {
	execProgram string
	expression  ParsedExpression
}

func NewCommandLineDisplay(execProgram string, expression ParsedExpression) CommandLineDisplay {
	return CommandLineDisplay{
		execProgram: execProgram,
		expression:  expression,
	}
}

func (c CommandLineDisplay) Display() {

	for _, v := range c.expression.Fields {

		fmt.Printf("%-14s", v.name)

		// print elements in slice
		for _, n := range v.Values {
			fmt.Printf("%d ", n)
		}
		fmt.Println()
	}
	fmt.Printf("%-14s%s\n", "command", c.execProgram)

}
