package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
	st "github.com/sithuaung/task-queue-redis/structs"
)

var rdb *redis.Client

func init() {
	redisAddr := os.Getenv("REDIS_ADDR")
	rdb = redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})
}

func main() {
	ctx := context.Background()
	for {
		taskJSON, err := rdb.RPop(ctx, "task_queue").Result()
		if err == redis.Nil {
			log.Println("No tasks in queue. Retrying...")
			time.Sleep(2 * time.Second)
			continue
		} else if err != nil {
			log.Fatalf("Failed to fetch task: %v", err)
		}

		var task st.Task
		if err := json.Unmarshal([]byte(taskJSON), &task); err != nil {
			log.Printf("Failed to parse task: %v", err)
			continue
		}

		log.Printf("Processing task: %+v", task)

		processTask(ctx, task.ID)

		// Simulate task processing
		log.Printf("Task %s completed successfully!", task.Type)
	}
}

func processTask(ctx context.Context, taskID string) error {
	taskKey := fmt.Sprintf("task:%s", taskID)

	// Fetch task details
	taskData, err := rdb.Get(ctx, taskKey).Result()
	if err != nil {
		return fmt.Errorf("error fetching task: %w", err)
	}

	var task st.Task
	if err := json.Unmarshal([]byte(taskData), &task); err != nil {
		return fmt.Errorf("error unmarshalling task: %w", err)
	}

	// Update status
	task.Status = "in-progress"
	taskData, err = marshalAndSet(ctx, rdb, taskKey, task)
	if err != nil {
		return fmt.Errorf("error updating task status: %w", err)
	}

	// Process task
	time.Sleep(time.Duration(rand.Intn(5)+1) * time.Second)

	// Update completion
	task.Status = "completed"
	task.Result = "Task processed successfully!"
	_, err = marshalAndSet(ctx, rdb, taskKey, task)
	if err != nil {
		return fmt.Errorf("error saving task result: %w", err)
	}

	return nil
}

func marshalAndSet(
	ctx context.Context,
	rdb *redis.Client,
	key string,
	value interface{},
) (string, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return "", err
	}

	return rdb.Set(ctx, key, string(data), 0).Result()
}
