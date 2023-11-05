package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func main() {
	handle(func(ctx context.Context, log *zap.SugaredLogger, t *Request) error {
		log.Info(t.Body)
		time.Sleep(time.Second)
		return nil
	})
}

type Request struct {
	Path string
	Body string
}

type Typed interface{}

type Handler[T Request] func(context.Context, *zap.SugaredLogger, *Request) error

func handle[T Request](h Handler[T]) {
	zz, _ := zap.NewProduction()
	z := zz.Sugar()
	z.Info("Logging")

	h2s := &http2.Server{}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body Request
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			fmt.Fprintf(w, "[%s] Incorrect Payload", r.URL.Path)
			w.WriteHeader(http.StatusInternalServerError)
		}

		h(r.Context(), z, &body)

	})

	server := &http.Server{
		Addr:    "0.0.0.0:9999",
		Handler: h2c.NewHandler(handler, h2s),
	}

	z.Info("Listening [0.0.0.0:9999]...\n")
	z.Fatal(server.ListenAndServe(), "while listening")
}
