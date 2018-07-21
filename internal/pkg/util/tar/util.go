package tar

import (
  "strings"
  "github.com/jhoonb/archivex"
)

// Create a new tar file from a directory.
func Create(dir string) (string, error) {
  tarfile := new(archivex.TarFile)
  tarfilePath := strings.TrimRight(dir, "/") + ".tar"

  if err := tarfile.Create(tarfilePath); err != nil {
    return "", err
  }

  if err := tarfile.AddAll(dir, false); err != nil {
    return "", err
  }

  if err := tarfile.Close(); err != nil {
    return "", err
  }

  return tarfilePath, nil
}
