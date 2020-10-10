package powershell_test

import (
	"fmt"
	"testing"

	"github.com/matryer/is"
	"github.com/simonjanss/go-powershell"
)

func TestNew(t *testing.T) {
	is := is.New(t)

	ps, err := powershell.New()
	is.NoErr(err)
	defer ps.Close()
	fmt.Println(ps.GetPid())
}

func TestExecute(t *testing.T) {
	is := is.New(t)

	ps, err := powershell.New()
	is.NoErr(err)

	var tt = []struct {
		name       string
		command    string
		shouldFail bool
		hasOutput  bool
	}{
		{name: "dir-command", command: "dir", shouldFail: false, hasOutput: true},
		{name: "cd-command", command: "cd ..", shouldFail: false, hasOutput: false},
		{name: "random-command", command: "this-should-fail", shouldFail: true, hasOutput: false},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			is := is.New(t)
			output, err := ps.Execute(tc.command)
			if !tc.shouldFail {
				is.NoErr(err)
			}
			if tc.hasOutput {
				is.True(output != nil)
			} else {
				is.True(output == nil)
			}
		})
	}

	err = ps.Close()
	is.NoErr(err)
}

func TestConcurrentExecute(t *testing.T) {
	is := is.New(t)

	ps, err := powershell.New()
	is.NoErr(err)

	go func() {
		_, err := ps.Execute("dir")
		is.Equal(err.Error(), "powershell: cannot execute command - powershell is busy")
	}()

	output, err := ps.Execute("echo hello")
	is.NoErr(err)
	is.True(string(output) == "hello")
}
