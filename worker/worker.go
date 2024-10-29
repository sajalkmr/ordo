package worker

import (
	"fmt"

	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"

	"github.com/sajalkmr/ordo/task"
)

type Worker struct {
	Name      string
	Queue     queue.Queue
	Db        map[uuid.UUID]task.Task
	TaskCount int
}

func (w *Worker) CollectStats() {
	fmt.Println("Collect stats")
}

func (w *Worker) RunTask() {
	fmt.Println("Runn task")
}

func (w *Worker) StartTask() {
	fmt.Println("Start task")
}

func (w *Worker) StopTask() {
	fmt.Println("Stop task")
}
