package targetconfig

import (
  "errors"
  "fmt"
  "io/ioutil"
  "strings"
  "github.com/sweettea-io/build-server/internal/pkg/util/fileutil"
  "gopkg.in/yaml.v2"
)

// ConfigFile is the name of the build target config file.
const FileName = ".sweettea.yml"

// Config is a struct representation of the contents of the config file.
type Config struct {
  Training struct {
    Buildpack string
    Dataset struct {
      Fetch  string
      Prepro string
    }
    Train string
    Test  string
    Eval  string
    Model struct {
      Path           string
      UploadCriteria *string
    }
  }
  Hosting struct {
    Buildpack string
    Predict   string
    Model struct {
      Path string
    }
  }
}

func New(targetDir string) (*Config, error) {
  // Build path to config file from provided directory.
  configPath := fmt.Sprintf("%s/%s", strings.TrimRight(targetDir, "/"), FileName)

  // Ensure config file exists.
  configExists, err := fileutil.Exists(configPath)

  if err != nil {
    return nil, err
  }

  if !configExists {
    return nil, errors.New("target config file not found")
  }

  // Read config file into byte array.
  configData, err := ioutil.ReadFile(configPath)

  if err != nil {
    return nil, err
  }

  // Unmarshal config file into struct.
  var config Config
  if err := yaml.Unmarshal(configData, &config); err != nil {
    return nil, err
  }

  return &config, nil
}