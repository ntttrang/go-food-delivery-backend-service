package datatype

// ElasticSearchConfig represents the configuration for Elasticsearch
type ElasticSearchConfig struct {
	Addresses []string // List of Elasticsearch node addresses
	Username  string   // Username for authentication
	Password  string   // Password for authentication
	CloudID   string   // Cloud ID for Elastic Cloud
	APIKey    string   // API key for authentication
	IndexName string   // Default index name
}
