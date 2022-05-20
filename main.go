// Command quickstart generates an audio file with the content "Hello, World!".
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"google.golang.org/api/option"
	"gopkg.in/ini.v1"
	"io/ioutil"
	"log"
	"net/http"

	texttospeech "cloud.google.com/go/texttospeech/apiv1"
	texttospeechpb "google.golang.org/genproto/googleapis/cloud/texttospeech/v1"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"google.golang.org/protobuf/encoding/protojson"
)

var grpcClient *texttospeech.Client

func myLog(request *texttospeechpb.SynthesizeSpeechRequest) {
	asJson, err := protojson.Marshal(request)
	if err != nil {
		log.Fatal(err)
	}
	var indented bytes.Buffer
	if err = json.Indent(&indented, asJson, "", "  "); err != nil {
		log.Fatal(err)
	}
	fmt.Printf(indented.String())
}

func v1Synthesize(c echo.Context) error {
	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return echo.ErrInternalServerError
	}
	var req texttospeechpb.SynthesizeSpeechRequest
	if err := protojson.Unmarshal(body, &req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Could not convert argument to SynthesizeSpeechRequest",
			"error":   fmt.Sprint(err),
		})
	}

	myLog(&req)

	ctx := context.Background()
	resp, err := grpcClient.SynthesizeSpeech(ctx, &req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Upstream Google TTS Server responded with an error.",
			"error":   fmt.Sprint(err),
		})
	}
	respJson, err := protojson.Marshal(resp)
	if err != nil {
		return echo.ErrInternalServerError
	}
	return c.JSONBlob(http.StatusOK, respJson)
}

func main() {
	config, err := ini.Load("config.ini")
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	grpcClient, err = texttospeech.NewClient(
		ctx,
		option.WithCredentialsFile(config.Section("tts").Key("credentialsFile").String()),
	)
	defer grpcClient.Close()
	if err != nil {
		log.Fatal(err)
	}

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.POST("v1/text:synthesize", v1Synthesize)
	e.Logger.Fatal(e.Start(config.Section("voxmachina").Key("hostport").String()))
}
