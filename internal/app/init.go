package app

import (
  "github.com/Sirupsen/logrus"
  "github.com/sweettea-io/build-server/internal/pkg/config"
  "github.com/sweettea-io/build-server/internal/pkg/logger"
  "github.com/sweettea-io/build-server/internal/pkg/redis"
  r "github.com/gomodule/redigo/redis"
)

var Config *config.Config
var Redis *r.Pool
var Log *logrus.Logger

func Init(cfg *config.Config) {
  // Set global config.
  Config = cfg

  // Create redis pool.
  Redis = redis.NewPool(
    Config.RedisAddress,
    Config.RedisPassword,
    Config.RedisPoolMaxActive,
    Config.RedisPoolMaxIdle,
    Config.RedisPoolWait,
  )

  // Create app logger.
  Log = logger.NewLogger()
}