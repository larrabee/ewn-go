package main

import (
	"fmt"
	ewn "github.com/larrabee/ewn-go/ewn"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cli := ewn.GetCliArgs()
	lock := ewn.Lock{Key: cli.DontDuplicateKey}
	msg := ewn.Message{Args: cli}

	if cli.InitConfig {
		err := ewn.InitConfig(cli.Config)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Config initialization failed with error: ", err)
			os.Exit(1)
		}
		fmt.Printf("Default config file successfully created in %s file\nPlease update it with your values\n", cli.Config)
		os.Exit(0)
	}

	cfg, err1 := ewn.GetConfig(cli.Config)
	if err1 != nil {
		panic(err1)
	}
	if len(cli.Recipients) != 0 {
		cfg.Set("email.recipients", cli.Recipients)
	}
	if cli.Daemon {
		fmt.Fprintln(os.Stderr, "Daemonization not implementet. Running in normal mode.")
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP)

	msg.Host, _ = os.Hostname()

	if cli.DontDuplicate {
		err2 := lock.Acquire()
		if err2 != nil {
			msg.GeneralError = err2
			_ = ewn.Notify(msg)
			os.Exit(1)
		}
	}

RetryLoop:
	for retryCounter := 1; retryCounter <= cli.Retry; retryCounter++ {
		fmt.Println(retryCounter)
		retry := ewn.Popen(cli.Command)
		msg.Retries = append(msg.Retries, retry)
		for _, v := range cli.ValidExitCode {
			if retry.ExitCode == v {
				break RetryLoop
			}
		}
	}
	_ = ewn.Notify(msg)
}
