package command

import (
	"errors"

	"github.com/spf13/viper"
)

// nolint: gochecknoinits
func init() {
	viper.SetEnvPrefix("ein")
	viper.BindEnv("runtime_path")
	viper.BindEnv("module_root_path")
}

func getRuntimePath() (string, error) {
	s := viper.GetString("runtime_path")

	if s == "" {
		return "", newEnvironmentVariableError("EIN_RUNTIME_PATH")
	}

	return s, nil
}

func getModulesRootPath() (string, error) {
	s := viper.GetString("module_root_path")

	if s == "" {
		return "", newEnvironmentVariableError("EIN_MODULE_ROOT_PATH")
	}

	return s, nil
}

func newEnvironmentVariableError(s string) error {
	return errors.New(s + " environment variable not set")
}
