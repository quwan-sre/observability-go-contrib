package mongo_driver

import (
	"context"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"go.mongodb.org/mongo-driver/event"

	"github.com/quwan-sre/observability-go-contrib/metrics/common"
)

func NewCommandMonitor() *event.CommandMonitor {
	return &event.CommandMonitor{
		Started:   Started,
		Succeeded: Succeeded,
		Failed:    Failed,
	}
}

func Started(ctx context.Context, evt *event.CommandStartedEvent) {
	return
}

func Succeeded(ctx context.Context, evt *event.CommandSucceededEvent) {
	common.DefaultDatabaseSendRequestMetric.With(prometheus.Labels{
		"sdk":             common.DatabaseSDKMongoDriver,
		"database_type":   common.DatabaseTypeMongoDB,
		"database_addr":   parseConnectionID(evt.ConnectionID),
		"response_status": common.DatabaseResponseStatusSuccess,
		"query_type":      evt.CommandName,
	}).Observe(float64(evt.DurationNanos) / float64(time.Millisecond))
}

func Failed(ctx context.Context, evt *event.CommandFailedEvent) {
	common.DefaultDatabaseSendRequestMetric.With(prometheus.Labels{
		"sdk":             common.DatabaseSDKMongoDriver,
		"database_type":   common.DatabaseTypeMongoDB,
		"database_addr":   parseConnectionID(evt.ConnectionID),
		"response_status": common.DatabaseResponseStatusError,
		"query_type":      evt.CommandName,
	}).Observe(float64(evt.DurationNanos) / float64(time.Millisecond))
}

func parseConnectionID(connectionID string) string {
	potentialAddr := strings.Split(connectionID, "[")
	if len(potentialAddr) >= 1 {
		return potentialAddr[0]
	}
	return connectionID
}
