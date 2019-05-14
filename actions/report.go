package actions

import (
	"encoding/json"

	"github.com/urfave/cli"

	log "github.com/sirupsen/logrus"

	"qmetry_uploader/commands"
	"qmetry_uploader/modules/config"
	"qmetry_uploader/modules/prompt"
)

// Report ...
func Report(c *cli.Context) error {

	input := c.String("input")
	input = prompt.Dir("Input Dir", input, config.Vars.Dir.Input)

	options := commands.ReportOptions{
		Input: input,
	}

	data, err := commands.Report(options)
	if data != nil {
		printJSON(data)
	}
	return err
}

func printJSON(data interface{}) {

	output, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		panic(err)
	}
	log.Debug(string(output))
}
