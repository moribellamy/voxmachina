package storage

import (
	"context"
	lionrockpb "github.com/moribellamy/voxmachina/proto/lionrock"
	"go.uber.org/multierr"
	texttospeechpb "google.golang.org/genproto/googleapis/cloud/texttospeech/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/proto"
)

type lionrockWrapper struct {
	conn   *grpc.ClientConn
	client lionrockpb.TransactionalKeyValueStoreClient
	ctx    context.Context
	name   string
	prefix string
}

func NewLionrock(addr string, name string, prefix string) (*lionrockWrapper, error) {
	var err error
	lionrock := &lionrockWrapper{}
	lionrock.prefix = prefix
	lionrock.name = name
	lionrock.ctx = context.Background()
	lionrock.conn, err = grpc.Dial(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}
	lionrock.client = lionrockpb.NewTransactionalKeyValueStoreClient(lionrock.conn)
	return lionrock, nil
}

func (lionrock *lionrockWrapper) Close() error {
	return lionrock.conn.Close()
}

func (lionrock *lionrockWrapper) Store(
	request *texttospeechpb.SynthesizeSpeechRequest,
	response *texttospeechpb.SynthesizeSpeechResponse) error {
	requestBytes, requestErr := proto.Marshal(request)
	responseBytes, responseErr := proto.Marshal(response)
	if requestErr != nil || responseErr != nil {
		return multierr.Combine(requestErr, responseErr)
	}
	key := append([]byte(lionrock.prefix), requestBytes...)

	req := lionrockpb.DatabaseRequest{
		DatabaseName: lionrock.name,
		Request: &lionrockpb.DatabaseRequest_SetValue{
			SetValue: &lionrockpb.SetValueRequest{
				Key:   key,
				Value: responseBytes,
			},
		},
	}
	_, err := lionrock.client.Execute(lionrock.ctx, &req)
	if err != nil {
		return err
	}
	return nil
}

func (lionrock *lionrockWrapper) Get(request *texttospeechpb.SynthesizeSpeechRequest) (
	*texttospeechpb.SynthesizeSpeechResponse, error) {
	requestBytes, requestErr := proto.Marshal(request)
	if requestErr != nil {
		return nil, requestErr
	}
	key := append([]byte(lionrock.prefix), requestBytes...)
	req := lionrockpb.DatabaseRequest{
		DatabaseName: lionrock.name,
		Request: &lionrockpb.DatabaseRequest_GetValue{
			GetValue: &lionrockpb.GetValueRequest{
				Key: key,
			},
		},
	}
	resp, err := lionrock.client.Execute(lionrock.ctx, &req)
	if err != nil {
		return nil, nil
	}
	respBytes := resp.GetGetValue().Value
	speechResp := texttospeechpb.SynthesizeSpeechResponse{}
	if err = proto.Unmarshal(respBytes, &speechResp); err != nil {
		return nil, err
	}

	return &speechResp, nil
}
