package elasticsearch

// FoodIndexMapping defines the Elasticsearch mapping for the food index
const FoodIndexMapping = `
{
  "restaurant": {
    "aliases": {

    },
    "mappings": {
      "_meta": {
        "created_by": "file-data-visualizer"
      },
      "properties": {
        "city": {
          "type": "keyword"
        },
        "cuisine": {
          "type": "keyword"
        },
        "lat": {
          "type": "double"
        },
        "location": {
          "type": "geo_point"
        },
        "lon": {
          "type": "double"
        },
        "name": {
          "type": "text",
          "fields": {
            "keyword": {
              "type": "keyword",
              "ignore_above": 256
            }
          }
        },
        "price": {
          "type": "keyword"
        },
        "region": {
          "type": "keyword"
        },
        "star": {
          "type": "long"
        },
        "url": {
          "type": "keyword"
        },
        "year": {
          "type": "long"
        },
        "zipCode": {
          "type": "keyword"
        }
      }
    }
  }
}
`

// const FoodIndexMapping = `
// {
//   "settings": {
//     "number_of_shards": 1,
//     "number_of_replicas": 0,
//     "analysis": {
//       "analyzer": {
//         "food_analyzer": {
//           "type": "custom",
//           "tokenizer": "standard",
//           "filter": ["lowercase", "asciifolding", "word_delimiter", "stop"]
//         },
//         "ngram_analyzer": {
//           "type": "custom",
//           "tokenizer": "ngram_tokenizer",
//           "filter": ["lowercase", "asciifolding"]
//         }
//       },
//       "tokenizer": {
//         "ngram_tokenizer": {
//           "type": "ngram",
//           "min_gram": 2,
//           "max_gram": 4,
//           "token_chars": ["letter", "digit"]
//         }
//       }
//     }
//   },
//   "mappings": {
//     "properties": {
//       "id": {
//         "type": "keyword"
//       },
//       "name": {
//         "type": "text",
//         "analyzer": "food_analyzer",
//         "fields": {
//           "keyword": {
//             "type": "keyword"
//           },
//           "ngram": {
//             "type": "text",
//             "analyzer": "ngram_analyzer"
//           }
//         }
//       },
//       "description": {
//         "type": "text",
//         "analyzer": "food_analyzer"
//       },
//       "price": {
//         "type": "double"
//       },
//       "category_id": {
//         "type": "keyword"
//       },
//       "status": {
//         "type": "keyword"
//       },
//       "images": {
//         "type": "keyword"
//       },
//       "avg_rating": {
//         "type": "double"
//       },
//       "rating_count": {
//         "type": "integer"
//       },
//       "created_at": {
//         "type": "date"
//       },
//       "updated_at": {
//         "type": "date"
//       }
//     }
//   }
// }`
