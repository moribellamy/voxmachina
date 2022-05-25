package main

import (
	"github.com/moribellamy/voxmachina/server"
	"github.com/moribellamy/voxmachina/storage"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"

	"github.com/moribellamy/voxmachina/utils"
)

func main() {
	utils.InfoLogger.Println("Running as PID", os.Getpid())
	var config utils.Config
	configBytes, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		utils.ErrorLogger.Fatal(err)
	}
	if err := yaml.Unmarshal(configBytes, &config); err != nil {
		utils.ErrorLogger.Fatal(err)
	}

	var store storage.Storage
	store, err = storage.FromConfig(config.Storage)
	if err != nil {
		utils.ErrorLogger.Fatal(err)
	}
	utils.InfoLogger.Println(store)

	vox, err := server.NewVox(config.Server, store)
	if err != nil {
		utils.ErrorLogger.Fatal(err)
	}
	utils.ErrorLogger.Fatal(vox.Run())
}
