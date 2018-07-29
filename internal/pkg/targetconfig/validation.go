package targetconfig

import (
  "github.com/sweettea-io/build-server/internal/pkg/util/targetutil"
  "fmt"
  "github.com/sweettea-io/build-server/internal/pkg/util/buildpack"
  "github.com/sweettea-io/build-server/internal/pkg/util/model"
)

func (c *Config) Validate(targetCluster string) error {
  switch targetCluster {
  case targetutil.Train:
    return c.ValidateTraining()
  case targetutil.API:
    return c.ValidateHosting()
  default:
    return fmt.Errorf("target cluster \"%s\" unknown", targetCluster)
  }
}

func (c *Config) ValidateTraining() error {
  // Validate buildpack.
  if err := buildpack.Validate(c.Training.Buildpack); err != nil {
    return err
  }

  // Validate train method.
  if err := validateNonEmptyStr(c.Training.Train, "training.train"); err != nil {
    return err
  }

  // Validate model path.
  if err := validateNonEmptyStr(c.Training.Model.Path, "training.model.path"); err != nil {
    return err
  }

  // No more checks needed if model upload criteria not provided.
  if c.Training.Model.UploadCriteria == nil {
    return nil
  }

  // Validate model upload criteria.
  if err := model.ValidateUploadCriteria(*c.Training.Model.UploadCriteria); err != nil {
    return err
  }

  // If model upload criteria is based on the success of the
  // eval method, ensure the eval method is present.
  if *c.Training.Model.UploadCriteria == model.UploadCriteria.Eval {
    if err := validateNonEmptyStr(c.Training.Train, "training.eval"); err != nil {
      return err
    }
  }

  return nil
}

func (c *Config) ValidateHosting() error {
  // Validate buildpack.
  if err := buildpack.Validate(c.Hosting.Buildpack); err != nil {
    return err
  }

  // Validate predict method.
  if err := validateNonEmptyStr(c.Hosting.Predict, "hosting.predict"); err != nil {
    return err
  }

  // Validate model path.
  if err := validateNonEmptyStr(c.Hosting.Model.Path, "hosting.model.path"); err != nil {
    return err
  }

  return nil
}

func validateNonEmptyStr(val string, keyPath string) error {
  if val == "" {
    return fmt.Errorf("key at path \"%s\" cannot be empty", keyPath)
  }

  return nil
}