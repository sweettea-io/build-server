package buildpack

import "fmt"

// Buildpack struct representation
type Buildpack struct {
  Name        string
  Url         string
  Sha         string
  AccessToken string
  FileExt     string
}

// Supported Buildpacks
const (
  PythonTrainBuildpack = "python-train"
  PythonJsonApiBuildpack = "python-json-api"
)

// Validate buildpack is supported for provided name.
func Validate(name string) error {
  switch name {
  case PythonTrainBuildpack, PythonJsonApiBuildpack:
    return nil
  default:
    return fmt.Errorf("buildpack \"%s\" not supported", name)
  }
}