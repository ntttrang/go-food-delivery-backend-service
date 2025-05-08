package elasticsearch

// RestaurantIndexMapping defines the Elasticsearch mapping for restaurants
const RestaurantIndexMapping = `
{
  "settings": {
    "number_of_shards": 1,
    "number_of_replicas": 0,
    "analysis": {
      "analyzer": {
        "autocomplete": {
          "tokenizer": "autocomplete",
          "filter": ["lowercase"]
        },
        "autocomplete_search": {
          "tokenizer": "lowercase"
        }
      },
      "tokenizer": {
        "autocomplete": {
          "type": "edge_ngram",
          "min_gram": 2,
          "max_gram": 10,
          "token_chars": ["letter", "digit"]
        }
      }
    }
  },
  "mappings": {
    "properties": {
      "id": {
        "type": "keyword"
      },
      "name": {
        "type": "text",
        "analyzer": "autocomplete",
        "search_analyzer": "autocomplete_search",
        "fields": {
          "keyword": {
            "type": "keyword",
            "ignore_above": 256
          }
        }
      },
      "address": {
        "type": "text",
        "analyzer": "autocomplete",
        "search_analyzer": "autocomplete_search"
      },
      "city_id": {
        "type": "integer"
      },
      "lat": {
        "type": "double"
      },
      "lng": {
        "type": "double"
      },
      "location": {
        "type": "geo_point"
      },
      "logo": {
        "type": "object",
        "enabled": false
      },
      "cover": {
        "type": "object",
        "enabled": false
      },
      "shipping_fee_per_km": {
        "type": "double"
      },
      "status": {
        "type": "keyword"
      },
      "avg_rating": {
        "type": "double"
      },
      "rating_count": {
        "type": "integer"
      },
      "cuisines": {
        "type": "keyword"
      },
      "popularity_score": {
        "type": "double"
      },
      "delivery_time": {
        "type": "integer"
      },
      "created_at": {
        "type": "date"
      },
      "updated_at": {
        "type": "date"
      }
    }
  }
}
`
