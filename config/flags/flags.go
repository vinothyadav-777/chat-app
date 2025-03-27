package flags

import (
	"github.com/vinothyadav-777/chat-app/constants"
	"os"

	flag "github.com/spf13/pflag"
)

var (
	port           = flag.Int(constants.PortKey, constants.PortDefaultValue, constants.PortUsage)
	baseConfigPath = flag.String(constants.BaseConfigPathKey, constants.BaseConfigPathDefaultValue,
		constants.BaseConfigPathUsage)
)

func init() {
	flag.Parse()
}

// Env is the application.yml runtime environment
func Env() string {
	env := os.Getenv(constants.EnvKey)
	if env == "" {
		return constants.EnvDefaultValue
	}
	return env
}

// BaseConfigPath is the path that holds the configuration files
func BaseConfigPath() string {
	return *baseConfigPath + Env()
}

// Port is the application.yml port number where the process will be started
func Port() int {
	return *port
}
