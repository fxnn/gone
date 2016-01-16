package config

import "time"

// Config is the configuration for the application.
// A Config instance might be created from different sources, like a
// configuration file or command line switches.
type Config interface {

	// Command is the command to be executed by the application.
	// This defaults to the DefaultCommand constant.
	Command() Command

	// BindAddress is the network address the application will listen on.
	// This defaults to the DefaultListenAddress constant.
	BindAddress() string

	// TemplatePath is the path to the directory containing custom templates.
	// This defaults to the empty string, meaning the static templates
	// delivered with the application are used.
	TemplatePath() string

	// BruteforceMaxDelay is the maximum amount of time a login request is
	// delayed in order to prevent bruteforce attacks.
	BruteforceMaxDelay() time.Duration

	// BruteforceDelayStep configures how fast login requests will take longer
	// after failed login attempts.
	// After a failed attempt, the next attempt will take BruteforceDelayStep()
	// longer for the same user; BruteforceDelayStep() / 5 longer for the same
	// IP address and BruteforceDelayStep() / 20 longer independent of user and
	// IP address.
	BruteforceDelayStep() time.Duration

	// BruteforceDropDelayAfter configures after what time after the last failed
	// login attempt to drop the delays.
	BruteforceDropDelayAfter() time.Duration
}

const (
	DefaultCommand                  = CommandListen
	DefaultBindAddress              = ":8080"
	DefaultTemplatePath             = ""
	DefaultBruteforceMaxDelay       = 20 * time.Second
	DefaultBruteforceDelayStep      = 1 * time.Second
	DefaultBruteforceDropDelayAfter = 4 * time.Hour
)
