package ewn

import (
  "github.com/alexflint/go-arg"
)

type Args struct {
	Command string `arg:"required, -c"`
	Recipients []string `arg:"-r,separate"`
	Comment string `help:"Comment for email Message"`
	ValidExitCode []int `arg:"--valid-exitcodes, separate" help:"Valid exitcodes for executed command"`
	Daemon bool `arg:"-d" help:"Daemonize process after start"`
	DontDuplicate bool `arg:"--dont-duplicate" help:"Not run process when process with same key already run"`
	DontDuplicateKey string `arg:"--dont-duplicate-key" help:"Default: --command value"`
	Retry int `help:"Retry run N times on fail. Default: 1 (no retries)"`
	RetrySleep int `arg:"--retry-sleep" help:"Sleep between retries (seconds). Default: 0"`
	Config string `help:"Path to config file. Default: /etc/ewn.conf"`
  InitConfig bool `help:"Write default config to --config path and exit"`
}

func GetCliArgs() (cli Args) {
  cli.ValidExitCode = []int{0,}
  cli.Retry = 1
  cli.RetrySleep = 0
  cli.Config = "/etc/ewn.conf"
  
  arg.MustParse(&cli)
  
  if cli.DontDuplicateKey == "" {
		cli.DontDuplicateKey = cli.Command
	}
  return
}
