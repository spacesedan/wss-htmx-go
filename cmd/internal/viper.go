package internal

import (
	"log/slog"

	"github.com/spacesedan/wss-htmx-go/internal"
	"github.com/spf13/viper"
)

// NewViper reads config file
func NewViper(logger *slog.Logger) (*viper.Viper, error) {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")


	if err := v.ReadInConfig(); err != nil {
		logger.Error("Failed to read config file", slog.String("err", err.Error()))
		return nil,  internal.WrapErrorf(err,internal.ErrorCodeUnknown, "viper.ReadInConfig FAILED TO READ" )
	}

	logger.Info("Viper Config Loaded: OK")
	return v, nil
}
