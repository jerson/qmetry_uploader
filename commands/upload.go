package commands

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"

	log "github.com/sirupsen/logrus"
)

type UploadOptions struct {
	File string
}
type UploadNexusOptions struct {
	UploadOptions
	Server   string
	Project  string
	Password string
	Username string
}

// UploadNexus
//
// initial script for bash
//#!/bin/bash
//# install: ln -sf $PWD/upload.sh /usr/local/bin/upload
//# usage: upload file.apk
//# usage2: upload file.ipa
//
//PLATFORM="ios"
//
//if [[ "$1" =~ ".apk" ]]
//then
//  PLATFORM="android"
//fi
//URL="http://mb-nexus.westus.cloudapp.azure.com/repository/mi-banco-$PLATFORM/builds/$1"
//echo "uploading: $1 => $PLATFORM";
//curl -u "username:pass"  --upload-file $1 "$URL";
//echo "URL: $URL";
func UploadNexus(options UploadNexusOptions) error {

	extension := path.Ext(options.File)
	platform := "common"
	if extension == ".apk" {
		platform = "android"
	} else if extension == ".ipa" || extension == ".zip" {
		platform = "ios"
	}

	url := fmt.Sprintf(options.Server, options.Project, platform, options.File)

	log.WithFields(log.Fields{
		"project":  options.Project,
		"platform": platform,
		"file":     options.File,
	}).Infof("uploading to %s", url)

	data, err := os.Open(options.File)
	if err != nil {
		log.Fatal(err)
	}
	client := &http.Client{}
	req, err := http.NewRequest("PUT", url, data)
	if err != nil {
		return err
	}
	req.SetBasicAuth(options.Username, options.Password)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	bodyText, err := ioutil.ReadAll(resp.Body)

	log.Infof("success %s", string(bodyText))

	return nil
}
