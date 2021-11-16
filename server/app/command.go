package app

import "github.com/urfave/cli/v2"

func beforeCommand(c *cli.Context) error {
	if err := getServerName(c); err != nil {
		return err
	}
	return nil
}
