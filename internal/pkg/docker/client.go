package docker

import (
  "github.com/docker/docker/client"
  "os"
  "context"
  "github.com/sweettea-io/build-server/internal/pkg/util/tar"
  "github.com/docker/docker/api/types"
)

var dockerClient *client.Client

func Init(host string, apiVersion string, httpHeaders map[string]string) error {
  var err error
  dockerClient, err = client.NewClient(host, apiVersion, nil, httpHeaders)

  if err != nil {
    return err
  }

  return nil
}

func Build(dir string, tag string) error {
  // Create tar file from build dir (ex. '/my/path' --> '/my/path.tar')
  buildContextPath, err := tar.Create(dir)

  if err != nil {
    return err
  }

  // Open created build context tar file and defer its closing.
  buildContext, err := os.Open(buildContextPath)
  defer buildContext.Close()

  // Create docker build options
  buildOpts := types.ImageBuildOptions{
    SuppressOutput: false,
    Remove: true,
    ForceRemove: true,
    PullParent: true,
    Tags: []string{tag},
  }

  // Build the docker image.
  if _, err := dockerClient.ImageBuild(context.Background(), buildContext, buildOpts); err != nil {
    return err
  }

  return nil
}