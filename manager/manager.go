package manager

import "fmt"

type Manager struct {
	Pending       queue.Queue
	TaskDb        map[string][]Task
	Event         map[string][]TaskEvent
	EventDb       map[String][]TaskEvent
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
