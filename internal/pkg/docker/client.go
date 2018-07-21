package docker

import (
  "bufio"
  "context"
  "errors"
  "encoding/json"
  "encoding/base64"
  "fmt"
  "io"
  "os"
  "github.com/docker/docker/api/types"
  "github.com/docker/docker/client"
  "github.com/sweettea-io/build-server/internal/pkg/util/tar"
)

// Wrapper type around Docker *client.Client, providing custom build/push functionality.
type Client struct {
  client *client.Client
  auth   string
}

// Build a Docker image from the specified directory containing a Dockerfile.
func (dc *Client) Build(dir string, tag string, buildArgs map[string]*string) error {
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
  resp, err := dc.client.ImageBuild(context.Background(), buildContext, buildOpts)

  if err != nil {
    return err
  }

  // Read the response.
  if err := readRespLines(resp.Body); err != nil {
    return err
  }

  return nil
}

// Push a specified Docker image or repository to a registry.
func (dc *Client) Push(name string) error {
  // Create Docker push options.
  pushOpts := types.ImagePushOptions{
    RegistryAuth: dc.auth,
  }

  // Push image.
  output, err := dc.client.ImagePush(context.Background(), name, pushOpts)

  if err != nil {
    return err
  }

  // Read the response.
  if err := readRespLines(output); err != nil {
    return err
  }

  return nil
}

// New creates and returns a new `Client` instance pointer.
func New(host string, apiVersion string, httpHeaders map[string]string, username string, password string) (*Client, error) {
  // Create a new Docker client.
  internalClient, err := client.NewClient(host, apiVersion, nil, httpHeaders)

  if err != nil {
    return nil, err
  }

  // Marshal JSON auth config.
  jsonAuth, err := json.Marshal(types.AuthConfig{
    Username: username,
    Password: password,
  })

  if err != nil {
    return nil, err
  }

  // Create new `Client` instance.
  dc := &Client{
    client: internalClient,
    auth: base64.URLEncoding.EncodeToString(jsonAuth),
  }

  return dc, nil
}

// readRespLines reads and JSON-parses lines from a Docker API response body
// and returns an error if an "error" key is found in the line.
func readRespLines(body io.ReadCloser) error {
  // Pass response body into new reader.
  reader := bufio.NewReader(body)

  for {
    // Read a line of bytes from the response.
    lineBytes, err := reader.ReadBytes('\n')

    if err != nil {
      // Return successfully when EOF reached.
      if err == io.EOF {
        return nil
      }

      return err
    }

    // Log the line as a string.
    fmt.Print(string(lineBytes[:]))

    // Parse the line into JSON.
    var lineJSON map[string]string
    json.Unmarshal(lineBytes, &lineJSON)

    // Check if this line indicates an error.
    if lineJSON["error"] != "" {
      return errors.New(lineJSON["error"])
    }
  }
}