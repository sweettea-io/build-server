package targetutil

// Supported target clusters
const Train = "train"
const API = "api"

var validTargetClusters = map[string]bool {
  Train: true,
  API: true,
}

// IsValidTargetCluster returns whether the provided target cluster is supported.
func IsValidTargetCluster(tc string) bool {
  return validTargetClusters[tc] == true
}