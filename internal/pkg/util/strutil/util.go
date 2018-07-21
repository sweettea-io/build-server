package strutil

// InSlice returns a boolean indicating whether the provided
// string exists inside the provided slice of strings.
func InSlice(str string, strList []string) bool {
  for _, entry := range strList {
    if entry == str {
      return true
    }
  }

  return false
}
