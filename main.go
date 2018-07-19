package main

import (
  "github.com/sweettea-io/build-server/internal/app"
  "github.com/sweettea-io/build-server/internal/pkg/config"
)

func main() {
  // Initialize the app.
  app.Init(config.New())

  // Git clone build target.

  // Validate build target's config file.

  // Git clone buildpack.

  // Attach buildpack to target.

  // Build augmented target into Docker image.

  // Push Docker image to external repository.
}