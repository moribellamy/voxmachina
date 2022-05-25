package utils

import (
	"bytes"
	"encoding/json"
	"log"
	"os"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

var ErrorLogger = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
var InfoLogger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)

func AsJson(request proto.Message) string {
	asJson, err := protojson.Marshal(request)
	if err != nil {
		log.Fatal(err)
	}
	var indented bytes.Buffer
	if err = json.Indent(&indented, asJson, "", "  "); err != nil {
		log.Fatal(err)
	}
	return indented.String()
}
