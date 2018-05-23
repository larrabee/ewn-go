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
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
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
func Popen(command string, timeout time.Duration, tty bool) (result Retry, err error) {
	var outB bytes.Buffer
	var timer *time.Timer
	result.StartTime = time.Now().UTC()
	cmd := exec.Command("/bin/bash", "-c", command)

	if isatty.IsTerminal(os.Stdout.Fd()) {
		cmd.Stdout = io.MultiWriter(os.Stdout, &outB)
		cmd.Stderr = io.MultiWriter(os.Stdout, &outB)
	} else {
		cmd.Stdout = &outB
		cmd.Stderr = &outB
	}

	if tty {
		ptmx, err := pty.Start(cmd)
		if err != nil {
			return result, err
		}
		defer ptmx.Close()

		if err := pty.InheritSize(os.Stdin, ptmx); err != nil {
			fmt.Printf("error resizing pty: %s", err)
		}
		
		oldState, err := terminal.MakeRaw(int(os.Stdin.Fd()))
		if err != nil {
			return result, err
		}
		defer terminal.Restore(int(os.Stdin.Fd()), oldState)
		go func() { io.Copy(cmd.Stdout, ptmx) }()
		cmd.Stdout = io.MultiWriter(os.Stdout, &outB)
		cmd.Stderr = io.MultiWriter(os.Stdout, &outB)
	}

	if timeout > time.Duration(0) {
		timer = time.AfterFunc(timeout, func() {
			cmd.Process.Kill()
		})
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
