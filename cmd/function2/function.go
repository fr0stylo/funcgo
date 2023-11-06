package main

import (
	"context"
	"encoding/json"

	"go.uber.org/zap"

	"github.com/fr0stylo/funcgo/pkg/funcgo"
)

func main() {
	funcgo.Handler(func(ctx context.Context, log *zap.SugaredLogger, t *funcgo.Request) (*funcgo.Response, error) {
		log.Info(t.Body)
		b, _ := json.Marshal(t)
		return &funcgo.Response{
			StatusCode: 201,
			Body:       string(b),
		}, nil
	})
}

type Typed interface{}
