package asynqplus

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	"log"
	"reflect"
	"testing"
	"time"
)

func TestMarshalCtx(t *testing.T) {
	ctx := context.Background()
	c, e := json.Marshal(ctx)
	fmt.Println(c, e)
}

func Get(result ...any) {
	for i, a := range result {
		fmt.Println(i, a)
	}
	Get2(result)
}

func Get2(result ...any) {
	for i, a := range result {
		fmt.Println(i, a)
	}
}

func TestGet(t *testing.T) {
	Get(1, "A")
}

func TestFutureBase(t *testing.T) {
	fu := futureBase{
		ctx:       nil,
		inspector: nil,
	}
	var r string
	var f float64
	fmt.Println(fu.unmarshal(&asynq.TaskInfo{Result: []byte(`{"result":{"result0":"sayHello: ","result1":3.14,"result2":null}}`)}, []any{&r, &f}))
	fmt.Println(r, f)
}

func TestUnmarshal(t *testing.T) {
	s := `{"result":{"result0":"sayHello: ","result1":3.14,"result2":null}}`
	var res Result
	err := json.Unmarshal([]byte(s), &res)
	if err != nil {
		panic(err)
	}
	var r string
	var f float64
	result := []any{
		&r, &f,
	}

	for i, v := range result {
		err = json.Unmarshal(res.Result[resultName(i)], v)
		if err != nil {
			panic(err)
		}
	}
	fmt.Println(r, f)
}

func Say() *error {
	//return reflect.TypeOf((*error)(nil)).Elem()
	var err *error
	return err
}

func TestErr(t *testing.T) {
	s := reflect.ValueOf(Say)
	r := s.Call(nil)
	fmt.Println(r)
	for _, v := range r {
		if err, ok := v.Interface().(*error); ok && err != nil {
			fmt.Println(err)
		}
	}
	//var err2 error
	//
	//fmt.Println(reflect.ValueOf(err2).Interface().(error))
}

func TestClient(t *testing.T) {
	client := NewClient(asynq.RedisClientOpt{Addr: "127.0.0.1:6379"})
	defer client.Close()

	var result string
	var f float64

	err := client.Enqueue(context.Background(), SayHello, "chsendev", 20).Get(&result, &f)
	fmt.Println(result, f, err)
}

func SayHello(ctx context.Context, name string, price float64) (string, float64, error) {
	log.Println("start")
	time.Sleep(time.Second * 3)
	log.Println("sayHello: ", name)
	return "sayHello: " + name, price * 10, nil
}

func TestClient2(t *testing.T) {
	client := NewClient(asynq.RedisClientOpt{Addr: "127.0.0.1:6379"})
	defer client.Close()

	input := &SayHelloInput{
		Name:  "chsendev",
		Price: 30,
	}

	var ouput SayHelloOutput

	err := client.Enqueue(context.Background(), SayHello2, input).Get(&ouput)
	fmt.Println(ouput, err)
}

type SayHelloInput struct {
	Name  string
	Price float64
}

type SayHelloOutput struct {
	Name  string
	Price float64
}

func SayHello2(ctx context.Context, input *SayHelloInput) (*SayHelloOutput, error) {
	log.Println("start")
	time.Sleep(time.Second * 3)
	log.Println("input: ", input)
	output := &SayHelloOutput{
		Name:  "sayHello: " + input.Name,
		Price: input.Price * 100,
	}
	return output, nil
}
