package main

import (
  "github.com/sweettea-io/build-server/internal/pkg/config"
  "github.com/sweettea-io/build-server/internal/pkg/redis"
  "github.com/sweettea-io/build-server/internal/pkg/logger"
)

func main() {
  // Get new config instance.
  cfg := config.New()

  // Create new redis pool.
  pool := redis.NewPool(
    cfg.RedisAddress,
    cfg.RedisPassword,
    cfg.RedisPoolMaxActive,
    cfg.RedisPoolMaxIdle,
    cfg.RedisPoolWait,
  )

  // Create new logger.
  log := logger.New(pool, cfg.LogStreamKey())

  // Git clone build target.
  log.Infof("Cloning %s...", cfg.BuildTargetUrl)

  // Validate build target's config file.

  // Git clone buildpack.

  // Attach buildpack to target.

  // Build augmented target into Docker image.

  // Push Docker image to external repository.
}