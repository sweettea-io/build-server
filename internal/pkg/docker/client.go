package docker

import (
  "github.com/docker/docker/client"
  "os"
  "context"
  "github.com/sweettea-io/build-server/internal/pkg/util/tar"
  "github.com/docker/docker/api/types"
  "io"
  "encoding/json"
  "encoding/base64"
)

var dockerClient *Client

type Client struct {
  client *client.Client
  auth   string
}

// Init assigns a new Docker `Client` instance to the global `dockerClient`.
func Init(host string, apiVersion string, httpHeaders map[string]string, username string, password string) error {
  // Create a new Docker client.
  dc, err := client.NewClient(host, apiVersion, nil, httpHeaders)

  if err != nil {
    return err
  }

  // Marshal JSON auth config.
  jsonAuth, err := json.Marshal(types.AuthConfig{
    Username: username,
    Password: password,
  })

  if err != nil {
    return err
  }

  // Create new `DockerClient` instance.
  dockerClient = &Client{
    client: dc,
    auth: base64.URLEncoding.EncodeToString(jsonAuth),
  }

  return nil
}

// Build a Docker image from the specified directory containing a Dockerfile.
func Build(dir string, tag string, buildArgs map[string]*string) error {
  // Create tar file from build dir (ex. '/my/path' --> '/my/path.tar').
  buildContextPath, err := tar.Create(dir)

  if err != nil {
    return err
  }

  // Open created build context tar file and defer its closing.
  buildContext, err := os.Open(buildContextPath)
  defer buildContext.Close()

  // Create Docker build options.
  buildOpts := types.ImageBuildOptions{
    SuppressOutput: false,
    Remove: true,
    ForceRemove: true,
    PullParent: true,
    Tags: []string{tag},
    BuildArgs: buildArgs,
  }

  // Build the Docker image.
  resp, err := dockerClient.client.ImageBuild(context.Background(), buildContext, buildOpts)

  if err != nil {
    return err
  }

  // Copy log output from push command to stdout.
  io.Copy(os.Stdout, resp.Body)

  return nil
}

// Push a specified Docker image or repository to a registry.
func Push(name string, ) error {
  // Create Docker push options.
  pushOpts := types.ImagePushOptions{
    RegistryAuth: dockerClient.auth,
  }

  // Push image.
  output, err := dockerClient.client.ImagePush(context.Background(), name, pushOpts)

  if err != nil {
    return err
  }

  // Copy log output from push command to stdout.
  io.Copy(os.Stdout, output)

  return nil
}