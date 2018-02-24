package ewn

import (
	"github.com/alexflint/go-arg"
)

// Args cli args structure
type Args struct {
	Command          string   `arg:"-c"`
	Recipients       []string `arg:"-r"`
	Comment          string   `help:"Comment for email Message"`
	ValidExitCode    []int    `arg:"--valid-exitcodes" help:"Valid exitcodes for executed command"`
	Daemon           bool     `arg:"-d" help:"Daemonize process after start. DONT implemented"`
	DontDuplicate    bool     `arg:"--dont-duplicate" help:"Not run process when process with same key already run"`
	DontDuplicateKey string   `arg:"--dont-duplicate-key" help:"Default: --command value"`
	Retry            int      `help:"Retry run N times on fail."`
	RetrySleep       int      `arg:"--retry-sleep" help:"Sleep between retries (seconds)"`
	Config           string   `help:"Path to config file."`
	InitConfig       bool     `help:"Write default config to --config path and exit"`
}

// GetCliArgs return cli args structure
func GetCliArgs() (cli Args) {
	cli.ValidExitCode = []int{0}
	cli.Retry = 1
	cli.RetrySleep = 0
	cli.Config = "/etc/ewn.conf"

	arg.MustParse(&cli)

	if cli.DontDuplicateKey == "" {
		cli.DontDuplicateKey = cli.Command
	}
	return
}
