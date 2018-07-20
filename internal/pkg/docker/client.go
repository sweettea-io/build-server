package docker

import (
  "github.com/docker/docker/client"
  "os"
  "context"
  "github.com/sweettea-io/build-server/internal/pkg/util/tar"
  "github.com/docker/docker/api/types"
)

var dockerClient *client.Client

// Init creates a new Docker client.
func Init(host string, apiVersion string, httpHeaders map[string]string) error {
  var err error
  dockerClient, err = client.NewClient(host, apiVersion, nil, httpHeaders)

  if err != nil {
    return err
  }

  return nil
}

// Build a Docker image from the specified directory containing a Dockerfile.
func Build(dir string, tag string) error {
  // Create tar file from build dir (ex. '/my/path' --> '/my/path.tar')
  buildContextPath, err := tar.Create(dir)

  if err != nil {
    return err
  }

  // Open created build context tar file and defer its closing.
  buildContext, err := os.Open(buildContextPath)
  defer buildContext.Close()

  // Create Docker build options
  buildOpts := types.ImageBuildOptions{
    SuppressOutput: false,
    Remove: true,
    ForceRemove: true,
    PullParent: true,
    Tags: []string{tag},
  }

  // Build the Docker image.
  if _, err := dockerClient.ImageBuild(context.Background(), buildContext, buildOpts); err != nil {
    return err
  }

  return nil
}

// Push a specified Docker image or repository to a registry.
func Push(name string) error {
  return nil
}