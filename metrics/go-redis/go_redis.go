package go_redis

import (
	"context"
	"sort"
	"strings"
	"time"

	rdsV9 "github.com/go-redis/redis/v9"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/quwan-sre/observability-go-contrib/metrics/common"
)

type metricsHook struct {
	labels prometheus.Labels
}

type metricsCtxType struct{}

var (
	metricsCtxKey = metricsCtxType{}
)

func NewRedisHook(dst []string) rdsV9.Hook {
	sort.Slice(dst, func(i, j int) bool {
		if dst[i] < dst[j] {
			return true
		}
		return false
	})

	host := "unknown"
	if len(dst) >= 1 {
		host = strings.Join(dst, ";")
	}

	return metricsHook{
		labels: prometheus.Labels{
			"cache_addr": host,
		},
	}
}

func (h metricsHook) BeforeProcess(ctx context.Context, cmd rdsV9.Cmder) (context.Context, error) {
	ctx = context.WithValue(ctx, metricsCtxKey, time.Now())
	return ctx, nil
}

func (h metricsHook) AfterProcess(ctx context.Context, cmd rdsV9.Cmder) error {
	val := ctx.Value(metricsCtxKey)
	if val == nil {
		return nil
	}

	startTime, ok := val.(time.Time)
	if !ok {
		return nil
	}

	latency := time.Now().Sub(startTime)
	responseStatus := common.CacheResponseStatusSuccess
	if cmd.Err() != nil && cmd.Err() != rdsV9.Nil {
		responseStatus = common.CacheResponseStatusError
	}

	labels := prometheus.Labels{
		"sdk":             common.CacheSDKGoRedis,
		"cache_type":      common.CacheTypeRedis,
		"response_status": responseStatus,
		"command":         cmd.Name(),
	}

	// Added cache_host
	for k, v := range h.labels {
		labels[k] = v
	}

	common.DefaultCacheRequestMetric.With(labels).Observe(latency.Seconds())
	return nil
}

func (h metricsHook) BeforeProcessPipeline(ctx context.Context, cmds []rdsV9.Cmder) (context.Context, error) {
	ctx = context.WithValue(ctx, metricsCtxKey, time.Now())
	return ctx, nil
}

func (h metricsHook) AfterProcessPipeline(ctx context.Context, cmds []rdsV9.Cmder) error {
	val := ctx.Value(metricsCtxKey)
	if val == nil {
		return nil
	}

	startTime, ok := val.(time.Time)
	if !ok {
		return nil
	}

	latency := time.Now().Sub(startTime)
	responseStatus := common.CacheResponseStatusSuccess
	for i := range cmds {
		if cmds[i].Err() != nil && cmds[i].Err() != rdsV9.Nil {
			responseStatus = common.CacheResponseStatusError
			break
		}
	}

	labels := prometheus.Labels{
		"sdk":             common.CacheSDKGoRedis,
		"cache_type":      common.CacheTypeRedis,
		"response_status": responseStatus,
		"command":         common.RedisCmdPipeline,
	}

	// Added cache_host
	for k, v := range h.labels {
		labels[k] = v
	}

	common.DefaultCacheRequestMetric.With(labels).Observe(latency.Seconds())
	return nil
}
