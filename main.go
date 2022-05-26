package main

import (
	"github.com/moribellamy/voxmachina/server"
	"github.com/moribellamy/voxmachina/utils"
)

func main() {
	if err := server.RunFromConfig("config.yaml"); err != nil {
		utils.ErrorLogger.Fatal(err)
	}
}
