## etcd/client

etcd/client 使用 gRPC 通信，因此可以复用 gRPC 提供的 Interceptor 进行插装。

etcd Client 端插装代码示例：
```go
func initEtcdClient() {
	...
	cfg := clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
		DialOptions: []grpc.DialOption{
			grpc.WithBlock(),
			
			// Interceptor 作为 grpc.DialOption 使用
			grpc.WithUnaryInterceptor(metrics.NewUnaryClientInterceptor()),
			grpc.WithStreamInterceptor(metrics.NewStreamClientInterceptor()),
		},
	}
	var err error
	if etcdClient, err = clientv3.New(cfg); err != nil {
		panic(fmt.Sprintf("etcd client init err: %v", err))
	}
	...
}
```

更多指引可以参考 `grpc` 目录内容。