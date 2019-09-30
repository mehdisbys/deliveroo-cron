package domain

import (
	"errors"
	"os"
)

type InputFetcher interface {
	GetData() string
}

type CommandLineInput struct {
	Parser           CronParser
	programName      string
	programToExecute string
}

func (c *CommandLineInput) GetProgramToExec() string {
	return c.programToExecute
}

func (c *CommandLineInput) GetDataFromCommandLine() ([]string, error) {
	return c.GetData(os.Args)
}

func (c *CommandLineInput) GetData(argsWithProg []string) ([]string, error) {

	if len(argsWithProg) < (CronElementsSize + 2) {
		return nil, errors.New("invalid input")
	}

	c.programName = argsWithProg[0]

	l := len(argsWithProg)

	c.programToExecute = argsWithProg[l-1]

	cronArgs := make([]string, 0)

	start := l - 1 - CronElementsSize
	end := l - 2

	for i := start; i <= end; i++ {
		cronArgs = append(cronArgs, argsWithProg[i])
	}

	return cronArgs, nil
}
