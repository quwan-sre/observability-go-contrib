package resty

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	metrics "github.com/quwan-sre/observability-go-contrib/metrics/resty"
	"io"
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

func TestRestyClient(t *testing.T) {
	restyClient := resty.New()
	metrics.NewMetricsMiddleware(restyClient)

	testCases := []string{
		"http://127.0.0.1:8080/health",
		"http://127.0.0.1:8080/not_exist",
	}

	for _, tc := range testCases {
		for i := 0; i < 10; i++ {
			restyClient.R().Get(tc)
		}
	}

	resp, err := http.Get("http://127.0.0.1:8080/metrics")
	if err != nil || resp.StatusCode != 200 {
		t.Fatalf("test failed, err: %v", err)
	}
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("test failed read body, err: %v", err)
	}
	fmt.Println(string(bodyBytes))
}
