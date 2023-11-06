package funcgo

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type Request struct {
	Params      map[string]string
	QueryParams map[string][]string
	Path        string
	Body        string
	Method      string
}

type Response struct {
	StatusCode int
	Body       string
	Headers    map[string]string
}

type HandlerFunc[T Request] func(context.Context, *zap.SugaredLogger, *Request) (*Response, error)

func Handler[T Request](h HandlerFunc[T]) {
	zz, _ := zap.NewProduction()
	z := zz.Sugar()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t := time.Now()
		var body Request
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "[%s] Incorrect Payload", r.URL.Path)
		}

		res, err := h(r.Context(), z, &body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, err.Error())
			return
		}

		if res.Headers == nil {
			res.Headers = make(map[string]string)
		}
		res.Headers["X-Internal-Time"] = time.Since(t).String()
		b, err := json.Marshal(res)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, err.Error())
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, string(b))
	})

	server := &http.Server{
		Addr:    "0.0.0.0:9999",
		Handler: handler,
	}

	z.Info("Listening [0.0.0.0:9999]...\n")
	z.Fatal(server.ListenAndServe(), "while listening")
}
