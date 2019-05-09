package commands

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path/filepath"
	"qmetry_uploader/modules/config"
	"strings"
)

// MergeImages ...
func MergeImages() error {

	baseDir := config.Vars.Dir.Input
	files, err := ioutil.ReadDir(baseDir)
	if err != nil {
		return err
	}
	var filePaths []string
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		path := fmt.Sprintf("%s/%s", baseDir, file.Name())
		extension := filepath.Ext(path)
		if !contains([]string{".jpg", ".png"}, strings.ToLower(extension)) {
			log.Warnf("ignored file: %s", path)
			continue
		}
		filePaths = append(filePaths, path)

	}

	if len(filePaths) < 1 {
		return errors.New("images not found")
	}

	_ = os.MkdirAll(config.Vars.Dir.Output, 0777)
	output := fmt.Sprintf("%s/%s.png", config.Vars.Dir.Output, "merged")
	err = mergeImages(filePaths, output)
	if err != nil {
		return err
	}
	log.Infof("output: %s\n", output)

	return nil
}
