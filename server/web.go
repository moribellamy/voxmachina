package server

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	texttospeechpb "google.golang.org/genproto/googleapis/cloud/texttospeech/v1"
	"google.golang.org/protobuf/encoding/protojson"
)

type localWebserverWrapper struct {
	getter SynthesizeFunc
	server *echo.Echo
}

func NewLocalWebserver(getter SynthesizeFunc) *localWebserverWrapper {
	wrapper := localWebserverWrapper{}
	wrapper.getter = getter
	server := echo.New()
	server.Use(middleware.Logger())
	server.Use(middleware.Recover())
	server.POST("v1/text:synthesize", wrapper.v1Synthesize)
	wrapper.server = server
	return &wrapper
}

func setJsonMessage(c echo.Context, status int, messages ...any) error {
	return c.JSON(status, map[string]string{
		"message": fmt.Sprintln(messages...),
	})
}

func (wrapper *localWebserverWrapper) v1Synthesize(c echo.Context) error {
	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return err
	}
	var req texttospeechpb.SynthesizeSpeechRequest
	if err := protojson.Unmarshal(body, &req); err != nil {
		return setJsonMessage(
			c, http.StatusBadRequest,
			"Could not convert argument to SynthesizeSpeechRequest", err,
		)
	}

	resp, err := wrapper.getter(&req)
	if err != nil {
		return setJsonMessage(c, http.StatusInternalServerError, err)
	}
	respJson, err := protojson.Marshal(resp)
	if err != nil {
		return setJsonMessage(c, http.StatusInternalServerError, err)
	}
	return c.JSONBlob(http.StatusOK, respJson)
}

func (wrapper *localWebserverWrapper) Start(addr string) error {
	return wrapper.server.Start(addr)
}
