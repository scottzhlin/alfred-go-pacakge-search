package alfred_go_pacakge_search

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"time"
)

// search demo:
// GET https://api.godoc.org/search?q=elasticsearch
// response:
// {
//    "results": [
//        {
//            "name": "elasticsearch",
//            "path": "github.com/elastic/beats/libbeat/outputs/elasticsearch",
//            "import_count": 468,
//            "stars": 9233,
//            "score": 0.92178405
//        },
//        {
//            "name": "elastic",
//            "path": "github.com/olivere/elastic",
//            "import_count": 286,
//            "synopsis": "Package elastic provides an interface to the Elasticsearch server (https://www.elastic.co/products/elasticsearch).",
//            "stars": 5513,
//            "score": 1
//        },
//        {
//            "name": "elasticsearchservice",
//            "path": "github.com/aws/aws-sdk-go/service/elasticsearchservice",
//            "import_count": 235,
//            "synopsis": "Package elasticsearchservice provides the client and types for making API requests to Amazon Elasticsearch Service.",
//            "stars": 6181,
//            "score": 0.9801
//        }
//    ]
// }

const goPackageSearchApiFormat = "https://api.godoc.org/search?q=%s"

const defaultPageLimit = 10

// GoPackageSearchResult define api.godoc.org/search result set
type GoPackageSearchResult struct {
	Results GoPackages `json:"results"`
}

// GoPackage define api.godoc.org/search result item
type GoPackage struct {
	Name        string  `json:"name"`
	Path        string  `json:"path"`
	Synopsis    string  `json:"synopsis"`
	ImportCount int     `json:"import_count"`
	Stars       int     `json:"stars"`
	Score       float64 `json:"score"`
}

type GoPackages []*GoPackage

func (g GoPackages) Len() int {
	return len(g)
}

func (g GoPackages) Less(i, j int) bool {
	return g[i].Score > g[j].Score
}

func (g GoPackages) Swap(i, j int) {
	g[i], g[j] = g[j], g[i]
}

func SearchGoPackages(q string, limit int) ([]*GoPackage, error) {
	if len(q) == 0 {
		return nil, nil
	}

	if limit <= 0 {
		limit = defaultPageLimit
	}

	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf(goPackageSearchApiFormat, q), nil)
	if err != nil {
		return nil, err
	}

	// set request timeout
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Second)
	defer cancel()
	request = request.WithContext(ctx)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("failed to send a http GET request: %w", err)
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	results := new(GoPackageSearchResult)
	if err := json.Unmarshal(body, results); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body to json: %w", err)
	}

	// sort result item by score
	sort.Sort(results.Results)
	if len(results.Results) < limit {
		limit = len(results.Results)
	}

	return results.Results[:limit], nil
}
