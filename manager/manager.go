package manager

import (
	"fmt"

	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
	"github.com/sajalkmr/ordo/task"
)

type Manager struct {
	Pending       queue.Queue
	TaskDb        map[string][]task.Task
	Event         map[string][]task.TaskEvent
	EventDb       map[string][]task.TaskEvent // Fixed String -> string
	Workers       []string
	WorkerTaskMap map[string][]uuid.UUID
	TaskWorkerMap map[uuid.UUID]string
}

func (m *Manager) SelectWorker() {
	fmt.Println("Select worker")
}

func (m *Manager) UpdateTasks() {
	fmt.Println("Assign task")
}

func (m *Manager) SendWork() {
	fmt.Println("Send wrork to worker")
}
