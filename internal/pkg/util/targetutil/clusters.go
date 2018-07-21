package targetutil

// Supported target clusters
const (
  // Cluster used for model training.
  Train = "train"

  // Cluster used for serving model predictions from an API.
  API = "api"
)

// validTargetClusters is a map of all supported target clusters. It exists
// simply to provide faster lookup times during calls to `IsValidTargetCluster`.
var validTargetClusters = map[string]bool {
  Train: true,
  API: true,
}

// IsValidTargetCluster returns a boolean indicating
// whether the provided target cluster is supported.
func IsValidTargetCluster(tc string) bool {
  return validTargetClusters[tc] == true
}