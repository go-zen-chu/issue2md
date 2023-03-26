package runner

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/go-zen-chu/issue2md/internal/config"
	"github.com/go-zen-chu/issue2md/internal/log"
)

// Runner defines interface for running general applications
type Runner interface {
	// Hand command line args to config for setup
	LoadConfigFromCommandArgs(args []string) error
	// Hand environment vars to config for setup
	LoadConfigFromEnvVars(envVars []string) error
	SetCommandHandler(ch CommandHandler)
	Run() error
}

type CommandHandler func(c *config.Config) error

type runner struct {
	appName string
	args    []string
	flgSet  *flag.FlagSet
	// general flag
	debug          bool
	help           bool
	cnf            config.Config
	commandHandler CommandHandler
}

func NewRunner(appName string) Runner {
	flgSet := flag.NewFlagSet(appName, flag.ExitOnError)
	return &runner{
		appName:        appName,
		flgSet:         flgSet,
		debug:          false,
		help:           false,
		cnf:            config.NewConfig(),
		commandHandler: nil,
	}
}

func (r *runner) LoadConfigFromCommandArgs(args []string) error {
	r.args = args
	debugVal := r.flgSet.Bool("debug", false, "Enable debug")
	helpVal := r.flgSet.Bool("help", false, "Show help")
	if !r.flgSet.Parsed() {
		if len(args) <= 1 {
			return fmt.Errorf("invalid args len: %d", len(os.Args))
		}
		if err := r.flgSet.Parse(args[1:]); err != nil {
			return fmt.Errorf("parse args: %s", err)
		}
		// visit specified flag
		r.flgSet.Visit(func(f *flag.Flag) {
			switch f.Name {
			case "debug":
				r.debug = *debugVal
			case "help":
				r.help = *helpVal
			}
		})
		if err := r.cnf.LoadFromCommandArgs(args); err != nil {
			return fmt.Errorf("while parsing args: %w", err)
		}
	}
	return nil
}

func (r *runner) LoadConfigFromEnvVars(envVars []string) error {
	return r.cnf.LoadFromEnvVars(envVars)
}

func (r *runner) buildHelpString() string {
	var sb strings.Builder
	sb.WriteString("usage: ")
	sb.WriteString(r.appName)
	sb.WriteString(" <flags>\n")
	// set print setting to string builder
	op := r.flgSet.Output()
	r.flgSet.SetOutput(&sb)
	r.flgSet.PrintDefaults()
	r.flgSet.SetOutput(op)
	return sb.String()
}

func (r *runner) SetCommandHandler(ch CommandHandler) {
	r.commandHandler = ch
}

func (r *runner) Run() error {
	if r.help {
		fmt.Println(r.buildHelpString())
		return nil
	}
	if err := log.Init(r.debug); err != nil {
		return fmt.Errorf("initializing log: %w", err)
	}
	log.Debugf("[Run] config: %+v", r.cnf)
	return nil
}

// pathValue is defined to handle path type argument
type pathValue struct {
	path string
}

// implements Value interface for flag value argument
func (pv *pathValue) String() string {
	return pv.path
}

// implements Value interface for flag value argument
func (pv *pathValue) Set(path string) error {
	if _, err := os.Stat(path); err != nil {
		return fmt.Errorf("not valid path %s: %w", path, err)
	}
	pv.path = path
	return nil
}
