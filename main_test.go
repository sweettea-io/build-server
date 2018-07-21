package main

import (
  "testing"
  "os"
  "github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
  code := m.Run()
  os.Exit(code)
}

func TestFirstCase(t *testing.T) {
  assert.Equal(t, true, true, "first test")
}