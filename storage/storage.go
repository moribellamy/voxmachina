package storage

import (
	"errors"
	"github.com/moribellamy/voxmachina/utils"
	texttospeechpb "google.golang.org/genproto/googleapis/cloud/texttospeech/v1"
)

type Storage interface {
	Store(*texttospeechpb.SynthesizeSpeechRequest, *texttospeechpb.SynthesizeSpeechResponse) error
}

func FromConfig(config utils.Storage) (Storage, error) {
	if config.Sqllite.Fpath != "" {
		sqllite, err := NewSqlLite(config.Sqllite.Fpath)
		return sqllite, err
	}
	return nil, errors.New("cannot find storage type")
}
