# Go Contrib

[![CI](https://github.com/quwan-sre/observability-go-contrib/actions/workflows/e2e-test.yml/badge.svg)](https://github.com/quwan-sre/observability-go-contrib/actions?query=branch%3Amaster)
[![codecov](https://codecov.io/gh/quwan-sre/observability-go-contrib/graph/badge.svg?token=SQMXEVX0R4)](https://codecov.io/gh/quwan-sre/observability-go-contrib)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/quwan-sre/observability-go-contrib)](https://pkg.go.dev/github.com/quwan-sre/observability-go-contrib)
[![Go Report Card](https://goreportcard.com/badge/github.com/quwan-sre/observability-go-contrib)](https://goreportcard.com/report/github.com/quwan-sre/observability-go-contrib)

为 Go Packages 提供可观测性的扩展插件。

## 目标与边界
针对特定的 Go Package 提供 `middleware`、`interceptor`、`plugin`、`hook`、`callback` 等形式的可观测性埋点支持，Go 应用在启动、初始化时加载注册本项目组件，使特定操作（如网络调用、中间件调用）上报可观测性数据。

扩展插件只实现于特定 Package 提供的扩展接口之上，并且统一定义数据模型（如指标名、指标 Labels），但不干预用户在服务内做自定义的埋点。

扩展插件使用以下 Package 进行埋点：
- metrics：`github.com/prometheus/client_golang`；
- traces：`github.com/open-telemetry/opentelemetry-go`；
- logs： #TODO

## 接入使用
请阅读具体扩展插件目录下的 `README.md` 了解如何使用，如 [`go-contrib/metrics/gin/README.md`](https://github.com/quwan-sre/observability-go-contrib/blob/master/metrics/gin/README.md) 提供了对 `github.com/gin-gonic/gin` HTTP 框架 metrics 插装的指引。

## Metrics 定义
### RPC Metrics
**接收请求指标**

指标名称：
- `apm_rpc_receive_request_duration_milliseconds_count`
- `apm_rpc_receive_request_duration_milliseconds_sum`
- `apm_rpc_receive_request_duration_milliseconds_bucket`
包含 Labels：
- `sdk`：插装 SDK，枚举值见 [quwan-sre/observability-go-contrib](https://github.com/quwan-sre/observability-go-contrib/blob/master/metrics/common/rpc_enum.go)
- `request_protocol`：请求协议，枚举值：`http` | `grpc`；
- `request_target`：请求目标，HTTP 请求为域名，gRPC 请求为服务名，示例值：`www.baidu.com` | `RouteGuide` | `unknown`；
- `request_path`：请求路径，示例值：`/api/v1/foo/bar` | `/routeguide.RouteGuide/RouteChat`；
- `grpc_response_status`：RPC 状态码，遵循 [GRPC Core: Status codes and their use in gRPC](https://grpc.github.io/grpc/core/md_doc_statuscodes.html)，通常 `0` 为 OK；
- `response_code`：HTTP 状态码，遵循 [HTTP/1.1: Status Code Definitions](https://www.w3.org/Protocols/rfc2616/rfc2616-sec10.html)，通常 `200` 为 OK

代码定义：
```go
	DefaultRPCReceiveRequestMetric = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:                            DefaultRPCReceiveRequestMetricName,
		Buckets:                         []float64{0.5, 1, 5, 10, 25, 50, 100, 250, 500, 1000, 2500, 5000, 10000, 30000, 60000, 300000, 600000, 1800000, 3600000},
		NativeHistogramBucketFactor:     1.4,
		NativeHistogramZeroThreshold:    0.5,
		NativeHistogramMinResetDuration: 5 * time.Minute,
		NativeHistogramMaxZeroThreshold: 50,
		NativeHistogramMaxBucketNumber:  20,
	}, []string{"sdk", "request_protocol", "request_target", "request_path", "grpc_response_status", "response_code"})
```

---
**发送请求指标**

指标名称：
- `apm_rpc_send_request_duration_milliseconds_count`
- `apm_rpc_send_request_duration_milliseconds_sum`
- `apm_rpc_send_request_duration_milliseconds_bucket`
  包含 Labels：
- `sdk`：插装 SDK，枚举值见 [quwan-sre/observability-go-contrib](https://github.com/quwan-sre/observability-go-contrib/blob/master/metrics/common/rpc_enum.go)；
- `request_protocol`：请求协议，枚举值：http | grpc；
- `request_target`：请求目标，HTTP 请求为域名，gRPC 请求为服务名，示例值：www.baidu.com | RouteGuide | unknown；
- `request_path`：请求路径，示例值：127.0.0.1:8080/health | /routeguide.RouteGuide/RouteChat；
- `grpc_response_status`：RPC 状态码，遵循 GRPC Core: Status codes and their use in gRPC，通常 0 为 OK；
- `response_code`：HTTP 状态码，遵循 HTTP/1.1: Status Code Definitions，通常 200 为 OK。

代码定义：
```go
	DefaultRPCSendRequestMetric = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:                            DefaultRPCSendRequestMetricName,
		Buckets:                         []float64{0.5, 1, 5, 10, 25, 50, 100, 250, 500, 1000, 2500, 5000, 10000, 30000, 60000, 300000, 600000, 1800000, 3600000},
		NativeHistogramBucketFactor:     1.4,
		NativeHistogramZeroThreshold:    0.5,
		NativeHistogramMinResetDuration: 5 * time.Minute,
		NativeHistogramMaxZeroThreshold: 50,
		NativeHistogramMaxBucketNumber:  20,
	}, []string{"sdk", "request_protocol", "request_target", "request_path", "grpc_response_status", "response_code"})
```

备注
1. 由于 etcd 同样使用 gRPC 协议通信，并且使用 gRPC Interceptor 埋点采集，因此指标中包含与 etcd 进行通信的数据；
2. 对于 gRPC Streaming 类型请求，同样记录在以上指标中，但是 gRPC Streaming 是双向通信（可单向接收、可单向发送、可同时收发），因此在请求耗时上会不同（测试时显著低）于 Unary gRPC 请求响应；
3. 对于发送请求指标中的 request_target， gRPC 请求填入的值为 gRPC 请求方法（格式为：`/package.service/method`）中的 service，由 proto 定义，而非 Kubernetes、Istio 中的 Service，它们可能并不相同；
4. 对于 `request_path` 和 `request_target`，若无法获取，则填入 `unknown`。
   
### Database Metrics
**发送请求指标**

指标名称：
- `apm_database_send_request_duration_milliseconds_count`
- `apm_database_send_request_duration_milliseconds_sum`
- `apm_database_send_request_duration_milliseconds_bucket`

包含 Labels：
- `sdk`：插装 SDK，枚举值见 [quwan-sre/observability-go-contrib](https://github.com/quwan-sre/observability-go-contrib/blob/master/metrics/common/database_enum.go)；
- `database_type`：数据库类型，枚举值：`mysql` | `mongodb`；
- `database_addr`：数据库地址，通过 SDK 元数据获取或用户填入，示例值：`127.0.0.1:3306`；
- `response_status`：请求响应状态：枚举值：成功 `0` | 错误 `1`；
- query_type：语句类型，随不同数据库类型而不同，示例值： `create` | `dropDatabase` | `ping`；

代码定义：
```go
	DefaultDatabaseSendRequestMetric = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:                            DefaultDatabaseSendRequestMetricName,
		Buckets:                         []float64{1, 2.5, 5, 10, 20, 50, 100, 500, 1000, 2500, 5000, 7500, 10000},
		NativeHistogramBucketFactor:     1.4,
		NativeHistogramZeroThreshold:    1,
		NativeHistogramMinResetDuration: 5 * time.Minute,
		NativeHistogramMaxZeroThreshold: 50,
		NativeHistogramMaxBucketNumber:  20,
	}, []string{"sdk", "database_type", "database_addr", "response_status", "query_type"})
```
### 缓存指标
**发送请求指标**

指标名称：
- `apm_cache_send_request_duration_milliseconds_count`
- `apm_cache_send_request_duration_milliseconds_sum`
- `apm_cache_send_request_duration_milliseconds_bucket`

包含 Labels：
- `sdk`：插装 SDK，枚举值见 [quwan-sre/observability-go-contrib](https://github.com/quwan-sre/observability-go-contrib/blob/master/metrics/common/cache_enum.go)；
- `cache_type`：数据库类型，枚举值：`redis` | `memcached`；
- `cache_addr`：数据库地址，通过 SDK 元数据获取或用户填入，示例值：`127.0.0.1:6379`；
- `response_status`：请求响应状态：枚举值：成功 `0` | 错误 `1`；
- `command`：语句类型，随不同缓存类型而不同，示例值： `get` | `pipeline`；

代码定义：
```go
	DefaultCacheRequestMetric = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:                            DefaultCacheRequestMetricName,
		Buckets:                         []float64{0.25, 0.5, 1, 2, 5, 10, 25, 50, 100, 250, 500, 1000, 3000, 5000},
		NativeHistogramBucketFactor:     1.4,
		NativeHistogramZeroThreshold:    0.1,
		NativeHistogramMinResetDuration: 5 * time.Minute,
		NativeHistogramMaxZeroThreshold: 10,
		NativeHistogramMaxBucketNumber:  20,
	}, []string{"sdk", "cache_type", "cache_addr", "command", "response_status"})
```

### 消息队列指标
**发送消息指标**：
- `apm_mq_send_msg_count`

包含 Labels：
- `sdk`：插装 SDK，枚举值见 [quwan-sre/observability-go-contrib](https://github.com/quwan-sre/observability-go-contrib/blob/master/metrics/common/mq_enum.go)；
- `mq_type`：MQ 类型，枚举值：`kafka`；
- `mq_addr`：MQ Broker 地址，通过 SDK 元数据获取或用户填入，示例值：`127.0.0.1:9092`；
- `mq_topic`：MQ Topic，示例值：`my_kafka_topic`；
- `mq_partition`：MQ 分区，若发送时未知则为 `unknown`；

代码定义：
```go
	DefaultMQSendMsgMetric = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: DefaultMQSendMsgMetricName,
	}, []string{"sdk", "mq_type", "mq_addr", "mq_topic", "mq_partition"})
```

**消费消息指标**：
- `apm_mq_receive_msg_count`

包含 Labels：
- `sdk`：插装 SDK，枚举值见 [quwan-sre/observability-go-contrib](https://github.com/quwan-sre/observability-go-contrib/blob/master/metrics/common/mq_enum.go)；
- `mq_type`：MQ 类型，枚举值：`kafka`；
- `mq_addr`：MQ Broker 地址，通过 SDK 元数据获取或用户填入，示例值：`127.0.0.1:9092`；
- `mq_topic`：MQ Topic，示例值：`my_kafka_topic`；
- `mq_partition`：MQ 分区，示例值：`0` | `1` | `2`；

代码定义：
```go
	DefaultMQReceiveMsgMetric = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: DefaultMQReceiveMsgMetricName,
	}, []string{"sdk", "mq_type", "mq_addr", "mq_topic", "mq_partition"})
```

## Go Packages 支持列表

| 分类    | Package                      | 数据类型  & 稳定性状态                                                                                           |
|---------|------------------------------|---------------------------------------------------------------------------------------------------------|
| HTTP    | github.com/gin-gonic/gin     | alpha:  [metrics](https://github.com/quwan-sre/observability-go-contrib/tree/master/metrics/gin)        |
| HTTP    | github.com/go-resty/resty/v2 | alpha:  [metrics](https://github.com/quwan-sre/observability-go-contrib/tree/master/metrics/resty)        |
| RPC     | google.golang.org/grpc       | alpha:  [metrics](https://github.com/quwan-sre/observability-go-contrib/tree/master/metrics/grpc)         |
| Redis   | github.com/go-redis/redis/v9 | alpha:  [metrics](https://github.com/quwan-sre/observability-go-contrib/tree/master/metrics/go-redis)     |
| Kafka   | github.com/Shopify/sarama    | alpha:  [metrics](https://github.com/quwan-sre/observability-go-contrib/tree/master/metrics/sarama)       |
| MongoDB | go.mongodb.org/mongo-driver  | alpha:  [metrics](https://github.com/quwan-sre/observability-go-contrib/tree/master/metrics/mongo-driver) |
| etcd    | go.etcd.io/etcd/client/v3    | alpha:  [metrics](https://github.com/quwan-sre/observability-go-contrib/tree/master/metrics/etcd)         |
