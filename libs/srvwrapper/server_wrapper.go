package srvwrapper

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
)

type Wrapper[Request any, Response any] struct {
	handleFn func(ctx context.Context, req Request) (Response, error)
}

func New[Request any, Response any](
	fn func(ctx context.Context, req Request) (Response, error),
) *Wrapper[Request, Response] {
	return &Wrapper[Request, Response]{
		handleFn: fn,
	}
}

func (w *Wrapper[Request, Response]) ServeHTTP(resWriter http.ResponseWriter, httpReq *http.Request) {
	var req Request

	err := json.NewDecoder(httpReq.Body).Decode(&req)
	if err != nil {
		log.Printf("%s", err.Error())
		resWriter.WriteHeader(http.StatusInternalServerError)
		writeErrorText(resWriter, "parse request", err)
		return
	}

	res, err := w.handleFn(httpReq.Context(), req)
	if err != nil {
		log.Printf("%s", err.Error())
		resWriter.WriteHeader(http.StatusInternalServerError)
		writeErrorText(resWriter, "exec handler", err)
		return
	}

	resBytes, err := json.Marshal(res)
	if err != nil {
		log.Printf("%s", err.Error())
		resWriter.WriteHeader(http.StatusInternalServerError)
		writeErrorText(resWriter, "decode response", err)
		return
	}

	_, _ = resWriter.Write(resBytes)
}

func writeErrorText(w http.ResponseWriter, text string, err error) {
	buf := bytes.NewBufferString(text)

	buf.WriteString(": ")
	buf.WriteString(err.Error())
	buf.WriteByte('\n')

	w.Write(buf.Bytes())
}
