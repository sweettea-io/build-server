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

  // Calculate main executable's unique name (what it will be updated to).
  curr_main_exec := fmt.Sprintf("main.%s", bp.FileExt)
  new_main_exec := fmt.Sprintf("main_%s.%s", targetUid, bp.FileExt)

  // Rename current main executable to new name.
  if err := os.Rename(
    fmt.Sprintf("%s/%s", bpPath, curr_main_exec),
    fmt.Sprintf("%s/%s", bpPath, new_main_exec),
  ); err != nil {
    return err
  }

  // Read in Dockerfile of buildpack.
  dfilePath := fmt.Sprintf("%s/Dockerfile", bpPath)
  dfileContents, err := ioutil.ReadFile(dfilePath)

  if err != nil {
    return err
  }

  // Get lines of Dockerfile as an array of strings.
  dfileLines := strings.Split(string(dfileContents), "\n")

  // Update last line of Dockerfile, modifying the main executable's name to its new one.
  dfileLines[len(dfileLines) - 1] = strings.Replace(
    dfileLines[len(dfileLines) - 1],
    curr_main_exec,
    new_main_exec,
    1,
  )

  // Save new Dockerfile contents in place.
  if err := ioutil.WriteFile(dfilePath, []byte(strings.Join(dfileLines, "\n")), 0644); err != nil {
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