package storage

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	texttospeechpb "google.golang.org/genproto/googleapis/cloud/texttospeech/v1"
)

type SqlLite struct {
	db *sql.DB
}

func NewSqlLite(fpath string) (*SqlLite, error) {
	var err error
	sqllite := SqlLite{}
	sqllite.db, err = sql.Open("sqlite3", fpath)
	if err != nil {
		return nil, err
	}
	return &sqllite, nil
}

func (sqllite *SqlLite) Store(
	request *texttospeechpb.SynthesizeSpeechRequest,
	response *texttospeechpb.SynthesizeSpeechResponse) error {
	return nil
}
