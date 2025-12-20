package task

import (
	"encoding/json"
	"errors"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

var ErrNotFound = errors.New("task not found")

type Repo struct {
	mu       sync.RWMutex
	filePath string
}

func NewRepo(filePath string) *Repo {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		_ = os.WriteFile(filePath, []byte("{}"), 0644)
	}
	return &Repo{filePath: filePath}
}

// readAll читает все задачи из файла и возвращает map[id]Task
func (r *Repo) readAll() (map[string]Task, error) {
	data, err := os.ReadFile(r.filePath)
	if err != nil {
		return nil, err
	}

	tasks := make(map[string]Task)
	if len(data) == 0 {
		return tasks, nil
	}

	if err := json.Unmarshal(data, &tasks); err != nil {
		return nil, err
	}

	return tasks, nil
}

// writeAll записывает map[id]Task в файл
func (r *Repo) writeAll(tasks map[string]Task) error {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(r.filePath, data, 0644)
}

func (r *Repo) List(title string, page, limit int) ([]Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tasksMap, err := r.readAll()
	if err != nil {
		return nil, err
	}

	tasks := make([]Task, 0, len(tasksMap))
	for _, t := range tasksMap {
		if title == "" || strings.Contains(strings.ToLower(t.Title), strings.ToLower(title)) {
			tasks = append(tasks, t)
		}
	}

	start := (page - 1) * limit
	if start > len(tasks) {
		return []Task{}, nil
	}
	end := start + limit
	if end > len(tasks) {
		end = len(tasks)
	}

	return tasks[start:end], nil
}

func (r *Repo) Get(id string) (*Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tasks, err := r.readAll()
	if err != nil {
		return nil, err
	}

	t, ok := tasks[id]
	if !ok {
		return nil, ErrNotFound
	}
	return &t, nil
}

func (r *Repo) Create(title string) (*Task, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	tasks, err := r.readAll()
	if err != nil {
		return nil, err
	}

	now := time.Now()
	t := Task{
		ID:        uuid.NewString(),
		Title:     title,
		CreatedAt: now,
		UpdatedAt: now,
		Done:      false,
	}

	tasks[t.ID] = t

	if err := r.writeAll(tasks); err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *Repo) Update(id, title string, done bool) (*Task, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	tasks, err := r.readAll()
	if err != nil {
		return nil, err
	}

	t, ok := tasks[id]
	if !ok {
		return nil, ErrNotFound
	}

	t.Title = title
	t.Done = done
	t.UpdatedAt = time.Now()
	tasks[id] = t

	if err := r.writeAll(tasks); err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *Repo) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	tasks, err := r.readAll()
	if err != nil {
		return err
	}

	if _, ok := tasks[id]; !ok {
		return ErrNotFound
	}

	delete(tasks, id)

	return r.writeAll(tasks)
}
