package pinecone

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Pinecone struct {
	ApiKey      string
	Environment string
}

const (
	baseUrl = "https://controller.%s.pinecone.io"
)

// list_collections
// GET
// https://controller.{environment}.pinecone.io/collections
// This operation returns a list of your Pinecone collections.
//
// 200 Array of String - This operation returns a list of all the collections in your current project.

type CreateCollectionBodyParams struct {
	// The name of the collection to be created.
	Name string `json:"name"`
	// The name of the source index to be used as the source for the collection.
	Source string `json:"source"`
}

// create_collection
// POST
// https://controller.{environment}.pinecone.io/collections
// This operation creates a Pinecone collection. Not supported by projects on the gcp-starter environment.
//
// 201 String - The collection has been successfully created.
// 400 String - Bad request. Request exceeds quota or collection name is invalid.
// 409 String - A collection with the name provided already exists.
// 500 String - Internal error. Can be caused by invalid parameters.
func (p *Pinecone) CreateCollection(bodyParams CreateCollectionBodyParams) (*string, error) {
	url := fmt.Sprintf(baseUrl+"/collections", p.Environment)

	// convert struct to byte[]
	payloadBytes, err := json.Marshal(bodyParams)
	if err != nil {
		return nil, err
	}

	// convert byte[] to io.Reader
	payload := bytes.NewReader(payloadBytes)

	// initialize a request
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return nil, err
	}

	req.Header.Add("accept", "text/plain")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("Api-Key", p.ApiKey)

	// fire off the request
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	fmt.Println(string(body))
	bodyString := string(body)

	switch {
	case res.StatusCode < 300: // 2xx
		return &bodyString, nil
	default: // non-2xx
		return nil, fmt.Errorf("CreateCollection failed with status code %d and message %q", res.StatusCode, bodyString)
	}
}

type DescribeCollectionResponse struct {
	Name string `json:"name"`
	// values: Ready,
	Status string `json:"status"`
	Size   int64  `json:"size"`
	// The dimension of the vectors stored in the collection.
	Dimension   int64 `json:"dimension"`
	VecotrCount int64 `json:"vector_count"`
}

// describe_collection
// GET
// https://controller.{environment}.pinecone.io/collections/{collectionName}
// Get a description of a collection.
//
// 200 JSON - Configuration information and deployment status of the index
// 404 String - Index not found.
// 500 String - Internal error. Can be caused by invalid parameters.
func (p *Pinecone) DescribeCollection(name string) (*DescribeCollectionResponse, error) {
	url := fmt.Sprintf(baseUrl+"/collections/%s", p.Environment, name)

	// initialize a request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("Api-Key", p.ApiKey)

	// fire off the request
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	bodyString := string(body)

	switch {
	case res.StatusCode < 300: // 2xx
		// unmarshal json to struct
		descCollectionResponse := &DescribeCollectionResponse{}
		err := json.Unmarshal(body, descCollectionResponse)
		if err != nil {
			return nil, err
		}
		return descCollectionResponse, nil
	default: // non-2xx
		return nil, fmt.Errorf("DescribeCollection failed with status code %d and message %q", res.StatusCode, bodyString)
	}
}

// delete_collection
// DELETE
// https://controller.{environment}.pinecone.io/collections/{collectionName}
// This operation deletes an existing collection.
//
// 202 String - The index has been successfully deleted.
// 404 String - Collection not found.
// 500 String - Internal error. Can be caused by invalid parameters.
func (p *Pinecone) DeleteCollection(name string) (*string, error) {
	url := fmt.Sprintf(baseUrl+"/collections/%s", p.Environment, name)

	// initialize a request
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("Api-Key", p.ApiKey)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	bodyString := string(body)

	switch {
	case res.StatusCode < 300: // 2xx
		return &bodyString, nil
	default: // non-2xx
		return nil, fmt.Errorf("DeleteCollection failed with status code %d and message %q", res.StatusCode, bodyString)
	}
}

type CreateIndexBodyParams struct {
	// The name of the index to be created. The maximum length is 45 characters.
	Name string `json:"name"`
	// The dimensions of the vectors to be inserted in the index
	Dimension int64 `json:"dimension"`
	// The distance metric to be used for similarity search. You can use 'euclidean', 'cosine', or 'dotproduct'.
	Metric string `json:"metric"`
	// The number of pods for the index to use,including replicas.
	Pods int64 `json:"pods"`
	// The number of replicas. Replicas duplicate your index. They provide higher availability and throughput.
	Replicas int64 `json:"replicas"`
	// The type of pod to use. One of s1, p1, or p2 appended with . and one of x1, x2, x4, or x8.
	PodType string `json:"pod_type"`
	// Configuration for the behavior of Pinecone's internal metadata index. By default, all metadata is indexed; when metadata_config is present, only specified metadata fields are indexed. To specify metadata fields to index, provide a JSON object of the following form:
	MetadataConfig *map[string]interface{} `json:"metadata_config"`
	// The name of the collection to create an index from
	SourceCollection string `json:"source_collection"`
}

