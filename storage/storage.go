package storage

import (
	"errors"

	"github.com/moribellamy/voxmachina/utils"
	texttospeechpb "google.golang.org/genproto/googleapis/cloud/texttospeech/v1"
)

type Storage interface {
	Store(*texttospeechpb.SynthesizeSpeechRequest, *texttospeechpb.SynthesizeSpeechResponse) error
	Get(*texttospeechpb.SynthesizeSpeechRequest) (*texttospeechpb.SynthesizeSpeechResponse, error)
	Close() error
}

var cannotFindStorage = errors.New("cannot find storage type")

func FromConfig(config utils.Storage) (Storage, error) {
	if config.Sqlite.Fpath != "" {
		return NewSqlite(config.Sqlite.Fpath)
	}
	if config.Lionrock.Hostport != "" {
		return NewLionrock(config.Lionrock.Hostport, config.Lionrock.Name, config.Lionrock.Prefix)
	}
	return nil, cannotFindStorage
}
