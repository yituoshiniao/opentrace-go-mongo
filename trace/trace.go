package trace

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	traceLog "github.com/opentracing/opentracing-go/log"
	"go.mongodb.org/mongo-driver/event"
	"sync"
)

type Trace struct {
	spans         sync.Map
	ctxs          sync.Map
	IsAccessTrace bool
}

func NewTrace(IsAccessTrace bool) *Trace {
	return &Trace{}
}

const (
	Prefix        = "mongodb."
	ComponentName = "golang-mongo"
	DbType        = "mongo-"
	ErrorMsg      = "errorMsg"
	DbExecInfo    = "mongoDbExecInfo"
)

var (
	//执行命令黑名单
	BlackCommandArr = []string{
		"ping",
	}
	//执行命令白名单
	WhiteCommandArr = []string{
		"find",
	}
)

//是否执行记录
func isExec(command string) bool {
	for _, v := range BlackCommandArr {
		if command == v {
			return true
		}
	}
	return false
}

func (t *Trace) HandleStartedEvent(ctx context.Context, evt *event.CommandStartedEvent) {
	//if !t.IsAccessTrace {
	//	return
	//}

	if isExec(evt.CommandName) {
		return
	}

	if evt == nil {
		return
	}


	//上报trace
	span, tmpCtx := opentracing.StartSpanFromContext(ctx, Prefix+evt.CommandName)
	ext.DBType.Set(span, DbType)
	ext.DBInstance.Set(span, evt.DatabaseName)
	span.SetTag("db.host", evt.ConnectionID)
	span.SetTag("CommandName", evt.CommandName)
	ext.SpanKind.Set(span, ext.SpanKindRPCClientEnum)
	ext.Component.Set(span, ComponentName)

	//span.LogKV("db.statement", evt.Command.String())
	span.LogFields(traceLog.String("DB Exec", evt.Command.String()))

	t.spans.Store(evt.RequestID, span) //上下文传递
	t.ctxs.Store(evt.RequestID, tmpCtx)


	//记录日志 带trace链路追踪的
	//log.L(tmpCtx).Info(
	//	DbType+"HandleStartedEvent开始事件",
	//	zap.String("DatabaseName", evt.DatabaseName),
	//)
}

func (t *Trace) HandleSucceededEvent(ctx context.Context, evt *event.CommandSucceededEvent) {
	//if !t.IsAccessTrace {
	//	return
	//}

	if isExec(evt.CommandName) {
		return
	}

	if evt == nil {
		return
	}

	if rawSpan, ok := t.spans.Load(evt.RequestID); ok {
		//获取正确的ctx上下文
		var isCtx bool = false
		if rawCtx, ok := t.ctxs.Load(evt.RequestID); ok {
			if ctxTmp, ok := rawCtx.(context.Context); ok {
				ctx = ctxTmp
				isCtx = true
			}
		}


		defer t.spans.Delete(evt.RequestID)
		if isCtx {
			defer t.ctxs.Delete(evt.RequestID)
		}
		if span, ok := rawSpan.(opentracing.Span); ok {
			defer span.Finish()
			//span.SetTag(Prefix+"reply", evt.Reply.String())
			span.SetTag(Prefix+"timeNs", evt.DurationNanos)
			span.SetTag(Prefix+"timeMs", evt.DurationNanos/(1000*1000))
		}
	}
}

func (t *Trace) HandleFailedEvent(ctx context.Context, evt *event.CommandFailedEvent) {
	//if !t.IsAccessTrace {
	//	return
	//}
	if isExec(evt.CommandName) {
		return
	}
	if evt == nil {
		return
	}
	if rawSpan, ok := t.spans.Load(evt.RequestID); ok {
		//获取正确的ctx上下文
		var isCtx bool = false
		if rawCtx, ok := t.ctxs.Load(evt.RequestID); ok {
			if ctxTmp, ok := rawCtx.(context.Context); ok {
				ctx = ctxTmp
				isCtx = true
			}
		}

		defer t.spans.Delete(evt.RequestID)
		if isCtx {
			defer t.ctxs.Delete(evt.RequestID)
		}
		if span, ok := rawSpan.(opentracing.Span); ok {
			defer span.Finish()
			ext.Error.Set(span, true)
			span.SetTag("errorMsg", evt.Failure)
		}
	}
}


