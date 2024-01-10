package metrics

import (
	"context"
	"fmt"
	"testing"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"

	metrics "github.com/quwan-sre/observability-go-contrib/metrics/grpc"
)

var (
	etcdClient *clientv3.Client
)

func initEtcdClient() {
	cfg := clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
		DialOptions: []grpc.DialOption{
			grpc.WithBlock(),
			grpc.WithUnaryInterceptor(metrics.NewUnaryClientInterceptor()),
			grpc.WithStreamInterceptor(metrics.NewStreamClientInterceptor()),
		},
	}
	var err error
	if etcdClient, err = clientv3.New(cfg); err != nil {
		panic(fmt.Sprintf("etcd client init err: %v", err))
	}
	return
}

func TestEtcdClient(t *testing.T) {
	initEtcdClient()
	etcdClient.Get(context.TODO(), "foo")
	etcdClient.Delete(context.TODO(), "asdf asdf")
}
