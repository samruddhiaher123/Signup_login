package repository

import (
	"Signup_login/model"

	"gorm.io/gorm"
)

type TaskRepository struct {
	DB *gorm.DB
}

func NewTaskRepository(db *gorm.DB) *TaskRepository {
	return &TaskRepository{DB: db}
}
func (r *TaskRepository) Create(task *model.Task) error {
	return r.DB.Create(task).Error
}

func (r *TaskRepository) ListByUser(userID uint) ([]model.Task, error) {
	tasks := make([]model.Task, 0)
	err := r.DB.Where("user_id = ?", userID).Order("id ASC").Find(&tasks).Error
	return tasks, err
}

func (r *TaskRepository) FindByIDForUser(userID, taskID uint) (*model.Task, error) {
	var task model.Task
	err := r.DB.Where("id = ? AND user_id = ?", taskID, userID).First(&task).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *TaskRepository) UpdateForUser(task *model.Task) (*model.Task, error) {
	result := r.DB.Model(&model.Task{}).
		Where("id = ? AND user_id = ?", task.ID, task.UserID).
		Updates(map[string]any{
			"subject":     task.Subject,
			"topic":       task.Topic,
			"description": task.Description,
			"date":        task.Date,
			"duration":    task.Duration,
			"priority":    task.Priority,
			"status":      task.Status,
		})
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return r.FindByIDForUser(task.UserID, task.ID)
}

func (r *TaskRepository) DeleteForUser(userID, taskID uint) error {
	result := r.DB.Where("id = ? AND user_id = ?", taskID, userID).
		Delete(&model.Task{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
