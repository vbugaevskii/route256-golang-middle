package cliwrapper

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type RequestPreparer interface {
	Prepare(ctx context.Context, netloc string) (*http.Request, error)
}

type Wrapper[Request RequestPreparer, Response any] struct {
	netloc string
}

func New[Request RequestPreparer, Response any](netloc string) *Wrapper[Request, Response] {
	return &Wrapper[Request, Response]{
		netloc: netloc,
	}
}

func (w *Wrapper[Request, Response]) Retrieve(ctx context.Context, req Request) (Response, error) {
	var res Response

	reqHttp, err := req.Prepare(ctx, w.netloc)
	if err != nil {
		return res, err
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
