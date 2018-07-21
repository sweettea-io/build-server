package fileutil

import "os"

// Exists returns whether the given file or directory exists or not
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

// RemoveIfExists removes a file or directory and all of its subcontents if it exists.
func RemoveIfExists(path string) error {
  exists, err := Exists(path)

  if err != nil {
    return err
  }

  if !exists {
    return nil
  }

  if err := os.RemoveAll(path); err != nil {
    return err
  }

  return nil
}