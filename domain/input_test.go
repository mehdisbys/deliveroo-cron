package domain

import (
	"strings"
	"testing"
)

func TestInput(t *testing.T) {
	inputGetter := CommandLineInput{}

	tests := []struct {
		name        string
		commandLine string
	}{
		{
			commandLine: "application-commands file-name -arguments */15 0 1,15 * 1-5 /usr/bin/find",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, _ = inputGetter.GetData(strings.Split(test.commandLine, " "))
		})
	}
}
