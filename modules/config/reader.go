//Package config ...
package config

import (
	"path/filepath"

	"github.com/jinzhu/configor"

	"qmetry_uploader/modules/utils"
)

// Dir ...
type Dir struct {
	Input  string `toml:"input" default:""`
	Output string `toml:"output" default:"./output"`
}

// Nexus ...
type Nexus struct {
	Username string `toml:"username" default:""`
	Password string `toml:"password" default:""`
	Project  string `toml:"project" default:"mi-banco"`
	// project platform name
	Server string `toml:"server" default:"http://mb-nexus.westus.cloudapp.azure.com/repository/%s-%s/builds/%s"`
}

// Binary ...
type Binary struct {
	ADB       string `toml:"adb" default:"adb"`
	XCode     string `toml:"xcode" default:"Xcode"`
	Automator string `toml:"automator" default:"/usr/bin/automator"`
}

//Vars ...
var Vars = struct {
	Debug   bool   `toml:"debug" default:"false"`
	Version string `toml:"version" default:"0.0.2"`
	Dir     Dir    `toml:"dir"`
	Nexus   Nexus  `toml:"nexus"`
	Binary  Binary `toml:"Binary"`
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

	config := configor.New(&configor.Config{ENVPrefix: "APP", Debug: false, Verbose: false})
	if utils.ExistsFile(file) {
		return config.Load(&Vars, file)
	}
	return config.Load(&Vars)
}
