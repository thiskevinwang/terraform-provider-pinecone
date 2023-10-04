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

// Create
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

// Read
func (p *Pinecone) DescribeIndex(name string) (*string, error) {
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

	fmt.Println(string(body))
	bodyString := string(body)

	switch {
	case res.StatusCode < 300: // 2xx
		return &bodyString, nil
	default: // non-2xx
		return nil, fmt.Errorf("DescribeIndex failed with status code %d and message %q", res.StatusCode, bodyString)
	}
}
