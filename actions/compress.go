package actions

import (
	"github.com/urfave/cli"

	"qmetry_uploader/commands"
	"qmetry_uploader/modules/config"
	"qmetry_uploader/modules/prompt"
)

// Compress ...
func Compress(c *cli.Context) error {
	input := c.String("input")
	output := c.String("output")

	input = prompt.Dir("Input Dir", input, config.Vars.Dir.Input)
	output = prompt.Dir("Output Dir", output, config.Vars.Dir.Output)

	options := commands.CompressOptions{
		Input:  input,
		Output: output,
	}
	err := commands.Compress(options)
	return err
}
