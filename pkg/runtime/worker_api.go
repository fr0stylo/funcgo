package runtime

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type WorkerApi struct {
	ip     string
	client *http.Client
}

func NewWorkerApi(ip string) *WorkerApi {
	return &WorkerApi{
		ip: ip,
		client: &http.Client{
			Transport: &http.Transport{},
		},
	}
}

func (r *WorkerApi) Execute(o any) ([]byte, error) {
	b, _ := json.Marshal(o)

	res, err := r.client.Post(fmt.Sprintf("http://%s:9999", r.ip), "application/json", bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	// Kills other types as it becomes map[string]interface
	return io.ReadAll(res.Body)
}
