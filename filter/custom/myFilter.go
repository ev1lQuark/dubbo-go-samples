package custom

import (
	"context"
	"dubbo.apache.org/dubbo-go/v3/common/extension"
	"dubbo.apache.org/dubbo-go/v3/filter"
	"dubbo.apache.org/dubbo-go/v3/protocol"
	"github.com/dubbogo/gost/log/logger"
	"time"
)

func init() {
	extension.SetFilter("myFilter", NewMyFilter)
}

func NewMyFilter() filter.Filter {
	return &MyFilter{}
}

type MyFilter struct {
}

func (f *MyFilter) Invoke(ctx context.Context, invoker protocol.Invoker, invocation protocol.Invocation) protocol.Result {
	begin := time.Now().UnixMilli()
	invocation.SetAttribute("rpc_begin", begin)
	return invoker.Invoke(ctx, invocation)
}
func (f *MyFilter) OnResponse(ctx context.Context, result protocol.Result, invoker protocol.Invoker, protocol protocol.Invocation) protocol.Result {
	rpcBegin, _ := protocol.GetAttribute("rpc_begin")
	var side = "consumer"
	if invoker != nil && invoker.GetURL() != nil {
		side = invoker.GetURL().GetParam("side", "consumer")
	}
	logger.Infof("gaia-resource-srv|rpcAccess|%v|%v|%v|%v|%v", side, protocol.MethodName(), time.Now().UnixMilli()-rpcBegin.(int64), protocol.Arguments(), result.Error())
	return result
}
