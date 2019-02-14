package command

import (
	"errors"

	"github.com/spf13/viper"
)

// nolint: gochecknoinits
func init() {
	viper.SetEnvPrefix("ein")
	viper.BindEnv("runtime_path")
}

func getRuntimePath() (string, error) {
	s := viper.GetString("runtime_path")

	if s == "" {
		return "", errors.New("EIN_RUNTIME_PATH environment variable not set")
	}

	return s, nil
}
