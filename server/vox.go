package server

import (
	"context"
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"

	texttospeech "cloud.google.com/go/texttospeech/apiv1"
	"github.com/moribellamy/voxmachina/storage"
	"github.com/moribellamy/voxmachina/utils"
	"go.uber.org/multierr"
	"google.golang.org/api/option"
	texttospeechpb "google.golang.org/genproto/googleapis/cloud/texttospeech/v1"
)

type Vox struct {
	config     utils.Server
	grpc       *localGrpcServerWrapper
	web        *localWebserverWrapper
	grpcClient *texttospeech.Client
	ctx        context.Context
	store      storage.Storage
}

func NewVox(config utils.Server, store storage.Storage) (*Vox, error) {
	var err error
	vox := &Vox{}
	vox.ctx = context.Background()
	vox.store = store
	vox.config = config

	vox.grpcClient, err = texttospeech.NewClient(
		vox.ctx,
		option.WithCredentialsFile(config.Credentials),
	)
	if err != nil {
		return nil, err
	}

	vox.grpc = newLocalGrpcServer(vox.upstreamGet)
	vox.web = newLocalWebserver(vox.upstreamGet)
	return vox, nil
}

func (vox *Vox) upstreamGet(request *texttospeechpb.SynthesizeSpeechRequest) (
	*texttospeechpb.SynthesizeSpeechResponse, error) {
	resp, err := vox.store.Get(request)
	// Cache hit
	if err == nil {
		return resp, nil
	}
	// Cache miss
	resp, err = vox.grpcClient.SynthesizeSpeech(vox.ctx, request)
	if err != nil {
		return nil, err
	}
	err = vox.store.Store(request, resp)
	return resp, err
}

func (vox *Vox) Run() error {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	servers := make(chan error, 2)
	go func() {
		servers <- vox.web.Start(vox.config.WebHostport)
	}()
	go func() {
		servers <- vox.grpc.Start(vox.config.GrpcHostport)
	}()

	var servingErr error
	select {
	case servingErr = <-servers:
	case sig := <-sigs:
		servingErr = errors.New(fmt.Sprintln("Caught signal", sig))
	}
	shutdownErr := vox.GracefulShutdown()
	return multierr.Combine(servingErr, shutdownErr)
}

func (vox *Vox) GracefulShutdown() error {
	vox.grpc.GracefulStop()
	return multierr.Combine(
		vox.web.GracefulStop(),
		vox.store.Close(),
	)
}

func RunFromConfig(fpath string) error {
	utils.InfoLogger.Println("Running as PID", os.Getpid())
	var config utils.Config
	configBytes, err := ioutil.ReadFile(fpath)
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

	vox, err := NewVox(config.Server, store)
	if err != nil {
		utils.ErrorLogger.Fatal(err)
	}
	return vox.Run()
}
