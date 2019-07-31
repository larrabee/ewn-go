package ewn

import (
	"bytes"
	"context"
	"git.wsmgroup.ru/go-modules/utils"
	"github.com/kr/pty"
	"github.com/mattn/go-isatty"
	"io"
	"os"
	"os/exec"
	"syscall"
	"time"
)

//Retry contains command execution result
type Retry struct {
	ExitCode  int
	Output    string
	StartTime time.Time
	EndTime   time.Time
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
func Popen(command string, timeout time.Duration) (Retry, error) {
	ctx := context.Background()
	outBuffer := &bytes.Buffer{}
	result := Retry{}
	if timeout > 0 {
		ctx, _ = context.WithTimeout(ctx, timeout)
	}

	cmd := utils.NewExec()
	cmd.WithContext(ctx)
	cmd.Args = append(cmd.Args, []string{"/bin/bash", "-c", command}...)

	if isatty.IsTerminal(os.Stdout.Fd()) {
		cmd.Stdout = io.MultiWriter(os.Stdout, outBuffer)
		cmd.Stderr = io.MultiWriter(os.Stdout, outBuffer)
	} else {
		cmd.Stdout = outBuffer
		cmd.Stderr = outBuffer
	}


	result.StartTime = time.Now().UTC()
	err := cmd.Exec()
	result.EndTime = time.Now().UTC()
	if err != nil {
		return result, err
	}

	result.Output = outBuffer.String()
	//result.Duration = result.EndTime.Sub(result.StartTime)
	return result, nil
}

//Popen execute given command and return retry structure
func PopenOld(command string, timeout time.Duration, tty bool) (result Retry, err error) {
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
	} else {
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
	//result.Duration = result.EndTime.Sub(result.StartTime)

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
