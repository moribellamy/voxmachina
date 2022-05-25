package storage

import (
	"database/sql"
	"errors"

	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/multierr"
	texttospeechpb "google.golang.org/genproto/googleapis/cloud/texttospeech/v1"
	"google.golang.org/protobuf/proto"
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
	_, err = sqlite.db.Exec("create table if not exists cache (" +
		"req blob not null primary key, " +
		"resp blob not null);")
	return &sqlite, err
}

func (sqlite *Sqlite) Close() error {
	return sqlite.db.Close()
}

func (sqlite *Sqlite) Store(
	request *texttospeechpb.SynthesizeSpeechRequest,
	response *texttospeechpb.SynthesizeSpeechResponse) error {
	reqText, reqErr := proto.Marshal(request)
	respText, respErr := proto.Marshal(response)
	if reqErr != nil || respErr != nil {
		return multierr.Combine(reqErr, respErr)
	}

	_, err := sqlite.db.Exec("insert into cache(req, resp) values (?, ?)", reqText, respText)
	return err
}

var cacheMiss = errors.New("cache miss")

func (sqlite *Sqlite) Get(request *texttospeechpb.SynthesizeSpeechRequest) (
	*texttospeechpb.SynthesizeSpeechResponse, error) {
	reqText, err := proto.Marshal(request)
	if err != nil {
		return nil, err
	}
	rows, readErr := sqlite.db.Query("select resp from cache where req = ?", reqText)
	if readErr != nil {
		return nil, readErr
	}
	defer rows.Close()
	if !rows.Next() {
		return nil, cacheMiss
	}
	var respText []byte
	if err = rows.Scan(&respText); err != nil {
		return nil, err
	}
	resp := texttospeechpb.SynthesizeSpeechResponse{}
	err = proto.Unmarshal(respText, &resp)
	return &resp, err
}
