package repo

import (
	"errors"
	"time"

	"github.com/dogefuzz/dogefuzz/entities"
	"github.com/dogefuzz/dogefuzz/pkg/common"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TaskRepo interface {
	Get(tx *gorm.DB, id string) (*entities.Task, error)
	Create(tx *gorm.DB, task *entities.Task) error
	Update(tx *gorm.DB, task *entities.Task) error
	FindNotFinishedTasksThatDontHaveIncompletedTransactions(tx *gorm.DB) ([]entities.Task, error)
	FindNotFinishedAndExpired(tx *gorm.DB) ([]entities.Task, error)
}

type taskRepo struct {
}

func NewTaskRepo(e Env) *taskRepo {
	return &taskRepo{}
}

func (r *taskRepo) Get(tx *gorm.DB, id string) (*entities.Task, error) {
	var task entities.Task
	if err := tx.First(&task, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotExists
		}
		return nil, err
	}
	return &task, nil
}

func (r *taskRepo) Create(tx *gorm.DB, task *entities.Task) error {
	task.Id = uuid.NewString()
	return tx.Create(task).Error
}

func (r *taskRepo) Update(tx *gorm.DB, updatedTask *entities.Task) error {
	var task entities.Task
	if err := tx.First(&task, updatedTask.Id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotExists
		}
		return err
	}
	return tx.Model(&task).Updates(updatedTask).Error
}

func (r *taskRepo) FindNotFinishedTasksThatDontHaveIncompletedTransactions(tx *gorm.DB) ([]entities.Task, error) {
	var tasks []entities.Task
	if err := tx.Raw("SELECT * FROM tasks t WHERE NOT EXISTS (SELECT * FROM transactions WHERE task_id = t.id AND status = ?) AND status = ?", common.TRANSACTION_RUNNING, common.TASK_RUNNING).Scan(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *taskRepo) FindNotFinishedAndExpired(tx *gorm.DB) ([]entities.Task, error) {
	var tasks []entities.Task
	if err := tx.Where("status = ?", common.TASK_RUNNING).Where("AND expiration < ?", time.Now()).Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}
