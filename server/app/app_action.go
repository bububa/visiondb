package app

import (
	"github.com/urfave/cli/v2"

	"github.com/bububa/visiondb/server/service"
)

func beforeAction(c *cli.Context) error {
	if err := loadConfigAction(c); err != nil {
		return err
	}
	if err := InitLogger(); err != nil {
		return err
	}
	if err := service.Init(&config); err != nil {
		return err
	}
	return nil
}

func afterAction(c *cli.Context) error {
	return nil
}

func deferFunc() {
	service.Close()
}
