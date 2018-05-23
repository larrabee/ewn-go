package ewn

import (
	"bytes"
	"github.com/mattn/go-isatty"
	"io"
	"os"
	"os/exec"
	"syscall"
	"time"
	"github.com/kr/pty"
)

//Retry contains command execution result
type Retry struct {
	ExitCode  int
	Output    string
	StartTime time.Time
	EndTime   time.Time
	Duration  time.Duration
	Retry     int
}

//Message is output message structure
type Message struct {
	Args         Args
	Host         string
	Retries      []Retry
	GeneralError error
}

//Popen execute given command and return retry structure
func Popen(command string, timeout time.Duration, tty bool) (result Retry, err error) {
	var outB bytes.Buffer
	var timer *time.Timer
	result.StartTime = time.Now().UTC()
	cmd := exec.Command("/bin/bash", "-c", command)

	if timeout > time.Duration(0) {
		timer = time.AfterFunc(timeout, func() {
			cmd.Process.Kill()
		})
	}

	if isatty.IsTerminal(os.Stdout.Fd()) && !tty {
		cmd.Stdout = io.MultiWriter(os.Stdout, &outB)
		cmd.Stderr = io.MultiWriter(os.Stdout, &outB)
	} else if isatty.IsTerminal(os.Stdout.Fd()) && tty {
		ptmx, err := pty.Start(cmd)
		if err != nil {
			return result, err
		}
		defer ptmx.Close()
		go func() { io.Copy(io.MultiWriter(os.Stdout, &outB), ptmx) }()
	} else if tty {
		ptmx, err := pty.Start(cmd)
		if err != nil {
			return result, err
		}
		defer ptmx.Close()
		go func() { io.Copy(&outB, ptmx) }()
	}else {
		cmd.Stdout = &outB
		cmd.Stderr = &outB
	}

	if !tty {
		err = cmd.Start()
		if err != nil {
			return
		}
	}

	err = cmd.Wait()

	if timeout > time.Duration(0) {
		timer.Stop()
	}

	result.Output = outB.String()
	result.EndTime = time.Now().UTC()
	result.Duration = result.EndTime.Sub(result.StartTime)

	if err != nil {
		exitError := err.(*exec.ExitError)
		ws := exitError.Sys().(syscall.WaitStatus)
		result.ExitCode = ws.ExitStatus()
	} else {
		ws := cmd.ProcessState.Sys().(syscall.WaitStatus)
		result.ExitCode = ws.ExitStatus()
	}
	return result, nil
}
