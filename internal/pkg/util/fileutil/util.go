package fileutil

import "os"

// Exists returns a boolean indicating whether a
// file or directory exists at the given path.
func Exists(path string) (bool, error) {
  _, err := os.Stat(path)

  if err == nil {
    return true, nil
  }

  if os.IsNotExist(err) {
    return false, nil
  }

  return true, err
}

// RemoveIfExists removes a file or directory (and all of its contents) if it exists.
func RemoveIfExists(path string) error {
  // Check if a file/dir exists at the provided path.
  exists, err := Exists(path)

  if err != nil {
    return err
  }

  if !exists {
    return nil
  }

  // Remove the file or dir (and all sub-directories if a dir).
  if err := os.RemoveAll(path); err != nil {
    return err
  }

  return nil
}