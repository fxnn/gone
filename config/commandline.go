package config

import (
	"flag"
	"fmt"
	"os"
)

var out = os.Stderr
var (
	help        bool
	bindAddress string
)

func init() {
	flag.BoolVar(&help, "help", false, "Displays this usage information")
	flag.BoolVar(&help, "h", false, "")
	flag.StringVar(&bindAddress, "bind", DefaultBindAddress, "The `address` and/or port to listen on")
	flag.Usage = func() {
		fmt.Fprintln(out)
		PrintUsage()
	}
}

type commandlineConfig struct {
	// NOTE: We don't store argument values inside the struct atm,
	// as those things are global per application instance anyways
	command Command
}

func (c *commandlineConfig) Command() Command {
	return c.command
}

func (c *commandlineConfig) BindAddress() string {
	return bindAddress
}

func FromCommandline() Config {
	var config = &commandlineConfig{}
	config.parseCommandline()
	return config
}

func (c *commandlineConfig) parseCommandline() {
	flag.Parse()

	if flag.NArg() > 1 {
		fmt.Fprintln(out, "No more than one command allowed")
		PrintUsage()
		os.Exit(2)
	} else if flag.NArg() == 1 {
		if cmd, err := StringToCommand(flag.Arg(0)); err == nil {
			c.command = cmd
		} else {
			fmt.Fprintln(out, err)
			PrintUsage()
			os.Exit(2)
		}
	} else if help {
		c.command = CommandHelp
	} else {
		c.command = DefaultCommand
	}
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