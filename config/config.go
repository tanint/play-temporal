package config

import (
	"os"

	"go.temporal.io/sdk/client"
)

// GetTemporalClientOptions returns the client options for connecting to the Temporal server
func GetTemporalClientOptions() client.Options {
	// Default to localhost if TEMPORAL_HOST is not set
	temporalHost := os.Getenv("TEMPORAL_HOST")
	if temporalHost == "" {
		temporalHost = "localhost:7233"
	}

	// Default to default namespace if TEMPORAL_NAMESPACE is not set
	namespace := os.Getenv("TEMPORAL_NAMESPACE")
	if namespace == "" {
		namespace = "default"
	}

	return client.Options{
		HostPort:  temporalHost,
		Namespace: namespace,
	}
}
