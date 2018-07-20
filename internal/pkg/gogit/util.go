package gogit

import (
  "io"
  "gopkg.in/src-d/go-git.v4"
  "gopkg.in/src-d/go-git.v4/plumbing"
)

func CloneAtSha(url string, sha string, dest string, prog io.Writer) error {
  // Create clone options.
  cloneOpts := &git.CloneOptions{URL: url, Progress: prog}

  // Clone repo to specified location.
  repo, err := git.PlainClone(dest, false, cloneOpts)

  if err != nil {
    return err
  }

  // Get repo working tree.
  workingTree, err := repo.Worktree()

  if err != nil {
    return err
  }

  // Create checkout options.
  checkoutOpts := &git.CheckoutOptions{Hash: plumbing.NewHash(sha)}

  // Checkout to the specified sha.
  if err = workingTree.Checkout(checkoutOpts); err != nil {
    return err
  }

  return nil
}