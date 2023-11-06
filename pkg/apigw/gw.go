package apigw

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/fr0stylo/funcgo/pkg/funcgo"
	"github.com/fr0stylo/funcgo/pkg/runtime"
)

type Handler struct {
	Runner runtime.Runnable
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	b, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Failed to read body")
		return
	}

	req := funcgo.Request{
		Params:      mux.Vars(r),
		QueryParams: r.URL.Query(),
		Path:        r.URL.Path,
		Body:        string(b),
		Method:      r.Method,
	}

	res, err := h.Runner.Execute(req)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Failed to invoke function")
		return
	}

	var fres funcgo.Response
	json.Unmarshal(res, &fres)

	for n, h := range fres.Headers {
		w.Header().Add(n, h)
	}

	w.Header().Add("X-Upstream-Time", time.Since(start).String())
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(fres.StatusCode)
	fmt.Fprint(w, fres.Body)
}
