package services

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"Signup_login/dto"
	"Signup_login/model"

	"gorm.io/gorm"
)

var (
	ErrInvalidTaskInput = errors.New("invalid task input")
	ErrTaskNotFound     = errors.New("task not found")
)

type TaskStore interface {
	Create(task *model.Task) error
	ListByUser(userID uint) ([]model.Task, error)
	FindByIDForUser(userID, taskID uint) (*model.Task, error)
	UpdateForUser(task *model.Task) (*model.Task, error)
	DeleteForUser(userID, taskID uint) error
}

type TaskService struct {
	Repo TaskStore
}

func NewTaskService(repo TaskStore) *TaskService {
	return &TaskService{Repo: repo}
}

func (s *TaskService) CreateTask(userID uint, req dto.TaskRequest) (*model.Task, error) {
	task, err := buildTask(userID, req)
	if err != nil {
		return nil, err
	}
	if err := s.Repo.Create(task); err != nil {
		return nil, err
	}
	return task, nil
}

func (s *TaskService) ListTasks(userID uint) ([]model.Task, error) {
	return s.Repo.ListByUser(userID)
}

func (s *TaskService) GetTask(userID, taskID uint) (*model.Task, error) {
	task, err := s.Repo.FindByIDForUser(userID, taskID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrTaskNotFound
	}
	return task, err
}

func (s *TaskService) UpdateTask(userID, taskID uint, req dto.TaskRequest) (*model.Task, error) {
	task, err := buildTask(userID, req)
	if err != nil {
		return nil, err
	}
	task.ID = taskID

	updated, err := s.Repo.UpdateForUser(task)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrTaskNotFound
	}
	return updated, err
}

func (s *TaskService) DeleteTask(userID, taskID uint) error {
	err := s.Repo.DeleteForUser(userID, taskID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrTaskNotFound
	}
	return err
}

func buildTask(userID uint, req dto.TaskRequest) (*model.Task, error) {
	subject := strings.TrimSpace(req.Subject)
	if subject == "" {
		return nil, invalidTaskInput("subject is required")
	}
	topic := strings.TrimSpace(req.Topic)
	if topic == "" {
		return nil, invalidTaskInput("topic is required")
	}
	description := strings.TrimSpace(req.Description)
	if description == "" {
		return nil, invalidTaskInput("description is required")
	}
	date := strings.TrimSpace(req.Date)
	parsedDate, err := time.Parse("2006-01-02", date)
	if err != nil || parsedDate.Format("2006-01-02") != date {
		return nil, invalidTaskInput("date must use YYYY-MM-DD format")
	}
	if req.Duration <= 0 {
		return nil, invalidTaskInput("duration must be a positive number of minutes")
	}

	priority, err := normalizePriority(req.Priority)
	if err != nil {
		return nil, err
	}
	status, err := normalizeStatus(req.Status)
	if err != nil {
		return nil, err
	}

	return &model.Task{
		UserID:      userID,
		Subject:     subject,
		Topic:       topic,
		Description: description,
		Date:        date,
		Duration:    req.Duration,
		Priority:    priority,
		Status:      status,
	}, nil
}

func normalizePriority(value string) (string, error) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "high":
		return "High", nil
	case "medium":
		return "Medium", nil
	case "low":
		return "Low", nil
	default:
		return "", invalidTaskInput("priority must be High, Medium, or Low")
	}
}

func normalizeStatus(value string) (string, error) {
	switch strings.ToLower(strings.Join(strings.Fields(value), " ")) {
	case "not started":
		return "Not Started", nil
	case "in progress":
		return "In Progress", nil
	case "completed":
		return "Completed", nil
	case "on hold":
		return "On Hold", nil
	default:
		return "", invalidTaskInput("status must be Not Started, In Progress, Completed, or On Hold")
	}
}

func invalidTaskInput(message string) error {
	return fmt.Errorf("%w: %s", ErrInvalidTaskInput, message)
}
