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

func FromConfig(config utils.Storage) (Storage, error) {
	if config.Sqlite.Fpath != "" {
		sqllite, err := NewSqlite(config.Sqlite.Fpath)
		return sqllite, err
	}
	return nil, errors.New("cannot find storage type")
}
