package buildpack

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

var validBuildpacks = map[string]bool {
  PythonTrainBuildpack: true,
  PythonJsonApiBuildpack: true,
}

// Check if buildpack name is supported.
func IsValidBuildpack(name string) bool {
  return validBuildpacks[name] == true
}