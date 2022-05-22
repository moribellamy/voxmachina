package main

import (
	texttospeech "cloud.google.com/go/texttospeech/apiv1"
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/moribellamy/voxmachina/utils"
	"google.golang.org/api/option"
	texttospeechpb "google.golang.org/genproto/googleapis/cloud/texttospeech/v1"
	"google.golang.org/protobuf/encoding/protojson"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"net/http"
)

var grpcClient *texttospeech.Client

func v1Synthesize(c echo.Context) error {
	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return err
	}
	var req texttospeechpb.SynthesizeSpeechRequest
	if err := protojson.Unmarshal(body, &req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Could not convert argument to SynthesizeSpeechRequest: " + fmt.Sprint(err),
		})
	}

	resp, err := cachedGet(&req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": fmt.Sprint(err),
		})
	}
	respJson, err := protojson.Marshal(resp)
	if err != nil {
		return err
	}
	return c.JSONBlob(http.StatusOK, respJson)
}

func cachedGet(request *texttospeechpb.SynthesizeSpeechRequest) (*texttospeechpb.SynthesizeSpeechResponse, error) {
	ctx := context.Background()
	// TODO cache lookup
	return grpcClient.SynthesizeSpeech(ctx, request)
}

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
	grpcClient, err = texttospeech.NewClient(
		ctx,
		option.WithCredentialsFile(config.Credentials),
	)
	defer grpcClient.Close()
	if err != nil {
		utils.ErrorLogger.Fatal(err)
	}

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.POST("v1/text:synthesize", v1Synthesize)
	e.Logger.Fatal(e.Start(config.Hostport))
}
