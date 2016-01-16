package config

import (
	"flag"
	"fmt"
	"os"
	"time"
)

var out = os.Stderr
var (
	help                            bool
	bindAddress                     string
	requireSSLHeader                string
	templatePath                    string
	bruteforceMaxDelayMillis        int
	bruteforceDelayStepMillis       int
	bruteforceDropDelayAfterMinutes int
)

func init() {
	flag.BoolVar(&help, "help", false,
		"Displays this usage information")
	flag.BoolVar(&help, "h", false, "")

	flag.StringVar(&bindAddress, "bind", DefaultBindAddress,
		"The `address` and/or port to listen on")
	flag.StringVar(&requireSSLHeader, "require-ssl-header", DefaultRequireSSLHeader,
		"The `name` of a header to be required when logging in")
	flag.StringVar(&templatePath, "template", DefaultTemplatePath,
		"The `path` to a directory containing custom templates")

	flag.IntVar(&bruteforceMaxDelayMillis, "bruteforce-max-delay",
		int(DefaultBruteforceMaxDelay/time.Millisecond),
		"The max number of `millis` to delay login requests.")
	flag.IntVar(&bruteforceDelayStepMillis, "bruteforce-delay-step",
		int(DefaultBruteforceDelayStep/time.Millisecond),
		"The number of `millis` to delay login requests per recently failed login attempt.")
	flag.IntVar(&bruteforceDropDelayAfterMinutes, "bruteforce-drop-delay-after",
		int(DefaultBruteforceDropDelayAfter/time.Minute),
		"The lifetime of each delay in `minutes` after the last failed login attempt.")

	flag.Usage = func() {
		fmt.Fprintln(out)
		PrintUsage()
	}
}

func FromCommandline() Config {
	var command = parseCommandline()

	var c = Config{}
	c.Command = command
	c.BindAddress = bindAddress
	c.RequireSSLHeader = requireSSLHeader
	c.TemplatePath = templatePath
	c.BruteforceMaxDelay = time.Duration(bruteforceMaxDelayMillis) * time.Millisecond
	c.BruteforceDelayStep = time.Duration(bruteforceDelayStepMillis) * time.Millisecond
	c.BruteforceDropDelayAfter = time.Duration(bruteforceDropDelayAfterMinutes) * time.Minute

	return c
}

func parseCommandline() Command {
	flag.Parse()

	if flag.NArg() > 1 {
		fmt.Fprintln(out, "No more than one command allowed")
		PrintUsage()
		os.Exit(2)
	}

	if flag.NArg() == 1 {
		if cmd, err := StringToCommand(flag.Arg(0)); err != nil {
			fmt.Fprintln(out, err)
			PrintUsage()
			os.Exit(2)
		} else {
			return cmd
		}
	}

	if help {
		return CommandHelp
	}

	return DefaultCommand
}

func PrintUsage() {
	fmt.Fprintf(out, "Usage: %s [-flags ...] [command]", os.Args[0])
	fmt.Fprintln(out)

	flag.PrintDefaults()
	fmt.Fprintln(out)

	fmt.Fprintln(out, "Valid commands are:")
	for _, cmd := range Commands() {
		fmt.Fprintf(out, "  %s", cmd)
		if cmd == DefaultCommand {
			fmt.Fprintf(out, " (default)")
		}
		fmt.Fprintln(out)
	}
}
