package gin

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/go-resty/resty/v2"
)

var (
	client *resty.Client
)

func TestMain(t *testing.M) {
	// prepare a gin HTTP server
	go main()

	// health check
	for {
		resp, err := http.Get(fmt.Sprintf("http://127.0.0.1:%d/health", httpPort))
		if err != nil || resp.StatusCode != 200 {
			time.Sleep(1 * time.Second)
		} else {
			break
		}
	}

	for {
		resp, err := http.Get(fmt.Sprintf("http://127.0.0.1:%d/health", httpInstrumentPort))
		if err != nil || resp.StatusCode != 200 {
			time.Sleep(1 * time.Second)
		} else {
			break
		}
	}

	client = resty.New()

	// ready, run the test case
	t.Run()
}

func BenchmarkGin(b *testing.B) {
	for n := 0; n < b.N; n++ {
		client.R().Get(fmt.Sprintf("http://127.0.0.1:%d/health", httpPort))
	}

}

func BenchmarkGinWithInstrument(b *testing.B) {
	for n := 0; n < b.N; n++ {
		client.R().Get(fmt.Sprintf("http://127.0.0.1:%d/health", httpInstrumentPort))
	}
}
