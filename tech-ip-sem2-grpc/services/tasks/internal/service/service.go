package service

import (
	"errors"
	"fmt"
	"sync"
)

var ErrTaskNotFound = errors.New("task not found")

type Task struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	DueDate     string `json:"due_date,omitempty"`
	Done        bool   `json:"done"`
}

type CreateTaskInput struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	DueDate     string `json:"due_date"`
}

type UpdateTaskInput struct {
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
	DueDate     *string `json:"due_date,omitempty"`
	Done        *bool   `json:"done,omitempty"`
}

type Service struct {
	mu      sync.RWMutex
	tasks   map[string]Task
	counter int
}

func New() *Service {
	return &Service{tasks: make(map[string]Task)}
}

func (s *Service) Create(input CreateTaskInput) (Task, error) {
	if input.Title == "" {
		return Task{}, errors.New("title is required")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.counter++
	id := fmt.Sprintf("t_%03d", s.counter)
	task := Task{
		ID:          id,
		Title:       input.Title,
		Description: input.Description,
		DueDate:     input.DueDate,
		Done:        false,
	}

	s.tasks[id] = task
	return task, nil
}

func (s *Service) List() []Task {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]Task, 0, len(s.tasks))
	for i := 1; i <= s.counter; i++ {
		id := fmt.Sprintf("t_%03d", i)
		task, ok := s.tasks[id]
		if ok {
			result = append(result, task)
		}
	}

	return result
}

func (s *Service) Get(id string) (Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	task, ok := s.tasks[id]
	if !ok {
		return Task{}, ErrTaskNotFound
	}

	return task, nil
}

func (s *Service) Update(id string, input UpdateTaskInput) (Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	task, ok := s.tasks[id]
	if !ok {
		return Task{}, ErrTaskNotFound
	}

	if input.Title != nil {
		if *input.Title == "" {
			return Task{}, errors.New("title cannot be empty")
		}
		task.Title = *input.Title
	}

	if input.Description != nil {
		task.Description = *input.Description
	}

	if input.DueDate != nil {
		task.DueDate = *input.DueDate
	}

	if input.Done != nil {
		task.Done = *input.Done
	}

	s.tasks[id] = task
	return task, nil
}

func (s *Service) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.tasks[id]; !ok {
		return ErrTaskNotFound
	}

	delete(s.tasks, id)
	return nil
}
