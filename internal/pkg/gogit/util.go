package gogit

import (
  "fmt"
  "io"
  "gopkg.in/src-d/go-git.v4"
  "gopkg.in/src-d/go-git.v4/plumbing"
  URL "net/url"
)

// CloneAtSha git-clones a remote repository and then checks out the provided sha.
func CloneAtSha(url string, sha string, accessToken string, dest string, prog io.Writer) error {
  var repoUrl string
  var err error

  // Create authed url from accessToken if provided.
  if accessToken != "" {
    repoUrl, err = createAuthedUrl(url, accessToken)
  }

  // Create clone options.
  cloneOpts := &git.CloneOptions{URL: repoUrl, Progress: prog}

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
  if err := workingTree.Checkout(checkoutOpts); err != nil {
    return err
  }

  return nil
}

// createAuthedUrl adds `<token>:x-oauth-basic` to an existing url, given an access token.
func createAuthedUrl(url string, accessToken string) (string, error) {
  // Parse provided url.
  u, err := URL.Parse(url)

  if err != nil {
    return "", err
  }

  // Format authed url (ex. "https://myaccesstoken:x-oauth-basic@github.com/my/repo")
  authedUrl := fmt.Sprintf(
    "%s://%s:x-oauth-basic@%s%s",
    u.Scheme,
    accessToken,
    u.Host,
    u.RequestURI(),
  )

  return authedUrl, nil
}