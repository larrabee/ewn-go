package ewn

import (
  "bytes"
  "syscall"
  "time"
  "os/exec"
)

type Retry struct {
	ExitCode int
	Output string
	StartTime time.Time
	EndTime time.Time
	Duration time.Duration
	Retry int
}

type Message struct {
	Args Args
	Host string
	Retries []Retry
	GeneralError error
}

func Popen(command string) (result Retry) {
	var outB bytes.Buffer
  result.StartTime = time.Now().UTC()
	cmd := exec.Command("/bin/bash", "-c", command)
	cmd.Stdout = &outB
	cmd.Stderr = &outB
	err := cmd.Run()
	
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
	return
}
