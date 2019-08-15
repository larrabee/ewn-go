package ewn

import (
	"bytes"
	"context"
	"git.wsmgroup.ru/go-modules/utils"
	"github.com/mattn/go-isatty"
	"io"
	"os"
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
	return result, nil
}
