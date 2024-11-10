package asynqplus

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hibiken/asynq"
	"reflect"
	"time"
)

var defaultOpts = []asynq.Option{
	asynq.Retention(time.Hour * 72),
	asynq.MaxRetry(2),
	asynq.Timeout(time.Hour),
}

type Client struct {
	*asynq.Client
	inspector *asynq.Inspector
	opts      []asynq.Option
}

func NewClient(conn asynq.RedisConnOpt, opts ...asynq.Option) *Client {
	optMap := make(map[asynq.Option]struct{})
	for _, opt := range defaultOpts {
		optMap[opt] = struct{}{}
	}
	for _, opt := range opts {
		optMap[opt] = struct{}{}
	}
	opts = []asynq.Option{}
	for k := range optMap {
		opts = append(opts, k)
	}
	return &Client{Client: asynq.NewClient(conn), inspector: asynq.NewInspector(conn), opts: opts}
}

func (c *Client) Enqueue(ctx context.Context, fun any, paramsOrOpts ...any) Future {
	if reflect.TypeOf(fun).Kind() != reflect.Func {
		return newFutureErr(fmt.Errorf("the second parameter requires a transfer function"))
	}
	typeName := GetFunctionName(fun)
	paramMap := make(map[string]json.RawMessage)
	optMap := make(map[asynq.Option]struct{})
	for i := range c.opts {
		optMap[c.opts[i]] = struct{}{}
	}
	for i, v := range paramsOrOpts {
		if opt, ok := v.(asynq.Option); ok {
			optMap[opt] = struct{}{}
		} else {
			paramMap[paramName(i)], _ = json.Marshal(v)
		}
	}

	payload := Payload{
		Params: paramMap,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return newFutureErr(err)
	}
	var opts []asynq.Option
	for k := range optMap {
		opts = append(opts, k)
	}
	taskInfo, err := c.EnqueueContext(ctx, asynq.NewTask(typeName, jsonData), opts...)
	if err != nil {
		if errors.Is(err, asynq.ErrTaskIDConflict) || errors.Is(err, asynq.ErrDuplicateTask) {
			return newFutureFetcher(ctx, c.inspector, typeName, opts)
		}
		return newFutureErr(err)
	}
	return newFuturePoller(ctx, taskInfo, c.inspector)
}
