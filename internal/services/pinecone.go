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

// create_index
// POST
// https://controller.us-west4-gcp-free.pinecone.io/databases
// This operation creates a Pinecone index. You can use it to specify the measure of similarity, the dimension of vectors to be stored in the index, the numbers of replicas to use, and more.
func (p *Pinecone) CreateIndex(name string, dimension int64, metric string) (*string, error) {
	url := fmt.Sprintf("https://controller.%s.pinecone.io/databases", p.Environment)

	// payload := strings.NewReader("{\"metric\":\"cosine\",\"pods\":1,\"replicas\":1,\"pod_type\":\"p1.x1\",\"name\":\"test\",\"dimension\":1536}")
	data := map[string]interface{}{
		"metric":    metric,
		"name":      name,
		"dimension": dimension,
		"pods":      1,
		"replicas":  1,
		"pod_type":  "p1.x1",
	}
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
// https://controller.us-west4-gcp-free.pinecone.io/databases/{indexName}
// Get a description of an index.
func (p *Pinecone) DescribeIndex(name string) (*DescribeIndexResponse, error) {
	url := fmt.Sprintf("https://controller.%s.pinecone.io/databases/%s", p.Environment, name)

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
// https://controller.us-west4-gcp-free.pinecone.io/databases/{indexName}
// This operation specifies the pod type and number of replicas for an index. Not supported by projects on the gcp-starter environment.
// (Not supported for free tier)
//
// 202 String - The index has been successfully updated
// 400 String - Bad request,not enough quota.
// 404 String - Index not found.
// 500 String - Internal error. Can be caused by invalid parameters.
func (p *Pinecone) ConfigureIndex(name string, data *ConfigureIndexRequest) (*string, error) {
	url := fmt.Sprintf("https://controller.%s.pinecone.io/databases/%s", p.Environment, name)

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
// https://controller.us-west4-gcp-free.pinecone.io/databases/{indexName}
// This operation deletes an existing index.
//
// 202 String - The index has been successfully deleted
// 404 String - Index not found.
// 500 String - Internal error. Can be caused by invalid parameters.
func (p *Pinecone) DeleteIndex(name string) (*string, error) {
	url := fmt.Sprintf("https://controller.%s.pinecone.io/databases/%s", p.Environment, name)

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
