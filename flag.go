package envflag

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// Option type for configing the parsing.
type Option option
type option func(*config)

// FlagConverterFunc functions are used to update flag name for envvar lookups.
type FlagConverterFunc func(flag string) string

// UsageUpdaterFunc functions are used to update flag usage strings.
type UsageUpdaterFunc func(key, usage string) string

type config struct {
	prefix        string
	getenv        func(string) string
	flagConverter FlagConverterFunc
	usageUpdater  UsageUpdaterFunc
}

// FlagConverter converts the flag name to an environment variable key.
func FlagConverter(f FlagConverterFunc) Option {
	return func(cfg *config) {
		cfg.flagConverter = f
	}
}

// Prefix sets the prefix for all envvars.
// This is applied after the FlagConverter function.
func Prefix(s string) Option {
	return func(cfg *config) {
		cfg.prefix = s
	}
}

// Getenv is used to fetch the environment, defaults to os.Getenv.
func Getenv(f func(string) string) Option {
	return func(cfg *config) {
		cfg.getenv = f
	}
}

// UsageUpdater enables configurable flag usage updates.
func UsageUpdater(f UsageUpdaterFunc) Option {
	return func(cfg *config) {
		cfg.usageUpdater = f
	}
}

// UsagePrefixer prefixes the flag usage with [envvar].
//
// Example output:
//
//	Usage of example:
//	  -verbose
//	        [VERBOSE] Be verbose
func UsagePrefixer() Option {
	return func(cfg *config) {
		cfg.usageUpdater = usagePrefixer
	}
}

// UsageSuffixer suffixes the flag usage with [envvar].
//
// Example output:
//
//	Usage of example:
//	  -verbose
//	        Be verbose [VERBOSE]
func UsageSuffixer() Option {
	return func(cfg *config) {
		cfg.usageUpdater = usageSuffixer
	}
}

// ParseWithEnv parses the command-line flags from os.Args[1:]. Must be called
// after all flags are defined and before flags are accessed by the program.
func Parse(opts ...Option) error {
	return parseFlagSetWithEnv(flag.CommandLine, os.Args[1:], opts...)
}

// ParseFlagSetWithEnv parses flag definitions from the argument list, which
// should not include the command name. Must be called after all flags in the
// FlagSet are defined and before flags are accessed by the program. The return
// value will be ErrHelp if -help or -h were set but not defined.
func ParseFlagSet(fs *flag.FlagSet, arguments []string, opts ...Option) error {
	return parseFlagSetWithEnv(fs, arguments, opts...)
}

func parseFlagSetWithEnv(fs *flag.FlagSet, arguments []string, opts ...Option) error {
	if fs.Parsed() {
		return fmt.Errorf("flag has already been parsed")
	}
	cfg := &config{
		flagConverter: flagToEnv(),
		getenv:        os.Getenv,
	}
	for _, o := range opts {
		o(cfg)
	}
	var nerr error
	fs.VisitAll(func(f *flag.Flag) {
		envKey := cfg.flagConverter(f.Name)
		envKey = cfg.prefix + envKey
		if cfg.usageUpdater != nil {
			f.Usage = cfg.usageUpdater(envKey, f.Usage)
		}
		if envVal := cfg.getenv(envKey); envVal != "" {
			err := f.Value.Set(envVal)
			if err != nil && nerr == nil {
				nerr = err
			}
		}
	})
	if nerr != nil {
		return nerr
	}
	return fs.Parse(arguments)
}

// convert this-format to THIS_FORMAT
func flagToEnv() func(string) string {
	re := regexp.MustCompile("[^a-zA-Z0-9_]+")
	return func(f string) string {
		return strings.ToUpper(re.ReplaceAllLiteralString(f, "_"))
	}
}

func usageSuffixer(key, usage string) string {
	return fmt.Sprintf("%s [%s]", usage, key)
}

func usagePrefixer(key, usage string) string {
	return fmt.Sprintf("[%s] %s", key, usage)
}
