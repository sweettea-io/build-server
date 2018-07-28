package targetutil

import (
  "fmt"
  "io/ioutil"
  "os"
  "strings"
  "github.com/sweettea-io/build-server/internal/pkg/util/strutil"
  "github.com/sweettea-io/build-server/internal/pkg/util/buildpack"
)

// AttachBuildpack moves files/dirs from the provided buildpack path to the provided build target path.
func AttachBuildpack(bp *buildpack.Buildpack, bpPath string, targetPath string, targetUid string) error {
  // Trim any trailing slashes from path args.
  bpPath = strings.TrimRight(bpPath, "/")
  targetPath = strings.TrimRight(targetPath, "/")

  // Rename buildpack's main file to include the targetUid.
  if err := os.Rename(
    fmt.Sprintf("%s/main.%s", bpPath, bp.FileExt),
    fmt.Sprintf("%s/main_%s.%s", bpPath, targetUid, bp.FileExt),
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