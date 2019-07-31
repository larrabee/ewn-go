package ewn

import (
	"fmt"
	"github.com/alexflint/go-arg"
	"os"
)

// Args cli args structure
type Args struct {
	Command          string   `arg:"-c"`
	Recipients       []string `arg:"-r"`
	Comment          string   `help:"Comment for email Message"`
	ValidExitCode    []int    `arg:"--valid-exitcodes" help:"Valid exitcodes for executed command"`
	DontDuplicate    bool     `arg:"--dont-duplicate" help:"Not run process when process with same key already run"`
	DontDuplicateKey string   `arg:"--dont-duplicate-key" help:"Default: --command value"`
	Retry            int      `help:"Retry run N times on fail."`
	Timeout          int      `arg:"-t" help:"Kill process after N seconds. [default: 0]"`
	RetrySleep       int      `arg:"--retry-sleep" help:"Sleep between retries (seconds)"`
	Config           string   `help:"Path to config file."`
	InitConfig       bool     `help:"Write default config to --config path and exit"`
	Tty              bool     `help:"Reserved arg"`
	Daemon           bool     `arg:"-d" help:"Reserved arg"`
}

// GetCliArgs return cli args structure
func GetCliArgs() (cli Args) {
	cli.ValidExitCode = []int{0}
	cli.Retry = 1
	cli.RetrySleep = 0
	cli.Config = "/etc/ewn.conf"
	cli.Timeout = 0

	p := arg.MustParse(&cli)

	if cli.DontDuplicateKey == "" {
		cli.DontDuplicateKey = cli.Command
	}

	if cli.InitConfig {
		err := InitConfig(cli.Config)
		if err != nil {
			p.Fail(fmt.Sprintf("Config initialization failed with error: %s", err))
		}
		fmt.Printf("Default config file successfully created: %s\nPlease update it with your values\n", cli.Config)
		os.Exit(0)
	}

	if (cli.InitConfig == false) && (cli.Command == "") {
		p.Fail("error: arg --command or --init-config are required")
	}

	if cli.Daemon {
		fmt.Fprintln(os.Stderr, "Daemonization (-d) not implemented. Running in normal mode.")
	}
	if cli.Tty {
		fmt.Fprintln(os.Stderr, "Fake TTY not implemented. Running in normal mode.")
	}

	return
}
