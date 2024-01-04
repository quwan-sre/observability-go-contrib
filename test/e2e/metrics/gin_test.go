package metrics

import (
	"fmt"
	"io"
	"net/http"
	"testing"
)

func TestExpose(t *testing.T) {
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

func TestNormalRequest(t *testing.T) {
	testCases := []string{
		"http://127.0.0.1:8080/exist",
		"http://127.0.0.1:8080/not_exist",
	}

	for _, tc := range testCases {
		http.Get(tc)
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
