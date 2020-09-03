package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/PatrickKvartsaniy/print-me-at/models"
	"github.com/go-redis/redis/v7"
	"github.com/sirupsen/logrus"
	"strconv"
	"time"
)

type (
	Repository struct {
		redis           *redis.Client
		key             string
		pollingInterval time.Duration
		runErr          error
	}

	RedisConfig struct {
		Addr            string
		Key             string
		PollingInterval time.Duration
	}
)

func NewRedisRepository(cfg RedisConfig) *Repository {
	client := redis.NewClient(&redis.Options{
		Addr:       cfg.Addr,
		MaxRetries: 3,
	})
	return &Repository{
		redis:           client,
		key:             cfg.Key,
		pollingInterval: cfg.PollingInterval,
	}
}

func (r Repository) AddNewTask(m models.Message) error {
	msg, err := m.ToJSONString()
	if err != nil {
		return fmt.Errorf("marshaling message: %w", err)
	}
	if err := r.redis.ZAdd(r.key, &redis.Z{Score: float64(m.Timestamp), Member: msg}).Err(); err != nil {
		return fmt.Errorf("adding msg to the queue: %w", err)
	}
	return nil
}

func (r Repository) Poll(ctx context.Context) {
	logrus.Info("start polling")
	go func() {
		if err := r.poll(ctx); err != nil {
			r.runErr = err
			logrus.WithError(err).Error("redis service")
		}
	}()
}

func (r Repository) poll(ctx context.Context) error {
	ticker := time.NewTicker(r.pollingInterval)
	for {
		select {
		case <-ticker.C:
			if err := r.checkForTasks(); err != nil {
				return err
			}
		case <-ctx.Done():
			return nil
		}
	}
}

func (r Repository) checkForTasks() error {
	res := r.redis.ZRangeByScore(r.key, &redis.ZRangeBy{Min: "0", Max: strconv.Itoa(int(time.Now().Unix()))})
	if err := res.Err(); err != nil {
		return fmt.Errorf("checking for the tasks: %w", err)
	}
	tasks, err := res.Result()
	if err != nil {
		return err
	}
	return r.processTask(tasks)
}

func (r Repository) processTask(tasks []string) error {
	for _, task := range tasks {
		var msg models.Message
		if err := json.Unmarshal([]byte(task), &msg); err != nil {
			return fmt.Errorf("unmarshaling message: %w", err)
		}
		msg.PrintOut()
		if err := r.redis.ZRem(r.key, task).Err(); err != nil {
			return fmt.Errorf("message has been received but can't be deleted: %w", err)
		}
	}
	return nil
}

func (r Repository) Close() {
	if err := r.redis.Close(); err != nil {
		logrus.WithError(err).Error("closing redis connection")
	}
}

func (r Repository) HealthCheck() error {
	if err := r.redis.Ping().Err(); err != nil {
		return fmt.Errorf("redis is not responding")
	}
	if r.runErr != nil {
		return r.runErr
	}
	return nil
}
