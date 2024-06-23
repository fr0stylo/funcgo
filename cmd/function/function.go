package main

import (
	"context"
	"io"
	"net/http"

	"go.uber.org/zap"

	"github.com/fr0stylo/funcgo/pkg/funcgo"
)

func main() {
	funcgo.Handler(func(ctx context.Context, log *zap.SugaredLogger, t *funcgo.Request) (*funcgo.Response, error) {
		log.Info(t.Body)

		res, err := http.Get("https://ifconfig.me")
		if err != nil {
			return &funcgo.Response{
				StatusCode: 503,
				Body:       err.Error(),
			}, err
		}
		buf, err := io.ReadAll(res.Body)
		if err != nil {
			return &funcgo.Response{
				StatusCode: 503,
				Body:       err.Error(),
			}, err
		}
		return &funcgo.Response{
			StatusCode: 206,
			Body:       string(buf),
		}, nil
	})
}

type Typed interface{}
