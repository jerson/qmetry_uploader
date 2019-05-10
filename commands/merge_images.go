package commands

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"

	"qmetry_uploader/modules/utils"
)

// MergeImagesOptions ...
type MergeImagesOptions struct {
	Input  string
	Output string
}

// MergeImages ...
func MergeImages(options MergeImagesOptions) error {

	files, err := ioutil.ReadDir(options.Input)
	if err != nil {
		return err
	}
	var filePaths []string
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		path := fmt.Sprintf("%s/%s", options.Input, file.Name())
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

	_ = os.MkdirAll(options.Output, 0777)
	output := fmt.Sprintf("%s/%s.png", options.Output, "merged")
	err = utils.MergeImages(filePaths, output)
	if err != nil {
		return err
	}
	log.Infof("output: %s\n", output)

	return nil
}
