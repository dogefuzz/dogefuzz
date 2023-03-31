package repo

import (
	"errors"

	"github.com/dogefuzz/dogefuzz/entities"
	"github.com/dogefuzz/dogefuzz/pkg/common"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type taskRepo struct {
}

func NewTaskRepo(e Env) *taskRepo {
	return &taskRepo{}
}

func (r *taskRepo) Get(tx *gorm.DB, id string) (*entities.Task, error) {
	var task entities.Task
	if err := tx.First(&task, "id = ?", id).Error; err != nil {
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
	if err := tx.First(&task, "id = ?", updatedTask.Id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotExists
		}
		return err
	}
	return tx.Model(&task).Updates(updatedTask).Error
}

func (r *taskRepo) FindNotFinishedTasksThatDontHaveIncompletedTransactions(tx *gorm.DB) ([]entities.Task, error) {
	var tasks []entities.Task
	if err := tx.Raw("SELECT * FROM tasks t WHERE NOT EXISTS (SELECT * FROM transactions tx WHERE tx.task_id = t.id AND tx.status = ?) AND t.status = ?", common.TRANSACTION_RUNNING, common.TASK_RUNNING).Scan(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *taskRepo) FindNotFinishedAndExpired(tx *gorm.DB) ([]entities.Task, error) {
	var tasks []entities.Task
	if err := tx.Where("status = ?", common.TASK_RUNNING).Where("expiration < ?", common.Now()).Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *taskRepo) FindNotFinishedAndHaveDeployedContract(tx *gorm.DB) ([]entities.Task, error) {
	var tasks []entities.Task
	if err := tx.Raw("SELECT * FROM tasks t WHERE EXISTS (SELECT * FROM contracts c WHERE c.task_id = t.id AND c.status = ?) AND t.status = ?", common.CONTRACT_DEPLOYED, common.TASK_RUNNING).Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *taskRepo) FindNotFinishedThatHaveDeployedContractAndLimitedPendingTransactions(tx *gorm.DB, limit int) ([]entities.Task, error) {
	var tasks []entities.Task
	query := `
		SELECT *
		FROM tasks t
		WHERE
			EXISTS (SELECT *
				FROM contracts c
				WHERE c.task_id = t.id
					AND c.status = ?)
			AND EXISTS (SELECT *
			 	FROM transactions tx
				WHERE tx.task_id = t.id
					AND tx.status = ?
				GROUP BY tx.status
				HAVING COUNT(*) <= ?)
			AND t.status = ?`
	if err := tx.Raw(query, common.CONTRACT_DEPLOYED, common.TRANSACTION_RUNNING, limit, common.TASK_RUNNING).Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}
