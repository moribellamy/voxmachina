package server

import (
	"context"
	"errors"
	"net"

	texttospeechpb "google.golang.org/genproto/googleapis/cloud/texttospeech/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type SynthesizeFunc func(*texttospeechpb.SynthesizeSpeechRequest) (
	*texttospeechpb.SynthesizeSpeechResponse, error)

type localGrpcServerWrapper struct {
	getter SynthesizeFunc
	server *grpc.Server
}

func newLocalGrpcServer(getter SynthesizeFunc) *localGrpcServerWrapper {
	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)
	return &localGrpcServerWrapper{
		getter,
		grpcServer,
	}
}

func (wrapper *localGrpcServerWrapper) SynthesizeSpeech(
	_ context.Context, req *texttospeechpb.SynthesizeSpeechRequest) (
	*texttospeechpb.SynthesizeSpeechResponse, error) {
	resp, err := wrapper.getter(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (*localGrpcServerWrapper) ListVoices(
	_ context.Context, _ *texttospeechpb.ListVoicesRequest) (
	*texttospeechpb.ListVoicesResponse, error) {
	return nil, errors.New("not implemented")
}

func (wrapper *localGrpcServerWrapper) Start(addr string) error {
	texttospeechpb.RegisterTextToSpeechServer(wrapper.server, wrapper)
	if conn, err := net.Listen("tcp", addr); err != nil {
		return err
	} else {
		return wrapper.server.Serve(conn)
	}
}

func (wrapper *localGrpcServerWrapper) GracefulStop() {
	wrapper.server.GracefulStop()
}
