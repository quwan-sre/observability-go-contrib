package metrics

import (
	"net/http"
	"testing"
	"time"
)

func TestMain(t *testing.M) {
	// prepare a gin HTTP server
	go main()

	// health check
	for {
		resp, err := http.Get("http://127.0.0.1:8080/health")
		if err != nil || resp.StatusCode != 200 {
			time.Sleep(1 * time.Second)
		} else {
			break
		}
	}

	// ready, run the test case
	t.Run()
}
