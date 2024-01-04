package metrics

import (
	"flag"
	"fmt"
	metrics "github.com/quwan-sre/observability-go-contrib/metrics/grpc"
	"github.com/quwan-sre/observability-go-contrib/test/e2e/metrics/grpc_server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"log"
	"net/http"
	"testing"
)

func TestGRPCClient(t *testing.T) {
	flag.Parse()
	var opts []grpc.DialOption

	opts = append(
		opts,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(metrics.NewUnaryClientInterceptor()),
		grpc.WithStreamInterceptor(metrics.NewStreamClientInterceptor()),
	)

	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := grpc_server.NewRouteGuideClient(conn)

	// Looking for a valid feature
	printFeature(client, &grpc_server.Point{Latitude: 409146138, Longitude: -746188906})

	// Feature missing.
	printFeature(client, &grpc_server.Point{Latitude: 0, Longitude: 0})

	// Looking for features between 40, -75 and 42, -73.
	printFeatures(client, &grpc_server.Rectangle{
		Lo: &grpc_server.Point{Latitude: 400000000, Longitude: -750000000},
		Hi: &grpc_server.Point{Latitude: 420000000, Longitude: -730000000},
	})

	// RecordRoute
	runRecordRoute(client)

	// RouteChat
	runRouteChat(client)

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
