package metrics

import (
	"flag"
	"log"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	metrics "github.com/quwan-sre/observability-go-contrib/metrics/grpc"
	"github.com/quwan-sre/observability-go-contrib/test/e2e/metrics/pb"
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
	client := pb.NewRouteGuideClient(conn)

	// Looking for a valid feature
	printFeature(client, &pb.Point{Latitude: 409146138, Longitude: -746188906})

	// Feature missing.
	printFeature(client, &pb.Point{Latitude: 0, Longitude: 0})

	// Looking for features between 40, -75 and 42, -73.
	printFeatures(client, &pb.Rectangle{
		Lo: &pb.Point{Latitude: 400000000, Longitude: -750000000},
		Hi: &pb.Point{Latitude: 420000000, Longitude: -730000000},
	})

	// RecordRoute
	runRecordRoute(client)

	// RouteChat
	runRouteChat(client)
}
