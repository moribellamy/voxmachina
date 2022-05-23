package main

import (
	"context"
	"io/ioutil"

	"github.com/moribellamy/voxmachina/server"

	texttospeech "cloud.google.com/go/texttospeech/apiv1"
	"google.golang.org/api/option"
	texttospeechpb "google.golang.org/genproto/googleapis/cloud/texttospeech/v1"
	"gopkg.in/yaml.v3"

	"github.com/moribellamy/voxmachina/utils"
)

func main() {
	var config utils.Config
	configBytes, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		utils.ErrorLogger.Fatal(err)
	}
	if err := yaml.Unmarshal(configBytes, &config); err != nil {
		utils.ErrorLogger.Fatal(err)
	}

	ctx := context.Background()
	var grpcClient *texttospeech.Client
	grpcClient, err = texttospeech.NewClient(
		ctx,
		option.WithCredentialsFile(config.Credentials),
	)
	defer grpcClient.Close()

	cachedGet := func(request *texttospeechpb.SynthesizeSpeechRequest) (
		*texttospeechpb.SynthesizeSpeechResponse, error) {
		ctx := context.Background()
		return grpcClient.SynthesizeSpeech(ctx, request)
	}

	if err != nil {
		utils.ErrorLogger.Fatal(err)
	}

	executor := make(chan error)
	webserver := server.NewLocalWebserver(cachedGet)
	go func() {
		executor <- webserver.Start(config.Cache.WebHostport)
	}()
	grpcServer := server.NewCachingTextToSpeechServer(cachedGet)
	go func() {
		executor <- grpcServer.Start(config.Cache.GrpcHostport)
	}()

	err = <-executor
	utils.ErrorLogger.Println("First server terminated:", err)
	err = <-executor
	utils.ErrorLogger.Println("Second server terminated:", err)
}
