package main

import (
	"deliveroo-cron/domain"
	"log"
	"os"
)

func main() {

	c := domain.CommandLineInput{}

	cronargs, err := c.GetDataFromCommandLine()
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	parsedExp, err := domain.ParseExpression(cronargs)
	if err != nil {
		log.Print(err)
		os.Exit(2)
	}

	d := domain.NewCommandLineDisplay(c.GetProgramToExec(), *parsedExp)

	d.Display()
}
