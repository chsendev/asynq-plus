# 基于asynq的异步工作流封装
## 快速入门
1、安装
```go
go get -u github.com/chsendev/asynq-plus
```

2、编写工作流任务
```go
func SayHello(ctx context.Context, name string, price float64) (string, float64, error) {
	log.Println("start")
	time.Sleep(time.Second * 3)
	log.Println("sayHello: ", name)
	return "sayHello: " + name, price * 10, nil
}
```

3、注册服务端
```go
func TestServer(t *testing.T) {
	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: "127.0.0.1:6379"},
		asynq.Config{
			// Specify how many concurrent workers to use
			Concurrency: 10,
			// Optionally specify multiple queues with different priority.
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
			// See the godoc for other configuration options
		},
	)

	mux := NewServeFuture()
	mux.HandleFunc(SayHello)
	// ...register other handlers...
	if err := srv.Run(mux); err != nil {
		log.Fatalf("could not run server: %v", err)
	}
}
```

4、使用客户端执行任务，并获取结果
```go
func TestClient(t *testing.T) {
	client := NewClient(asynq.RedisClientOpt{Addr: "127.0.0.1:6379"})
	defer client.Close()

	var result string
	var f float64

	err := client.Enqueue(context.Background(), SayHello, "chsendev", 20).Get(&result, &f)
	fmt.Println(result, f, err)
}
```

更多高级功能请参考：https://github.com/hibiken/asynq/wiki