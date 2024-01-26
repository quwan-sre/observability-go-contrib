## Resty

Resty 对外提供了 Middleware 机制，用于挂载 Middleware 扩展额外的请求逻辑。

```go
// OnBeforeRequest method appends a request middleware into the before request chain.
// The user defined middlewares get applied before the default Resty request middlewares.
// After all middlewares have been applied, the request is sent from Resty to the host server.
//
//	client.OnBeforeRequest(func(c *resty.Client, r *resty.Request) error {
//			// Now you have access to Client and Request instance
//			// manipulate it as per your need
//
//			return nil 	// if its success otherwise return error
//		})
func (c *Client) OnBeforeRequest(m RequestMiddleware) *Client {
	c.udBeforeRequestLock.Lock()
	defer c.udBeforeRequestLock.Unlock()
	
	c.udBeforeRequest = append(c.udBeforeRequest, m)
	
	return c
}

// OnAfterResponse method appends response middleware into the after response chain.
// Once we receive response from host server, default Resty response middleware
// gets applied and then user assigned response middlewares applied.
//
//	client.OnAfterResponse(func(c *resty.Client, r *resty.Response) error {
//			// Now you have access to Client and Response instance
//			// manipulate it as per your need
//
//			return nil 	// if its success otherwise return error
//		})
func (c *Client) OnAfterResponse(m ResponseMiddleware) *Client {
	c.afterResponseLock.Lock()
	defer c.afterResponseLock.Unlock()
	
	c.afterResponse = append(c.afterResponse, m)
	
	return c
}

```

在 Resty 框架中可以通过以下代码使用 go-contrib 埋点：
```go
package main

import (
	"github.com/go-resty/resty/v2"

	metrics "github.com/quwan-sre/observability-go-contrib/metrics/resty"
)

func main() {
	...
	restyClient := resty.New()
	
	// Interceptor 作为 grpc.DialOption 使用
	metrics.NewMetricsMiddleware(restyClient)
	...
}

//func NewMetricsMiddleware(c *resty.Client) {
//	if c == nil {
//		return
//	}
//
//	// 内部使用 Resty 提供的 OnBefore 和 OnAfter 方法增加埋点
//	c.OnBeforeRequest(NewBeforeRequest())
//	c.OnAfterResponse(NewAfterResponse())
//
//	return
//}
```