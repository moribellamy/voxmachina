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
	_, dropErr := sqlite.db.Exec("drop table if exists cache;")
	_, createErr := sqlite.db.Exec("create table cache (" +
		"req blob not null primary key, " +
		"resp blob not null);")
	return &sqlite, multierr.Combine(dropErr, createErr)
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

var cacheMiss error = errors.New("cache miss")

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
	rows.Next()
	var respText []byte
	if err = rows.Scan(&respText); err != nil {
		return nil, err
	}
	if respText == nil {
		return nil, cacheMiss
	}
	resp := texttospeechpb.SynthesizeSpeechResponse{}
	err = proto.Unmarshal(respText, &resp)
	return &resp, err
}
