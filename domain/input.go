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
	cronExpression   string
}

const (
	cronElementsSize = 5
)

func (c *CommandLineInput) GetDataFromCommandLine() (string, error) {
	return c.GetData(os.Args)
}

func (c *CommandLineInput) GetData(argsWithProg []string) (string, error) {

	if len(argsWithProg) < (cronElementsSize + 2) {
		return "", errors.New("invalid command")
	}

	c.programName = argsWithProg[0]

	l := len(argsWithProg)

	c.programToExecute = argsWithProg[l-1]

	cronArgs := make([]string, 0)
	
	start := l - 1 - cronElementsSize
	end := l - 2

	for i := start; i <= end; i++ {
		cronArgs = append(cronArgs, argsWithProg[i])
	}

	return "", nil
}
