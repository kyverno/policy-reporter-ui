package redis

import (
	"context"
	"encoding/json"
	"log"

	goredis "github.com/go-redis/redis/v8"

	"github.com/kyverno/policy-reporter-ui/pkg/report"
)

// Redis configuration
type Config struct {
	Enabled  bool   `mapstructure:"enabled"`
	Address  string `mapstructure:"address"`
	Prefix   string `mapstructure:"prefix"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Database int    `mapstructure:"database"`
}

type RedisStore struct {
	rdb *goredis.Client
	key string
}

func (r *RedisStore) Add(report report.Result) error {
	content, err := json.Marshal(report)
	if err != nil {
		return err
	}

	return r.rdb.LPush(context.Background(), r.key, string(content)).Err()
}

func (r *RedisStore) List() ([]report.Result, error) {
	list := make([]report.Result, 0)
	res, err := r.rdb.LRange(context.Background(), r.key, 0, -1).Result()
	if err != nil {
		return list, err
	}

	for _, i := range res {
		result := report.Result{}

		err = json.Unmarshal([]byte(i), &result)
		if err != nil {
			log.Printf("[WARNING] failed to unmarshel result: %s", err)
			continue
		}

		list = append(list, result)
	}

	return list, nil
}

func New(key string, rdb *goredis.Client) *RedisStore {
	return &RedisStore{rdb, key}
}

func NewFromConfig(conf Config) *RedisStore {
	return New(conf.Prefix+":results", goredis.NewClient(&goredis.Options{
		Addr:     conf.Address,
		Username: conf.Username,
		Password: conf.Password,
		DB:       conf.Database,
	}))
}
