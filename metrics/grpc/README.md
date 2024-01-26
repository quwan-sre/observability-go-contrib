## gRPC-Go

gRPC-Go 对外提供了 Interceptor 机制，用于挂载 Interceptor 扩展额外的请求响应逻辑。

```go
type UnaryClientInterceptor func(ctx context.Context, method string, req, reply any, cc *ClientConn, invoker UnaryInvoker, opts ...CallOption) error
type UnaryServerInterceptor func(ctx context.Context, req any, info *UnaryServerInfo, handler UnaryHandler) (resp any, err error)
type StreamClientInterceptor func(ctx context.Context, desc *StreamDesc, cc *ClientConn, method string, streamer Streamer, opts ...CallOption) (ClientStream, error)
type StreamServerInterceptor func(srv any, ss ServerStream, info *StreamServerInfo, handler StreamHandler) error
```

Interceptor 可以作为 Client 或 Server 启动时的 Options 传入。因此，Server 端插装代码示例：
```go
func RunGRPCServer() {
	...
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	
	// Interceptor 作为 grpc.ServerOption 使用
	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(metrics.NewUnaryServerInterceptor()),
		grpc.StreamInterceptor(metrics.NewStreamServerInterceptor()),
	}
	
	grpcServer := grpc.NewServer(opts...)
	...
}
```

Client 端插装代码示例：
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

gRPC 支持 Unary 和 Stream 类型的 RPC 调用，如果其中的部分类型没有使用，也可以不添加对应的 Interceptor。