// create_index
// POST
// https://controller.{environment}.pinecone.io/databases
// This operation creates a Pinecone index. You can use it to specify the measure of similarity, the dimension of vectors to be stored in the index, the numbers of replicas to use, and more.
func (p *Pinecone) CreateIndex(data CreateIndexBodyParams) (*string, error) {
	url := fmt.Sprintf(baseUrl+"/databases", p.Environment)

	// set default values
	if data.Metric == "" {
		data.Metric = "cosine"
	}
	if data.Pods == 0 {
		data.Pods = 1
	}
	if data.Replicas == 0 {
		data.Replicas = 1
	}
	if data.PodType == "" {
		data.PodType = "p1.x1"
	}

	// convert struct to byte[]
	payloadBytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	// convert byte[] to io.Reader
	payload := bytes.NewReader(payloadBytes)

	// initialize a request
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return nil, err
	}

	req.Header.Add("accept", "text/plain")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("Api-Key", p.ApiKey)

	// fire off the request
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	fmt.Println(string(body))
	bodyString := string(body)

	switch {
	case res.StatusCode < 300: // 2xx
		return &bodyString, nil
	default: // non-2xx
		return nil, fmt.Errorf("CreateIndex failed with status code %d and message %q", res.StatusCode, bodyString)
	}
}

type DescribeIndexResponse struct {
	Database struct {
		Name      string `json:"name"`
		Metric    string `json:"metric"`
		Dimension int64  `json:"dimension"`
		Replicas  int64  `json:"replicas"`
		Shards    int64  `json:"shards"`
		Pods      int64  `json:"pods"`
	} `json:"database"`
	Status struct {
		Waiting []interface{} `json:"waiting"`
		Crashed []interface{} `json:"crashed"`
		Host    string        `json:"host"`
		Port    int64         `json:"port"`
		// values: Initializing, Ready,
		State string `json:"state"`
		Ready bool   `json:"ready"`
	} `json:"status"`
}

// describe_index
// GET
// https://controller.{environment}.pinecone.io/databases/{indexName}
// Get a description of an index.
func (p *Pinecone) DescribeIndex(name string) (*DescribeIndexResponse, error) {
	if name == "" {
		return nil, fmt.Errorf("DescribeIndex failed: name argument was not specified")
	}
	url := fmt.Sprintf(baseUrl+"/databases/%s", p.Environment, name)

	// initialize a request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("Api-Key", p.ApiKey)

	// fire off the request
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	bodyString := string(body)

	switch {
	case res.StatusCode < 300: // 2xx
		// unmarshal json to struct
		descIndexResponse := &DescribeIndexResponse{}
		err := json.Unmarshal(body, descIndexResponse)
		if err != nil {
			return nil, err
		}
		return descIndexResponse, nil
	default: // non-2xx
		return nil, fmt.Errorf("DescribeIndex failed with status code %d and message %q", res.StatusCode, bodyString)
	}
}

type ConfigureIndexRequest struct {
	// The new pod type for the index. One of s1, p1, or p2 appended with . and one of x1, x2, x4, or x8.
	PodType string `json:"pod_type"`
	// The desired number of replicas for the index.
	Replicas int64 `json:"replicas"`
}

// configure_index
// PATCH
// https://controller.{environment}.pinecone.io/databases/{indexName}
// This operation specifies the pod type and number of replicas for an index. Not supported by projects on the gcp-starter environment.
// (Not supported for free tier)
//
// 202 String - The index has been successfully updated
// 400 String - Bad request,not enough quota.
// 404 String - Index not found.
// 500 String - Internal error. Can be caused by invalid parameters.
func (p *Pinecone) ConfigureIndex(name string, data *ConfigureIndexRequest) (*string, error) {
	url := fmt.Sprintf(baseUrl+"/databases/%s", p.Environment, name)

	payloadBytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	// convert byte[] to io.Reader
	payload := bytes.NewReader(payloadBytes)

	// initialize a request
	req, err := http.NewRequest("PATCH", url, payload)
	if err != nil {
		return nil, err
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("Api-Key", p.ApiKey)

	// fire off the request
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	bodyString := string(body)

	switch {
	case res.StatusCode < 300: // 2xx
		return &bodyString, nil
	default: // non-2xx
		return nil, fmt.Errorf("ConfigureIndex failed with status code %d and message %q", res.StatusCode, bodyString)
	}
}

// delete_index
// DELETE
// https://controller.{environment}.pinecone.io/databases/{indexName}
// This operation deletes an existing index.
//
// 202 String - The index has been successfully deleted
// 404 String - Index not found.
// 500 String - Internal error. Can be caused by invalid parameters.
func (p *Pinecone) DeleteIndex(name string) (*string, error) {
	url := fmt.Sprintf(baseUrl+"/databases/%s", p.Environment, name)

	// initialize a request
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("Api-Key", p.ApiKey)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	bodyString := string(body)

	switch {
	case res.StatusCode < 300: // 2xx
		return &bodyString, nil
	default: // non-2xx
		return nil, fmt.Errorf("DeleteIndex failed with status code %d and message %q", res.StatusCode, bodyString)
	}
}
