package asynqplus

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	"log"
	"reflect"
)

var (
	ctxType = reflect.TypeOf((*context.Context)(nil)).Elem()
)

type ServeFuture struct {
	*asynq.ServeMux
}

func NewServeFuture() *ServeFuture {
	return &ServeFuture{asynq.NewServeMux()}
}

func (mux *ServeFuture) HandleFunc(handler any) {
	if handler == nil {
		panic("asynq: nil handler")
	}
	mux.Handle(getFunctionName(handler), asynq.HandlerFunc(func(ctx context.Context, task *asynq.Task) error {
		var p Payload
		err := json.Unmarshal(task.Payload(), &p)
		if err != nil {
			return fmt.Errorf("unmarshal err: " + err.Error())
		}
		funcValue := reflect.ValueOf(handler)
		funcType := funcValue.Type()
		var params []reflect.Value
		paramIndex := 0
		for i := 0; i < funcType.NumIn(); i++ {
			paramType := funcType.In(i)
			if paramType.Implements(ctxType) {
				params = append(params, reflect.ValueOf(ctx))
			} else {
				paramValue := reflect.New(paramType)
				if val, ok := p.Params[paramName(paramIndex)]; ok {
					err = json.Unmarshal(val, paramValue.Interface())
					if err != nil {
						log.Printf("Error unmarshaling for param %d: %v\n", i, err)
					}
				}
				params = append(params, paramValue.Elem())
				paramIndex++
			}
		}
		callResult := funcValue.Call(params)
		resultJson := make(map[string]json.RawMessage)
		var result Result
		for i, v := range callResult {
			if err, ok := v.Interface().(error); ok && err != nil {
				log.Println(err)
				return err
			} else {
				bys, err := json.Marshal(v.Interface())
				if err != nil {
					log.Printf("Error marshaling for result %d: %v\n", i, err)
				}
				resultJson[resultName(i)] = bys
			}
		}
		result.Result = resultJson
		bys, err := json.Marshal(result)
		if err != nil {
			log.Printf("Error marshaling for result %v\n", err)
			return err
		}
		_, err = task.ResultWriter().Write(bys)
		if err != nil {
			log.Printf("Result write error %v\n", err)
			return err
		}
		return nil
	}))
}
