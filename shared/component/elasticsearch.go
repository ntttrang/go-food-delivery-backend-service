package shareComponent

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/ntttrang/go-food-delivery-backend-service/shared/datatype"
	"github.com/pkg/errors"
)

type ElasticsearchClient struct {
	client    *elasticsearch.Client
	indexName string
}

func NewElasticsearchClient(cfg datatype.ElasticSearchConfig) (*ElasticsearchClient, error) {
	// Configure the Elasticsearch client
	esCfg := elasticsearch.Config{
		Addresses: cfg.Addresses,
		Username:  cfg.Username,
		Password:  cfg.Password,
		CloudID:   cfg.CloudID,
		APIKey:    cfg.APIKey,
	}

	// Create the client
	client, err := elasticsearch.NewClient(esCfg)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// Set default index name if not provided
	indexName := cfg.IndexName
	if indexName == "" {
		indexName = "foods"
	}

	return &ElasticsearchClient{
		client:    client,
		indexName: indexName,
	}, nil
}

// IndexDocument indexes a document in Elasticsearch
func (es *ElasticsearchClient) IndexDocument(ctx context.Context, id string, document interface{}) error {
	// Convert document to JSON
	documentJSON, err := json.Marshal(document)
	if err != nil {
		return errors.WithStack(err)
	}

	// Create the request
	req := esapi.IndexRequest{
		Index:      es.indexName,
		DocumentID: id,
		Body:       bytes.NewReader(documentJSON),
		Refresh:    "true", // Make the document immediately available for search
	}

	// Execute the request
	res, err := req.Do(ctx, es.client)
	if err != nil {
		return errors.WithStack(err)
	}
	defer res.Body.Close()

	// Check response status
	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return errors.Errorf("error parsing the response body: %s", err)
		}
		return errors.Errorf("[%s] %s: %s", res.Status(),
			e["error"].(map[string]interface{})["type"],
			e["error"].(map[string]interface{})["reason"])
	}

	return nil
}

// DeleteDocument deletes a document from Elasticsearch
func (es *ElasticsearchClient) DeleteDocument(ctx context.Context, id string) error {
	// Create the request
	req := esapi.DeleteRequest{
		Index:      es.indexName,
		DocumentID: id,
		Refresh:    "true", // Make the deletion immediately visible
	}

	// Execute the request
	res, err := req.Do(ctx, es.client)
	if err != nil {
		return errors.WithStack(err)
	}
	defer res.Body.Close()

	// Check response status (404 is not an error, it means the document was not found)
	if res.StatusCode != 404 && res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return errors.Errorf("error parsing the response body: %s", err)
		}
		return errors.Errorf("[%s] %s: %s", res.Status(),
			e["error"].(map[string]interface{})["type"],
			e["error"].(map[string]interface{})["reason"])
	}

	return nil
}

// Search performs a search in Elasticsearch
func (es *ElasticsearchClient) Search(ctx context.Context, query map[string]interface{}, from, size int) ([]map[string]interface{}, int64, map[string]interface{}, error) {
	// Convert query to JSON
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, 0, nil, errors.WithStack(err)
	}

	// Perform the search request
	res, err := es.client.Search(
		es.client.Search.WithContext(ctx),
		es.client.Search.WithIndex(es.indexName),
		es.client.Search.WithBody(&buf),
		es.client.Search.WithTrackTotalHits(true),
		es.client.Search.WithFrom(from),
		es.client.Search.WithSize(size),
	)
	if err != nil {
		return nil, 0, nil, errors.WithStack(err)
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return nil, 0, nil, errors.Errorf("error parsing the response body: %s", err)
		}
		// Print the response status and error information
		return nil, 0, nil, errors.Errorf("[%s] %s: %s", res.Status(),
			e["error"].(map[string]interface{})["type"],
			e["error"].(map[string]interface{})["reason"])
	}

	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, 0, nil, errors.WithStack(err)
	}

	// Extract total hits
	total := int64(r["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64))

	// Extract hits
	hits := r["hits"].(map[string]interface{})["hits"].([]interface{})
	results := make([]map[string]interface{}, len(hits))

	for i, hit := range hits {
		source := hit.(map[string]interface{})["_source"].(map[string]interface{})
		id := hit.(map[string]interface{})["_id"].(string)
		source["id"] = id
		results[i] = source
	}

	// Extract aggregations if present
	var aggregations map[string]interface{}
	if aggs, ok := r["aggregations"].(map[string]interface{}); ok {
		aggregations = aggs
	}

	return results, total, aggregations, nil
}

// CreateIndex creates an index with the specified mapping
func (es *ElasticsearchClient) CreateIndex(ctx context.Context, mapping string) error {
	res, err := es.client.Indices.Exists([]string{es.indexName})
	if err != nil {
		return errors.WithStack(err)
	}

	// If index already exists, don't recreate it
	if res.StatusCode == 200 {
		log.Printf("Index %s already exists", es.indexName)
		return nil
	}

	// Create index with mapping
	res, err = es.client.Indices.Create(
		es.indexName,
		es.client.Indices.Create.WithBody(strings.NewReader(mapping)),
	)
	if err != nil {
		return errors.WithStack(err)
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return errors.Errorf("error parsing the response body: %s", err)
		}
		// Print the response status and error information
		return errors.Errorf("[%s] %s: %s", res.Status(),
			e["error"].(map[string]interface{})["type"],
			e["error"].(map[string]interface{})["reason"])
	}

	log.Printf("Index %s created successfully", es.indexName)
	return nil
}

