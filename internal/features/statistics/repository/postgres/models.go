package statistics_postgres_repository

import (
	"time"

	"github.com/DiaLyu/todoapp/internal/core/domain"
)

type TaskModel struct {
	ID           int
	Version      int
	Title        string
	Description  *string
	Completed    bool
	CreateAt     time.Time
	CompletedAt  *time.Time
	AuthorUserID int
}

func taskDomainFromModel(taskModel TaskModel) domain.Task {
	return domain.NewTask(
		taskModel.ID,
		taskModel.Version,
		taskModel.Title,
		taskModel.Description,
		taskModel.Completed,
		taskModel.CreateAt,
		taskModel.CompletedAt,
		taskModel.AuthorUserID,
	)
}

func taskDomainsFromModels(taskModels []TaskModel) []domain.Task {
	domains := make([]domain.Task, len(taskModels))

	for i, model := range taskModels {
		domains[i] = taskDomainFromModel(model)
	}

	return domains
}
