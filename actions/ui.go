package actions

import (
	"github.com/urfave/cli"

	"qmetry_uploader/commands"
)

// UI ...
func UI(c *cli.Context) error {

	return commands.UI()
}
