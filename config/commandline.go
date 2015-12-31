package config

import (
	"flag"
	"fmt"
	"os"
)

var out = os.Stderr
var help bool = false

func init() {
	flag.BoolVar(&help, "help", false, "Displays this usage information")
	flag.BoolVar(&help, "h", false, "")
	flag.Usage = func() {
		fmt.Fprintln(out)
		PrintUsage()
	}
}

type commandlineConfig struct {
	command Command
}

func (c *commandlineConfig) Command() Command {
	return c.command
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
	fmt.Fprintf(out, "Usage: %s [-flag1 -flag2 ...] [command]", os.Args[0])
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
