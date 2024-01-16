package grpc

import (
	"context"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/quwan-sre/observability-go-contrib/metrics/common"
)

func NewUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		startTime := time.Now()

		defer func() {
			latency := time.Now().Sub(startTime)
			code := codes.OK

			if err != nil {
				s, _ := status.FromError(err)
				code = s.Code()
			}

			requestPath, requestTarget := ParseFullMethod(info.FullMethod)
			common.DefaultRPCReceiveRequestMetric.With(prometheus.Labels{
				"sdk":                  common.RPCSDKGRPC,
				"request_protocol":     common.RPCProtocolGRPC,
				"request_target":       requestTarget,
				"request_path":         requestPath,
				"grpc_response_status": strconv.Itoa(int(code)),
				"response_code":        "0",
			}).Observe(latency.Seconds() * 1000)
		}()

		resp, err = handler(ctx, req)

		return resp, err
	}
}

func NewStreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		requestPath, requestTarget := ParseFullMethod(info.FullMethod)
		err = handler(srv, &wrapServerStream{ss, requestPath, requestTarget})
		return err
	}
}

func NewUnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) (err error) {
		startTime := time.Now()

		defer func() {
			latency := time.Now().Sub(startTime)
			code := codes.OK

			if err != nil {
				s, _ := status.FromError(err)
				code = s.Code()
			}

			requestPath, requestTarget := ParseFullMethod(method)
			common.DefaultRPCSendRequestMetric.With(prometheus.Labels{
				"sdk":                  common.RPCSDKGRPC,
				"request_protocol":     common.RPCProtocolGRPC,
				"request_target":       requestTarget,
				"request_path":         requestPath,
				"grpc_response_status": strconv.Itoa(int(code)),
				"response_code":        "0",
			}).Observe(latency.Seconds() * 1000)
		}()

		err = invoker(ctx, method, req, reply, cc, opts...)
		return err
	}
}

func NewStreamClientInterceptor() grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		s, err := streamer(ctx, desc, cc, method, opts...)
		requestPath, requestTarget := ParseFullMethod(method)
		return &wrapClientStream{
			ClientStream:  s,
			requestTarget: requestTarget,
			requestPath:   requestPath,
		}, err
	}
}

// wrapClientStream wraps grpc.ClientStream to record each Sent/Recv of message in histogram.
type wrapClientStream struct {
	grpc.ClientStream
	requestPath   string
	requestTarget string
}

func (w *wrapClientStream) SendMsg(m interface{}) error {
	startTime := time.Now()
	var err error

	defer func() {
		latency := time.Now().Sub(startTime)
		code := codes.OK

		if err != nil && err != io.EOF {
			s, _ := status.FromError(err)
			code = s.Code()
		}

		common.DefaultRPCSendRequestMetric.With(prometheus.Labels{
			"sdk":                  common.RPCSDKGRPC,
			"request_protocol":     common.RPCProtocolGRPC,
			"request_target":       w.requestTarget,
			"request_path":         w.requestPath,
			"grpc_response_status": strconv.Itoa(int(code)),
			"response_code":        "0",
		}).Observe(latency.Seconds() * 1000)
	}()

	err = w.ClientStream.SendMsg(m)
	return err
}

func (w *wrapClientStream) RecvMsg(m interface{}) error {
	var err error
	startTime := time.Now()

	defer func() {
		latency := time.Now().Sub(startTime)
		code := codes.OK

		if err != nil && err != io.EOF {
			s, _ := status.FromError(err)
			code = s.Code()
		}

		common.DefaultRPCReceiveRequestMetric.With(prometheus.Labels{
			"sdk":                  common.RPCSDKGRPC,
			"request_protocol":     common.RPCProtocolGRPC,
			"request_target":       w.requestTarget,
			"request_path":         w.requestPath,
			"grpc_response_status": strconv.Itoa(int(code)),
			"response_code":        "0",
		}).Observe(latency.Seconds() * 1000)
	}()

	err = w.ClientStream.RecvMsg(m)
	return err
}

// wrapServerStream wraps grpc.ServerStream to record each Sent/Recv of message in histogram.
type wrapServerStream struct {
	grpc.ServerStream
	requestPath   string
	requestTarget string
}

func (w *wrapServerStream) SendMsg(m interface{}) error {
	startTime := time.Now()
	var err error

	defer func() {
		latency := time.Now().Sub(startTime)
		code := codes.OK

		if err != nil && err != io.EOF {
			s, _ := status.FromError(err)
			code = s.Code()
		}

		common.DefaultRPCSendRequestMetric.With(prometheus.Labels{
			"sdk":                  common.RPCSDKGRPC,
			"request_protocol":     common.RPCProtocolGRPC,
			"request_target":       w.requestTarget,
			"request_path":         w.requestPath,
			"grpc_response_status": strconv.Itoa(int(code)),
			"response_code":        "0",
		}).Observe(latency.Seconds() * 1000)
	}()

	err = w.ServerStream.SendMsg(m)
	return err
}

func (w *wrapServerStream) RecvMsg(m interface{}) error {
	var err error
	startTime := time.Now()

	defer func() {
		latency := time.Now().Sub(startTime)
		code := codes.OK

		if err != nil && err != io.EOF {
			s, _ := status.FromError(err)
			code = s.Code()
		}

		common.DefaultRPCReceiveRequestMetric.With(prometheus.Labels{
			"sdk":                  common.RPCSDKGRPC,
			"request_protocol":     common.RPCProtocolGRPC,
			"request_target":       w.requestTarget,
			"request_path":         w.requestPath,
			"grpc_response_status": strconv.Itoa(int(code)),
			"response_code":        "0",
		}).Observe(latency.Seconds() * 1000)
	}()

	err = w.ServerStream.RecvMsg(m)
	return err
}

// ParseFullMethod returns a "/package.service/method" and "service" following the OpenTelemetry semantic
// conventions.
//
// Parsing is consistent with grpc-go implementation:
// https://github.com/grpc/grpc-go/blob/v1.57.0/internal/grpcutil/method.go#L26-L39
func ParseFullMethod(input string) (fullMethod, service string) {
	if !strings.HasPrefix(input, "/") {
		// Invalid format, does not follow `/package.service/method`.
		return input, common.RPCUnknownString
	}
	methodName := input[1:]

	pos := strings.LastIndex(methodName, "/")
	if pos < 0 {
		return input, common.RPCUnknownString
	}
	return input, methodName[:pos]
}
