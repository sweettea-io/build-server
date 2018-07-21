package targetutil

import (
  "errors"
  "fmt"
  "io/ioutil"
  "strings"
  "github.com/sweettea-io/build-server/internal/pkg/util/fileutil"
  "gopkg.in/yaml.v2"
)

const ConfigFile = ".sweettea.yml"

type Config struct {
  Key string
}

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

  fmt.Println(config.Key)

  return nil
}