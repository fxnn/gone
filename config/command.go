package config

import "fmt"

// Command represents a command that can be executed by the application.
type Command int

const (
	CommandHelp Command = iota
	CommandListen
	CommandExportTemplates
)

// String returns the string representation of the command, as it's to be used
// on the commandline.
func (c Command) String() string {
	switch c {
	case CommandHelp:
		return "help"
	case CommandListen:
		return "listen"
	case CommandExportTemplates:
		return "export-templates"
	}
	return ""
}

// Commands returns all valid command values.
func Commands() []Command {
	return []Command{CommandHelp, CommandListen, CommandExportTemplates}
}

// StringToCommand interprets the given string as a Command.
// It returns an error if the given string is no known command value.
func StringToCommand(s string) (Command, error) {
	for _, cmd := range Commands() {
		if s == cmd.String() {
			return cmd, nil
		}
	}
	return DefaultCommand, fmt.Errorf("invalid command: %s", s)
}
