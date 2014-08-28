package task_bbs

import (
	"github.com/cloudfoundry-incubator/runtime-schema/bbs/shared"
	"github.com/cloudfoundry-incubator/runtime-schema/models"
	"github.com/cloudfoundry/storeadapter"
	"github.com/pivotal-golang/lager"
)

func (bbs *TaskBBS) GetAllTasks() ([]models.Task, error) {
	node, err := bbs.store.ListRecursively(shared.TaskSchemaRoot)
	if err == storeadapter.ErrorKeyNotFound {
		return []models.Task{}, nil
	}

	if err != nil {
		return []models.Task{}, err
	}

	tasks := []models.Task{}
	for _, node := range node.ChildNodes {
		task, err := models.NewTaskFromJSON(node.Value)
		if err != nil {
			bbs.logger.Error("failed-to-unmarshal-task", err, lager.Data{
				"key":   node.Key,
				"value": node.Value,
			})
		} else {
			tasks = append(tasks, task)
		}
	}

	return tasks, nil
}

func (bbs *TaskBBS) GetAllPendingTasks() ([]models.Task, error) {
	all, err := bbs.GetAllTasks()
	return filterTasks(all, models.TaskStatePending), err
}

func (bbs *TaskBBS) GetAllClaimedTasks() ([]models.Task, error) {
	all, err := bbs.GetAllTasks()
	return filterTasks(all, models.TaskStateClaimed), err
}

func (bbs *TaskBBS) GetAllRunningTasks() ([]models.Task, error) {
	all, err := bbs.GetAllTasks()
	return filterTasks(all, models.TaskStateRunning), err
}

func (bbs *TaskBBS) GetAllCompletedTasks() ([]models.Task, error) {
	all, err := bbs.GetAllTasks()
	return filterTasks(all, models.TaskStateCompleted), err
}

func (bbs *TaskBBS) GetAllResolvingTasks() ([]models.Task, error) {
	all, err := bbs.GetAllTasks()
	return filterTasks(all, models.TaskStateResolving), err
}

func filterTasks(tasks []models.Task, state models.TaskState) []models.Task {
	result := make([]models.Task, 0)
	for _, model := range tasks {
		if model.State == state {
			result = append(result, model)
		}
	}
	return result
}
