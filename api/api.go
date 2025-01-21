package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	st "github.com/sithuaung/task-queue-redis/structs"
)

var rdb *redis.Client

func init() {
	rdb = redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
	})
}

// Add a new task to the queue
func createTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task st.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "Invalid task format", http.StatusBadRequest)
		return
	}

	task.ID = uuid.New().String()
	task.Status = "queued"

	// Store the task in Redis
	taskKey := fmt.Sprintf("task:%s", task.ID)
	taskData, _ := json.Marshal(task)
	if err := rdb.Set(r.Context(), taskKey, taskData, 0).Err(); err != nil {
		http.Error(w, "Failed to save task", http.StatusInternalServerError)
		return
	}

	// Add task ID to the queue
	if err := rdb.LPush(r.Context(), "task_queue", task.ID).Err(); err != nil {
		http.Error(w, "Failed to enqueue task", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}

// Get the status of a task
func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	taskID := r.URL.Query().Get("id")
	if taskID == "" {
		http.Error(w, "Task ID is required", http.StatusBadRequest)
		return
	}

	taskKey := fmt.Sprintf("task:%s", taskID)
	taskData, err := rdb.Get(r.Context(), taskKey).Result()
	if err == redis.Nil {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Failed to fetch task", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(taskData))
}

// HTTP server setup
func main() {
	http.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			createTaskHandler(w, r)
		} else if r.Method == http.MethodGet {
			getTaskHandler(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	addr := ":8080"
	log.Printf("API server running on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
