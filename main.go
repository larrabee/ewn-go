package main

import (
	"github.com/larrabee/ewn-go/ewn"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cli := ewn.GetCliArgs()
	lock := ewn.Lock{Key: cli.DontDuplicateKey}
	msg := ewn.Message{
		Args: cli,
		Retries: make([]ewn.Retry, 0, cli.Retry),
	}

	cfg, err := ewn.GetConfig(cli.Config)
	if err != nil {
		panic(err)
	}
	if len(cli.Recipients) != 0 {
		cfg.Set("email.recipients", cli.Recipients)
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP)

	msg.Host, _ = os.Hostname()

	if cli.DontDuplicate {
		err := lock.Acquire()
		if err != nil {
			msg.GeneralError = err
			ewn.Notify(&msg, cfg)
			os.Exit(1)
		}
		defer lock.Release()
	}

RetryLoop:
	for i := 1; i <= cli.Retry; i++ {
		retry := ewn.Retry{Retry:i}
		retry, err := ewn.Popen(cli.Command, time.Duration(cli.Timeout)*time.Second)
		if err != nil {
			msg.GeneralError = err
			ewn.Notify(&msg, cfg)
			os.Exit(1)
		}
		msg.Retries = append(msg.Retries, retry)
		for _, v := range cli.ValidExitCode {
			if retry.ExitCode == v {
				break RetryLoop
			}
		}
	}
	ewn.Notify(&msg, cfg)
}
