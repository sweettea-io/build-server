package buildpack

// Supported Buildpacks.
const (
  PythonTrainBuildpack = "python-train"
  PythonJsonApiBuildpack = "python-json-api"
)

var validBuildpacks = map[string]bool {
  PythonTrainBuildpack: true,
  PythonJsonApiBuildpack: true,
}

func IsValidBuildpack(name string) bool {
  return validBuildpacks[name] == true
}