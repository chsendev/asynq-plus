package asynqplus

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hibiken/asynq"
	"time"
)

type Future interface {
	Get(result ...any) error
}

type futureErr struct {
	err error
}

func newFutureErr(err error) Future {
	return &futureErr{err: err}
}

func (e *futureErr) Get(result ...any) error {
	return e.err
}

type futureBase struct {
	ctx       context.Context
	inspector *asynq.Inspector
}

func (f *futureBase) unmarshal(info *asynq.TaskInfo, result []any) error {
	if result == nil {
		return nil
	}
	var res Result
	err := json.Unmarshal(info.Result, &res)
	if err != nil {
		return err
	}
	for i, v := range result {
		err = json.Unmarshal(res.Result[resultName(i)], v)
		if err != nil {
			return err
		}
	}
	return err
}

type futurePoller struct {
	futureBase
	taskInfo *asynq.TaskInfo
}

func newFuturePoller(ctx context.Context, taskInfo *asynq.TaskInfo, inspector *asynq.Inspector) Future {
	return &futurePoller{futureBase: futureBase{ctx: ctx, inspector: inspector}, taskInfo: taskInfo}
}

func (f *futurePoller) Get(result ...any) error {
	for {
		select {
		case <-f.ctx.Done():
			return f.ctx.Err()
		case <-time.After(time.Second):
			info, err := f.inspector.GetTaskInfo(f.taskInfo.Queue, f.taskInfo.ID)
			if err != nil {
				return fmt.Errorf("getTaskInfo err: %w", err)
			}
			if info.State == asynq.TaskStateCompleted {
				return f.unmarshal(info, result)
			} else if info.State == asynq.TaskStateArchived {
				return errors.New(info.LastErr)
			}
		}
	}
}

type futureFetcher struct {
	futureBase
	typeName string
	taskId   string
}

func newFutureFetcher(ctx context.Context, inspector *asynq.Inspector, typeName string, opts []asynq.Option) Future {
	var taskId string
	for _, opt := range opts {
		if opt.Type() == asynq.TaskIDOpt {
			if id, ok := opt.Value().(string); ok {
				taskId = id
			}
		}
	}

	return &futureFetcher{futureBase: futureBase{ctx: ctx, inspector: inspector}, typeName: typeName, taskId: taskId}
}

func (f *futureFetcher) Get(result ...any) error {
	queues, err := f.inspector.Queues()
	if err != nil {
		return err
	}
	for _, queue := range queues {
		info, err := f.inspector.GetTaskInfo(queue, f.taskId)
		if err != nil {
			return err
		}
		if info.Type == f.typeName {
			return f.unmarshal(info, result)
		}
	}
	return errors.New("task not exists")
}
