package command

import (
	"errors"

	"github.com/spf13/viper"
)

// nolint: gochecknoinits
func init() {
	viper.SetEnvPrefix("ein")
	viper.BindEnv("root")
}

func getRuntimeRoot() (string, error) {
	s := viper.GetString("root")

	if s == "" {
		return "", errors.New("EIN_ROOT environment variable not set")
	}

	return s, nil
}
