package actions

import (
	"github.com/urfave/cli"

	"qmetry_uploader/commands"
)

// GUI ...
func GUI(c *cli.Context) error {

	return commands.GUI()
}
