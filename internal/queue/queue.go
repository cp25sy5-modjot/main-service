package queue

import "github.com/hibiken/asynq"

type Queue interface {
	Enqueue(task *asynq.Task, opts ...asynq.Option) (*asynq.TaskInfo, error)
}
