package app

import (
	"github.com/jinzhu/configor"
	"github.com/urfave/cli/v2"

	"github.com/bububa/visiondb/server/conf"
)

func loadConfigAction(c *cli.Context) error {
	configPath := c.String("config")
	if err := loadConfig(configPath, &config); err != nil {
		return err
	}
	if c.IsSet("port") {
		config.Port = c.Int("port")
	}
	if c.IsSet("debug") {
		config.Debug = c.Bool("debug")
	}
	if c.IsSet("log") {
		config.LogPath = c.String("log")
	}
	return nil
}

func loadConfig(configPath string, cfg *conf.Config) error {
	return configor.New(&configor.Config{
		Verbose:              false,
		ErrorOnUnmatchedKeys: true,
		Environment:          "production",
	}).Load(cfg, configPath)
}

func getServerName(c *cli.Context) error {
	serverName = "visiondb"
	return nil
}
