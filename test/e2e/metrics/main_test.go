package metrics

import (
	"net"
	"net/http"
	"testing"
	"time"
)

func TestMain(t *testing.M) {
	// prepare a gin HTTP server
	main()

	// Gin health check
	for {
		resp, err := http.Get("http://127.0.0.1:8080/health")
		if err != nil || resp.StatusCode != 200 {
			time.Sleep(1 * time.Second)
		} else {
			break
		}
	}

	// gRPC health check
	for {
		timeout := time.Second
		_, err := net.DialTimeout("tcp", net.JoinHostPort("127.0.0.1", "8081"), timeout)
		if err != nil {
			time.Sleep(1 * time.Second)
		} else {
			break
		}
	}

	// ready, run the test case
	t.Run()
}
