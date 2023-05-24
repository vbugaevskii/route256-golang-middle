package cliwrapper

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type Wrapper[Request any, Response any] struct {
	netloc   string
	endpoint string
	method   string
}

func New[Request any, Response any](netloc string, endpoint string, method string) *Wrapper[Request, Response] {
	return &Wrapper[Request, Response]{
		netloc:   netloc,
		endpoint: endpoint,
		method:   method,
	}
}

func (w *Wrapper[Request, Response]) Retrieve(ctx context.Context, req Request) (Response, error) {
	var res Response

	reqBytes, err := json.Marshal(&req)
	if err != nil {
		return res, fmt.Errorf("encode request: %w", err)
	}

	reqHttp, err := http.NewRequestWithContext(ctx, w.method, w.netloc+w.endpoint, bytes.NewBuffer(reqBytes))
	if err != nil {
		return res, fmt.Errorf("prepare request: %w", err)
	}

	resHttp, err := http.DefaultClient.Do(reqHttp)
	if err != nil {
		return res, fmt.Errorf("do request: %w", err)
	}
	defer resHttp.Body.Close()

	if resHttp.StatusCode != http.StatusOK {
		return res, fmt.Errorf("wrong status code: %d", resHttp.StatusCode)
	}

	err = json.NewDecoder(resHttp.Body).Decode(&res)
	if err != nil {
		return res, fmt.Errorf("decode request: %w", err)
	}

	return res, nil
}
