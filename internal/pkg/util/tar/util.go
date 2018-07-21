package tar

import (
  "strings"
  "github.com/jhoonb/archivex"
)

// Create a new tarfile from the provided directory. The path to
// the tarfile will be returned, along with any error that occurs.
func Create(dir string) (string, error) {
  tarfile := new(archivex.TarFile)

  // Create tarfile path from the dir path provided.
  tarfilePath := strings.TrimRight(dir, "/") + ".tar"

  // Create new, empty tarfile at path.
  if err := tarfile.Create(tarfilePath); err != nil {
    return "", err
  }

  // Add all files/dirs from the provided dir into the tarfile.
  if err := tarfile.AddAll(dir, false); err != nil {
    return "", err
  }

  // Close the tarfile.
  if err := tarfile.Close(); err != nil {
    return "", err
  }

  return tarfilePath, nil
}
