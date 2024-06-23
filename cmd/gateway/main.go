package gateway

import (
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"

	"github.com/fr0stylo/funcgo/cmd/gateway/handlers/admin"
)

func main() {
	m := mux.NewRouter()
	s := m.PathPrefix("/admin").Subrouter()
	s.Handle("/function", &admin.CreateFunction{M: m}).Methods(http.MethodPost)
	s.Handle("/function", &admin.CreateFunction{M: m}).Methods(http.MethodDelete)

	zap.S().Info("Listening on :8000")
	http.ListenAndServe("0.0.0.0:8000", m)
}
