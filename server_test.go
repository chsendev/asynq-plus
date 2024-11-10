package asynqplus

import (
	"github.com/hibiken/asynq"
	"log"
	"testing"
)

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
	mux.HandleFunc(SayHello2)
	// ...register other handlers...
	if err := srv.Run(mux); err != nil {
		log.Fatalf("could not run server: %v", err)
	}
}
