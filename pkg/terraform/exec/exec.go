package exec

import (
	"bufio"
	"bytes"
	"io"
	"os"
	goexec "os/exec"
)

func _runner(cmd string, dir string, args []string, stdout, stderr io.Writer) int {
	baseCommand := "terraform"
	cmdArgs := []string{}
	cmdArgs = append(cmdArgs, "-chdir=./"+dir)
	cmdArgs = append(cmdArgs, cmd)
	cmdArgs = append(cmdArgs, args...)
	c := goexec.Command(baseCommand, cmdArgs...)

	c.Stdout = stdout
	c.Stderr = stderr
	if err := c.Run(); err != nil {
		if exitError, ok := err.(*goexec.ExitError); ok {
			return exitError.ExitCode()
		}
	}
	return 0
}

// Apply is wrapper around `terraform apply` subcommand.
func Apply(datadir string, args []string) int {
	return _runner("apply", datadir, args, os.Stdout, os.Stderr)
}

// Destroy is wrapper around `terraform destroy` subcommand.
func Destroy(datadir string, args []string) int {
	return _runner("destroy", datadir, args, os.Stdout, os.Stderr)
}

// Destroy is wrapper around `terraform output` subcommand.
func Output(datadir string, args []string) (string, int) {
	var b bytes.Buffer
	bw := bufio.NewWriter(&b)
	exitstatus := _runner("output", datadir, args, bw, os.Stderr)
	if exitstatus != 0 {
		return "", exitstatus
	}
	return b.String(), exitstatus
}

// Init is wrapper around `terraform init` subcommand.
func Init(datadir string, args []string) int {
	return _runner("init", datadir, args, os.Stdout, os.Stderr)
}
