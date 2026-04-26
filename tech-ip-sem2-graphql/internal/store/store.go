package store

import "sync"

type Task struct {
	ID          string
	Title       string
	Description *string
	Done        bool
}

type Store struct {
	mu    sync.RWMutex
	tasks []*Task
}

func New() *Store {
	return &Store{
		tasks: []*Task{
			{ID: "t_001", Title: "Первая задача", Description: strPtr("Учебный пример"), Done: false},
			{ID: "t_002", Title: "Вторая задача", Description: strPtr("GraphQL API"), Done: true},
		},
	}
}

func strPtr(s string) *string {
	return &s
}

func (s *Store) List() []*Task {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]*Task, 0, len(s.tasks))
	for _, t := range s.tasks {
		copyTask := *t
		result = append(result, &copyTask)
	}

	return result
}

func (s *Store) GetByID(id string) (*Task, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, t := range s.tasks {
		if t.ID == id {
			copyTask := *t
			return &copyTask, true
		}
	}

	return nil, false
}

func (s *Store) Create(task *Task) *Task {
	s.mu.Lock()
	defer s.mu.Unlock()

	copyTask := *task
	s.tasks = append(s.tasks, &copyTask)
	return &copyTask
}

func (s *Store) CreateTask(id, title string, description *string) *Task {
	return s.Create(&Task{
		ID:          id,
		Title:       title,
		Description: description,
		Done:        false,
	})
}

func (s *Store) Update(id string, updateFn func(task *Task)) (*Task, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, t := range s.tasks {
		if t.ID == id {
			updateFn(t)
			copyTask := *t
			return &copyTask, true
		}
	}

	return nil, false
}

func (s *Store) UpdateTask(id string, updateFn func(title *string, description **string, done *bool)) (*Task, bool) {
	return s.Update(id, func(task *Task) {
		updateFn(&task.Title, &task.Description, &task.Done)
	})
}

func (s *Store) Delete(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, t := range s.tasks {
		if t.ID == id {
			s.tasks = append(s.tasks[:i], s.tasks[i+1:]...)
			return true
		}
	}

	return false
}

func (s *Store) DeleteTask(id string) bool {
	return s.Delete(id)
}
