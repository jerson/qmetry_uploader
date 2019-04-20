//Package config ...
package config

import (
	"github.com/jinzhu/configor"
	"path/filepath"
)

// Dir ...
type Dir struct {
	Images string `toml:"images" default:"./images"`
	Output string `toml:"output" default:"./output"`
	Result string `toml:"result" default:"./results"`
}

//Vars ...
var Vars = struct {
	Debug   bool   `toml:"debug" default:"false"`
	Version string `toml:"version" default:"latest"`
	Dir     Dir    `toml:"dir"`
}{}

//ReadDefault ...
func ReadDefault() error {
	file, err := filepath.Abs("./config.toml")
	if err != nil {
		return err
	}
	return Read(file)
}

//Read ...
func Read(file string) error {
	return configor.New(&configor.Config{ENVPrefix: "APP", Debug: false, Verbose: false}).Load(&Vars, file)
}
