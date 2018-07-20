package tar

import (
  "strings"
  "github.com/jhoonb/archivex"
)

// Create a new tar file from a directory.
func Create(dir string) (string, error) {
  tar := new(archivex.TarFile)

  if err := tar.Create(dir); err != nil {
    return "", err
  }

  if err := tar.AddAll(dir, false); err != nil {
    return "", err
  }

  if err := tar.Close(); err != nil {
    return "", err
  }

  return strings.TrimRight(dir, "/") + ".tar", nil
}
