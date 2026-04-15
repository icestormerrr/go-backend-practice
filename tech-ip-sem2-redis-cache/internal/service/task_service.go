package service

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"example.com/tech-ip-sem2-redis-cache/internal/cache"
	"example.com/tech-ip-sem2-redis-cache/internal/config"
	"example.com/tech-ip-sem2-redis-cache/internal/task"
	"github.com/redis/go-redis/v9"
)

type TaskService struct {
	repo  *task.Repo
	redis *redis.Client
	cfg   config.Config
}

func NewTaskService(repo *task.Repo, redisClient *redis.Client, cfg config.Config) *TaskService {
	return &TaskService{
		repo:  repo,
		redis: redisClient,
		cfg:   cfg,
	}
}

func (s *TaskService) ListTasks(ctx context.Context) ([]task.Task, error) {
	key := cache.TasksListKey()

	if s.redis != nil {
		cached, err := s.redis.Get(ctx, key).Result()
		if err == nil {
			var tasks []task.Task
			if err := json.Unmarshal([]byte(cached), &tasks); err == nil {
				log.Println("cache hit:", key)
				return tasks, nil
			}

			log.Println("cache decode error:", err)
		} else if !errors.Is(err, redis.Nil) {
			log.Println("redis read error:", err)
		} else {
			log.Println("cache miss:", key)
		}
	}

	tasks := s.repo.List()

	if s.redis != nil {
		bytes, err := json.Marshal(tasks)
		if err != nil {
			log.Println("cache encode error:", err)
			return tasks, nil
		}

		ttl := cache.TTLWithJitter(s.cfg.CacheTTL, s.cfg.CacheTTLJitter)
		if err := s.redis.Set(ctx, key, bytes, ttl).Err(); err != nil {
			log.Println("redis write error:", err)
		} else {
			log.Println("cache set:", key, "ttl:", ttl)
		}
	}

	return tasks, nil
}

func (s *TaskService) GetTaskByID(ctx context.Context, id int64) (task.Task, error) {
	key := cache.TaskByIDKey(id)

	if s.redis != nil {
		cached, err := s.redis.Get(ctx, key).Result()
		if err == nil {
			var t task.Task
			if err := json.Unmarshal([]byte(cached), &t); err == nil {
				log.Println("cache hit:", key)
				return t, nil
			}

			log.Println("cache decode error:", err)
		} else if !errors.Is(err, redis.Nil) {
			log.Println("redis read error:", err)
		} else {
			log.Println("cache miss:", key)
		}
	}

	t, err := s.repo.GetByID(id)
	if err != nil {
		return task.Task{}, err
	}

	if s.redis != nil {
		bytes, err := json.Marshal(t)
		if err != nil {
			log.Println("cache encode error:", err)
			return t, nil
		}

		ttl := cache.TTLWithJitter(s.cfg.CacheTTL, s.cfg.CacheTTLJitter)
		if err := s.redis.Set(ctx, key, bytes, ttl).Err(); err != nil {
			log.Println("redis write error:", err)
		} else {
			log.Println("cache set:", key, "ttl:", ttl)
		}
	}

	return t, nil
}

func (s *TaskService) UpdateTask(ctx context.Context, t task.Task) error {
	if err := s.repo.Update(t); err != nil {
		return err
	}

	if s.redis != nil {
		keys := []string{cache.TaskByIDKey(t.ID), cache.TasksListKey()}
		if err := s.redis.Del(ctx, keys...).Err(); err != nil {
			log.Println("redis delete error:", err)
		} else {
			for _, key := range keys {
				log.Println("cache invalidated:", key)
			}
		}
	}

	return nil
}

func (s *TaskService) DeleteTask(ctx context.Context, id int64) error {
	if err := s.repo.Delete(id); err != nil {
		return err
	}

	if s.redis != nil {
		keys := []string{cache.TaskByIDKey(id), cache.TasksListKey()}
		if err := s.redis.Del(ctx, keys...).Err(); err != nil {
			log.Println("redis delete error:", err)
		} else {
			for _, key := range keys {
				log.Println("cache invalidated:", key)
			}
		}
	}

	return nil
}
