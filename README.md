# Go Contrib

![ci](https://github.com/quwan-sre/observability-go-contrib/actions/workflows/e2e-test.yml/badge.svg) [![codecov](https://codecov.io/gh/quwan-sre/observability-go-contrib/graph/badge.svg?token=SQMXEVX0R4)](https://codecov.io/gh/quwan-sre/observability-go-contrib)

为 Go Packages 提供可观测性的扩展插件。

## 目标与边界
针对特定的 Go Package 提供 `middleware`、`interceptor`、`plugin`、`hook`、`callback` 等形式的可观测性埋点支持，Go 应用在启动、初始化时加载注册本项目组件，使特定操作（如网络调用、中间件调用）上报可观测性数据。

扩展插件只实现于特定 Package 提供的扩展接口之上，并且统一定义数据模型（如指标名、指标 Labels），但不干预用户在服务内做自定义的埋点。

扩展插件使用以下 Package 进行埋点：
- metrics：`github.com/prometheus/client_golang`；
- traces：`github.com/open-telemetry/opentelemetry-go`；
- logs： #TODO

## 接入使用
请阅读具体扩展插件目录下的 `README.md` 了解如何使用，如 `go-contrib/metrics/gin/README.md` 提供了对 `github.com/gin-gonic/gin` HTTP 框架 metrics 插装的指引。

### Metrics 定义
#### RPC Metrics
目前提供 RPC 调用的通用 Metrics：
- `apm_rpc_receive_request_duration_seconds`；
- `apm_rpc_send_request_duration_seconds`。

他们的定义及包含 Labels 可以查看 [rpc_metrics.go](./metrics/common/rpc_metrics.go)

## Go Packages 支持列表

| 分类    | Package | 数据类型 & 稳定性状态                                                                                  | 
|-------|---------|-----------------------------------------------------------------------------------------------|
| HTTP  | `github.com/gin-gonic/gin` | dev: [metrics](https://gitlab.ttyuyin.com/observability/go-contrib/-/blob/master/metrics/gin) |
| Kafka | `github.com/Shopify/sarama` | todo: metrics                                                                                 |
| MongoDB | `go.mongodb.org/mongo-driver` | todo: metrics                                                                                 |
| Redis | `github.com/go-redis/redis` | todo: metrics                                                                                 |
