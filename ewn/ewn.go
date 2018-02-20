package ewn

import (
	"bytes"
	"github.com/mattn/go-isatty"
	"io"
	"os"
	"os/exec"
	"syscall"
	"time"
)

// Retry contains command execution result
type Retry struct {
	ExitCode  int
	Output    string
	StartTime time.Time
	EndTime   time.Time
	Duration  time.Duration
	Retry     int
}

// Message is output message structure
type Message struct {
	Args         Args
	Host         string
	Retries      []Retry
	GeneralError error
}

// Popen execute given command and return retry structure
func Popen(command string) (result Retry, err error) {
	var outB bytes.Buffer
	result.StartTime = time.Now().UTC()
	cmd := exec.Command("/bin/bash", "-c", command)
	if isatty.IsTerminal(os.Stdout.Fd()) {
		cmd.Stdout = io.MultiWriter(os.Stdout, &outB)
		cmd.Stderr = io.MultiWriter(os.Stdout, &outB)
	} else {
		cmd.Stdout = &outB
		cmd.Stderr = &outB
	}

	err2 := cmd.Start()
	_ = cmd.Wait()

	result.Output = outB.String()
	result.EndTime = time.Now().UTC()
	result.Duration = result.EndTime.Sub(result.StartTime)

	if err2 != nil {
		exitError := err.(*exec.ExitError)
		ws := exitError.Sys().(syscall.WaitStatus)
		result.ExitCode = ws.ExitStatus()
	} else {
		ws := cmd.ProcessState.Sys().(syscall.WaitStatus)
		result.ExitCode = ws.ExitStatus()
	}
	return result, nil
}
