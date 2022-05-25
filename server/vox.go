package server

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	texttospeech "cloud.google.com/go/texttospeech/apiv1"
	"github.com/moribellamy/voxmachina/utils"
	"go.uber.org/multierr"
	"google.golang.org/api/option"
	texttospeechpb "google.golang.org/genproto/googleapis/cloud/texttospeech/v1"
)

type Vox struct {
	config     utils.Config
	grpc       *localGrpcServerWrapper
	web        *localWebserverWrapper
	grpcClient *texttospeech.Client
	ctx        context.Context
}

func NewVox(config utils.Config) (*Vox, error) {
	var err error
	vox := &Vox{}
	vox.ctx = context.Background()
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
	return vox.grpcClient.SynthesizeSpeech(vox.ctx, request)
}

func (vox *Vox) Run() error {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	servers := make(chan error, 2)
	go func() {
		servers <- vox.web.Start(vox.config.Cache.WebHostport)
	}()
	go func() {
		servers <- vox.grpc.Start(vox.config.Cache.GrpcHostport)
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
	return vox.web.GracefulStop()
}
