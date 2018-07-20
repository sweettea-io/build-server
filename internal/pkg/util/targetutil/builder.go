package targetutil

import (
  "errors"
  "fmt"
  "io/ioutil"
  "os"
  "strings"
  "github.com/sweettea-io/build-server/internal/pkg/util/strutil"
)

func AttachBuildpack(buildpack string, bpPath string, targetPath string, targetUid string) error {
  // Get file extension for buildpack's language.
  mainFileExt, err := extForBuildpack(buildpack)

  if err != nil {
    return err
  }

  // Trim any trailing slashes from path args.
  bpPath = strings.TrimRight(bpPath, "/")
  targetPath = strings.TrimRight(targetPath, "/")

  // Rename buildpack's main file to include the targetUid.
  if err := os.Rename(
    fmt.Sprintf("%s/main.%s", bpPath, mainFileExt),
    fmt.Sprintf("%s/main_%s.%s", bpPath, targetUid, mainFileExt),
  ); err != nil {
    return err
  }

  // Rename buildpack's src directory to be `targetUid`.
  if err := os.Rename(
    fmt.Sprintf("%s/src", bpPath),
    fmt.Sprintf("%s/%s", bpPath, targetUid),
  ); err != nil {
    return err
  }

  // Get all contents of the buildpack dir.
  bpEntries, err := ioutil.ReadDir(bpPath)

  if err != nil {
    return err
  }

  // Files & dirs we don't want to copy over from the buildpack.
  ignorables := []string{
    ".gitignore",
    ".git",
    "README.md",
  }

  // Move all buildpack files & dirs into the target project, except for those in `ignorables`.
  for _, file := range bpEntries {
    entryName := file.Name() // can be a file or directory

    // Ignore file/dir if in ignorables list.
    if strutil.InSlice(entryName, ignorables) {
      continue
    }

    // Move buildpack file/dir into target.
    if err := os.Rename(
      fmt.Sprintf("%s/%s", bpPath, entryName),
      fmt.Sprintf("%s/%s", targetPath, entryName),
    ); err != nil {
      return err
    }
  }

  return nil
}

func extForBuildpack(buildpack string) (string, error) {
  if strings.HasPrefix(buildpack, "python") {
    return "py", nil
  } else {
    return "", errors.New(fmt.Sprintf("language not recognized for buildpack: %s", buildpack))
  }
}