package storage

import "main.go/pkg/storage/postgres"

type Interface interface {
	Tasks(int, int) ([]postgres.Task, error)
	NewTask(postgres.Task) (int, error)
	UpdateTask(t postgres.Task) error
	DeleteTask(taskID int) error
}
