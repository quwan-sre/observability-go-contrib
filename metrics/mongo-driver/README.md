## MongoDB Go Driver

MongoDB Go Driver 提供了一组方法对命令执行情况进行监控。
```go
// CommandMonitor represents a monitor that is triggered for different events.
type CommandMonitor struct {
	Started		func(context.Context, *CommandStartedEvent)
	Succeeded	func(context.Context, *CommandSucceededEvent)
	Failed		func(context.Context, *CommandFailedEvent)
}
```

这组方法可以在 Client 初始化时通过 SetMonitor 方法附带：
```go
// SetMonitor specifies a CommandMonitor to receive command events. See the event.CommandMonitor documentation for more
// information about the structure of the monitor and events that can be received.
func (c *ClientOptions) SetMonitor(m *event.CommandMonitor) *ClientOptions {
	c.Monitor = m
	return c
}
```

因此，在 MongoDB Go Driver 中可以通过以下代码使用 go-contrib 埋点：
```go
package metrics

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	mongo_driver "github.com/quwan-sre/observability-go-contrib/metrics/mongo-driver"
)

var mongoClient *mongo.Client

func initMongoDBClient() {
	...
	// 在初始化后使用 SetMonitor 方法增加埋点
	mongoClient, err = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://test_user:password@localhost:27017").SetMonitor(mongo_driver.NewCommandMonitor()))
	if err != nil {
		panic(fmt.Sprintf("init mongo db client err: %v", err))
	}
	...
}
```