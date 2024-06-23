package admin

import (
	"net/http"

	"github.com/gorilla/mux"
)

type CreateFunction struct {
	M *mux.Router
}

func (r *CreateFunction) ServeHTTP(w http.ResponseWriter, req *http.Request) {

}
