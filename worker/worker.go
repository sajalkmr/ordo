package worker

import (
	"github.com/google/uuid"
	"github.com/golang-collections/collections/queue"

	"GoOrchestrate/task"
)

type Worker struct {{
	Name string
	Queue queue.Queue
	Db map[uuid.UUID]task.Task
	TaskCount int
}
