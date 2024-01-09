package metrics

import (
	"testing"

	"github.com/go-resty/resty/v2"

	metrics "github.com/quwan-sre/observability-go-contrib/metrics/resty"
)

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
}
