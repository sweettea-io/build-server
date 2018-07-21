package targetutil

import (
  "errors"
  "fmt"
  "io/ioutil"
  "strings"
  "github.com/sweettea-io/build-server/internal/pkg/util/fileutil"
  "gopkg.in/yaml.v2"
)

// ConfigFile is the name of the build target config file.
const ConfigFile = ".sweettea.yml"

// Config is a struct representation of the contents of the config file.
type Config struct {
  Key string
}

// ValidateConfig ensures the build target config file exists inside the
// provided directory, and that its contents make for a valid SweetTea application.
func ValidateConfig(dir string) error {
  // Build path to config file from provided directory.
  configPath := fmt.Sprintf("%s/%s", strings.TrimRight(dir, "/"), ConfigFile)

  // Ensure config file exists.
  configExists, err := fileutil.Exists(configPath)

  if err != nil {
    return err
  }

  if !configExists {
    return errors.New("target config file not found")
  }

  // Read config file into byte array.
  configData, err := ioutil.ReadFile(configPath)

  if err != nil {
    return err
  }

  // Unmarshal config file into struct.
  var config Config
  if err := yaml.Unmarshal(configData, &config); err != nil {
    return err
  }

  return nil
}