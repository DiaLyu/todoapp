package domain

import "time"

type Statistics struct {
	TaskCreated               int
	TaskCompleted             int
	TaskCompletedRate         *float64
	TaskAverageCompletionTime *time.Duration
}

func NewStatistics(
	taskCreated int,
	taskCompleted int,
	taskCompletedRate *float64,
	taskAverageCompletionTime *time.Duration,
) Statistics {
	return Statistics{
		TaskCreated:               taskCreated,
		TaskCompleted:             taskCompleted,
		TaskCompletedRate:         taskCompletedRate,
		TaskAverageCompletionTime: taskAverageCompletionTime,
	}
}