// BulkIndex indexes multiple documents in a single request
func (es *ElasticsearchClient) BulkIndex(ctx context.Context, documents map[string]interface{}) error {
	if len(documents) == 0 {
		return nil
	}

	var buf bytes.Buffer

	for id, doc := range documents {
		// Action line
		meta := map[string]interface{}{
			"index": map[string]interface{}{
				"_index": es.indexName,
				"_id":    id,
			},
		}
		if err := json.NewEncoder(&buf).Encode(meta); err != nil {
			return errors.WithStack(err)
		}

		// Document line
		if err := json.NewEncoder(&buf).Encode(doc); err != nil {
			return errors.WithStack(err)
		}
	}

	// Perform the bulk request
	res, err := es.client.Bulk(
		bytes.NewReader(buf.Bytes()),
		es.client.Bulk.WithIndex(es.indexName),
		es.client.Bulk.WithRefresh("true"),
	)
	if err != nil {
		return errors.WithStack(err)
	}
	defer res.Body.Close()

	if res.IsError() {
		var raw map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&raw); err != nil {
			return errors.Errorf("failure to parse response body: %s", err)
		}
		return errors.Errorf("bulk request failed: %s", res.Status())
	}

	var blk map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&blk); err != nil {
		return errors.Errorf("failure to parse response body: %s", err)
	}

	// Check if there were any errors in the bulk operation
	if blk["errors"].(bool) {
		// Extract the first error
		for _, item := range blk["items"].([]interface{}) {
			for _, v := range item.(map[string]interface{}) {
				if v.(map[string]interface{})["error"] != nil {
					e := v.(map[string]interface{})["error"].(map[string]interface{})
					return errors.Errorf("[%s] %s: %s", v.(map[string]interface{})["status"],
						e["type"], e["reason"])
				}
			}
		}
		return errors.New("bulk request contained errors")
	}

	return nil
}

// GetDocument retrieves a document by ID
func (es *ElasticsearchClient) GetDocument(ctx context.Context, id string) (map[string]interface{}, error) {
	// Create the request
	req := esapi.GetRequest{
		Index:      es.indexName,
		DocumentID: id,
	}

	// Execute the request
	res, err := req.Do(ctx, es.client)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer res.Body.Close()

	// Check if document was found
	if res.StatusCode == 404 {
		return nil, fmt.Errorf("document with ID %s not found", id)
	}

	// Check for other errors
	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return nil, errors.Errorf("error parsing the response body: %s", err)
		}
		return nil, errors.Errorf("[%s] %s: %s", res.Status(),
			e["error"].(map[string]interface{})["type"],
			e["error"].(map[string]interface{})["reason"])
	}

	// Parse the response
	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, errors.WithStack(err)
	}

	// Extract the source
	source := r["_source"].(map[string]interface{})
	source["id"] = id

	return source, nil
}

// WithIndex returns a new ElasticsearchClient with the specified index name
func (es *ElasticsearchClient) WithIndex(indexName string) *ElasticsearchClient {
	return &ElasticsearchClient{
		client:    es.client,
		indexName: indexName,
	}
}
