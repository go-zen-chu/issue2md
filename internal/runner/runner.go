package runner

import (
	"flag"
	"fmt"
	"strings"

	"github.com/go-zen-chu/issue2md/internal/config"
	"github.com/go-zen-chu/issue2md/internal/log"
)

/*
Runner defines interface for running general applications.
Runner depends on internal packages such as config, log, ... and work as wrapper to initialize and setup them.
*/
type Runner interface {
	// Hand environment vars to config for setup
	LoadConfigFromEnvVars(envVars []string) error
	// Hand command line args to config for setup
	LoadConfigFromCommandArgs(args []string) error
	// Run execute runner with other setup processes
	Run(ch ConfigHandler) error
}

type ConfigHandler func(c config.Config) error

type runner struct {
	appName string
	args    []string
	flgSet  *flag.FlagSet
	debug   bool
	help    bool
	cnf     config.Config
}

func NewRunner(appName string) Runner {
	return &runner{
		appName: appName,
		flgSet:  flag.NewFlagSet(appName, flag.ContinueOnError),
		cnf:     config.NewConfig(),
	}
}

func (r *runner) LoadConfigFromEnvVars(envVars []string) error {
	return r.cnf.LoadFromEnvVars(envVars)
}

func (r *runner) LoadConfigFromCommandArgs(args []string) error {
	if len(args) <= 1 {
		// no args given
		return nil
	}
	r.args = args
	debugVal := r.flgSet.Bool("debug", false, "Enable debug")
	helpVal := r.flgSet.Bool("help", false, "Show help")
	r.cnf.SetupCommandArgs(r.flgSet)
	var err error
	if !r.flgSet.Parsed() {
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
			default:
				err = r.cnf.LoadFromCommandArgs(f.Name)
			}
		})
	}
	return err
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

func (r *runner) Run(ch ConfigHandler) error {
	if r.help {
		fmt.Println(r.buildHelpString())
		return nil
	}
	if err := log.Init(r.debug); err != nil {
		return fmt.Errorf("initializing log: %w", err)
	}
	defer log.Close()
	log.Debugf("[Run] config: %+v", r.cnf)
	if err := ch(r.cnf); err != nil {
		return fmt.Errorf("handling config %s: %w", r.cnf, err)
	}
	return nil
}
