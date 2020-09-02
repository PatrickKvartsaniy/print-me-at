package repository

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v7"
	"github.com/sirupsen/logrus"
	"time"
)

type (
	Repository struct {
		redis           *redis.Client
		key             string
		pollingInterval time.Duration
	}

	RedisConfig struct {
		Addr            string
		Key             string
		PollingInterval time.Duration
	}
)

func NewRedisRepository(cfg RedisConfig) *Repository {
	client := redis.NewClient(&redis.Options{
		Addr: cfg.Addr,
	})
	return &Repository{redis: client}
}

func (r Repository) AddNewTask(ctx context.Context, msg string, ts time.Time) error {
	if err := r.redis.ZAdd(r.key, &redis.Z{Score: 0, Member: ts.Second()}, &redis.Z{Score: 1, Member: msg}).Err(); err != nil {
		return fmt.Errorf("adding msg to the queue: %w", err)
	}
	return nil
}

func (r Repository) Poll(ctx context.Context) error {
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
	res := r.redis.ZRange(r.key, 0, int64(time.Now().Second()))
	if err := res.Err(); err != nil {
		return fmt.Errorf("checking for the tasks: %w", err)
	}
	tasks, err := res.Result()
	if err != nil {
		return err
	}
	for _, task := range tasks {
		fmt.Println("scheduled message: " + task)
	}
	return nil
}

func (r Repository) Close() {
	if err := r.redis.Close(); err != nil {
		logrus.WithError(err).Error("closing redis connection")
	}
}

func (r Repository) HealthCheck() error {
	return r.redis.Ping().Err()
}
