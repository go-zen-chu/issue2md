package config

type Config interface {
	LoadFromEnvVars(envVars []string) error
	LoadFromCommandArgs(args []string) error
}

type config struct {
}

func NewConfig() Config {
	return &config{}
}

func (c *config) LoadFromEnvVars(envVars []string) error {
	return nil
}
func (c *config) LoadFromCommandArgs(args []string) error {
	return nil
}
