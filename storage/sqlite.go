package storage

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	texttospeechpb "google.golang.org/genproto/googleapis/cloud/texttospeech/v1"
)

type Sqlite struct {
	db *sql.DB
}

func NewSqlite(fpath string) (*Sqlite, error) {
	var err error
	sqlite := Sqlite{}
	sqlite.db, err = sql.Open("sqlite3", fpath)
	if err != nil {
		return nil, err
	}
	return &sqlite, nil
}

func (sqlite *Sqlite) Store(
	request *texttospeechpb.SynthesizeSpeechRequest,
	response *texttospeechpb.SynthesizeSpeechResponse) error {
	return nil
}
