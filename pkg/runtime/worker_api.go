package runtime

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net"
	"net/http"

	"golang.org/x/net/http2"
)

type WorkerApi struct {
	ip     string
	client *http.Client
}

func NewWorkerApi(ip string) *WorkerApi {
	return &WorkerApi{
		ip: ip,
		client: &http.Client{
			Transport: &http2.Transport{
				AllowHTTP: true,
				DialTLSContext: func(ctx context.Context, network, addr string, cfg *tls.Config) (net.Conn, error) {
					var d net.Dialer
					return d.DialContext(ctx, network, addr)
				},
			},
		},
	}
}

func (r *WorkerApi) Execute(o any) (any, error) {
	b, _ := json.Marshal(o)

	res, err := r.client.Post(fmt.Sprintf("http://%s:9999", r.ip), "application/json", bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	// Kills other types as it becomes map[string]interface
	var body any
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		return nil, err
	}

	return body, nil
}
