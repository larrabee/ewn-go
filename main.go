package main

import (
	"fmt"
	"./ewn"
	"syscall"
	"os/signal"
	"os"
	"time"
)



func main() {
	cli := ewn.GetCliArgs()
	lock := ewn.Lock{Key: cli.DontDuplicateKey}
	msg := ewn.Message{Args: cli}
	
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
	
	for retryCounter := 1; retryCounter <= cli.Retry ; retryCounter++  {
		fmt.Println(retryCounter)
		_ = ewn.Popen(cli.Command)
		
		//fmt.Printf("Out: %scode: %d\ndur: %s\n", output, exitCode, duration)
	}
	_ = ewn.Notify(msg)
	//fmt.Println(args)
	//fmt.Println(MessageHeader)
	//fmt.Printf("%+v", cfg)
	time.Sleep(1* time.Second)
}
