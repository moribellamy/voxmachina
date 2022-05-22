package main

import (
	texttospeech "cloud.google.com/go/texttospeech/apiv1"
	"context"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/moribellamy/voxmachina/utils"
	"google.golang.org/api/option"
	texttospeechpb "google.golang.org/genproto/googleapis/cloud/texttospeech/v1"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"net"
	"net/http"
)

var grpcClient *texttospeech.Client

func setJsonMessage(c echo.Context, status int, messages ...any) error {
	return c.JSON(status, map[string]string{
		"message": fmt.Sprintln(messages...),
	})
}

func v1Synthesize(c echo.Context) error {
	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return err
	}
	var req texttospeechpb.SynthesizeSpeechRequest
	if err := protojson.Unmarshal(body, &req); err != nil {
		return setJsonMessage(
			c, http.StatusBadRequest,
			"Could not convert argument to SynthesizeSpeechRequest", err,
		)
	}

	resp, err := cachedGet(&req)
	if err != nil {
		return setJsonMessage(c, http.StatusInternalServerError, err)
	}
	respJson, err := protojson.Marshal(resp)
	if err != nil {
		return setJsonMessage(c, http.StatusInternalServerError, err)
	}
	return c.JSONBlob(http.StatusOK, respJson)
}

func cachedGet(request *texttospeechpb.SynthesizeSpeechRequest) (*texttospeechpb.SynthesizeSpeechResponse, error) {
	ctx := context.Background()
	return grpcClient.SynthesizeSpeech(ctx, request)
}

func runWebServer(config utils.Config, c chan error) {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.POST("v1/text:synthesize", v1Synthesize)
	c <- e.Start(config.Cache.WebHostport)
}

type cachingTextToSpeechServer struct {
}

func (s *cachingTextToSpeechServer) SynthesizeSpeech(
	ctx context.Context, req *texttospeechpb.SynthesizeSpeechRequest) (
	*texttospeechpb.SynthesizeSpeechResponse, error) {
	resp, err := cachedGet(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *cachingTextToSpeechServer) ListVoices(
	ctx context.Context, req *texttospeechpb.ListVoicesRequest) (
	*texttospeechpb.ListVoicesResponse, error) {
	return nil, errors.New("not implemented")
}

func runGrpcServer(config utils.Config, c chan error) {
	grpcServer := grpc.NewServer()
	cachingServer := cachingTextToSpeechServer{}
	texttospeechpb.RegisterTextToSpeechServer(grpcServer, &cachingServer)

	conn, err := net.Listen("tcp", config.Cache.GrpcHostport)
	if err != nil {
		c <- err
	} else {
		c <- grpcServer.Serve(conn)
	}
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

	executor := make(chan error)
	go runWebServer(config, executor)
	go runGrpcServer(config, executor)

	err = <-executor
	utils.ErrorLogger.Println("First server terminated:", err)
	err = <-executor
	utils.ErrorLogger.Println("Second server terminated:", err)
}
