package cmd

import (
	"grubzo/internal/config"
	"sync"
)

var once sync.Once
var c *config.Config

func loadConfig() error {
	var e error
	once.Do(func(){
		if cfg, err := config.LoadConfig(); err != nil {
			e = err
		} else {
			c = cfg
		}
	})
	return e
}
