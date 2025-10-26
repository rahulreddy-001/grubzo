package cmd

import (
	"grubzo/internal/config"
)

var c *config.Config

func loadConfig() error {
	if cfg, err := config.LoadConfig(); err != nil {
		return err
	} else {
		c = cfg
	}
	return nil
}
