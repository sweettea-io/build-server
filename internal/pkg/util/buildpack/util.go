package buildpack

import (
  "github.com/sweettea-io/build-server/internal/pkg/config"
  "fmt"
)

func FromConfig(cfg *config.Config, bpName string) (*Buildpack, error) {
  // Ensure buildpack name is valid.
  if !IsValidBuildpack(bpName) {
    return nil, fmt.Errorf("invalid buildpack: \"%s\"", bpName)
  }

  // Return configured Buildpack instance by name.
  switch bpName {
  case PythonTrainBuildpack:
    return &Buildpack{
      Name: PythonTrainBuildpack,
      Url: cfg.PythonTrainBuildpackUrl,
      Sha: cfg.PythonTrainBuildpackSha,
      AccessToken: cfg.PythonTrainBuildpackAccessToken,
    }, nil
  case PythonJsonApiBuildpack:
    return &Buildpack{
      Name: PythonJsonApiBuildpack,
      Url: cfg.PythonJsonApiBuildpackUrl,
      Sha: cfg.PythonJsonApiBuildpackSha,
      AccessToken: cfg.PythonJsonApiBuildpackAccessToken,
    }, nil
  default:
    return nil, fmt.Errorf("can't configure buildpack for name: \"%s\"", bpName)
  }
}