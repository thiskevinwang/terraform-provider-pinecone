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
func (p *Pinecone) CreateIndex() {
	url := fmt.Sprintf("https://controller.%s.pinecone.io/databases", p.Environment)

	// payload := strings.NewReader("{\"metric\":\"cosine\",\"pods\":1,\"replicas\":1,\"pod_type\":\"p1.x1\",\"name\":\"test\",\"dimension\":1536}")
	data := map[string]interface{}{
		"metric":    "cosine",
		"pods":      1,
		"replicas":  1,
		"pod_type":  "p1.x1",
		"name":      "test",
		"dimension": 1536,
	}
	payloadBytes, _ := json.Marshal(data)

	// convert byte[] to io.Reader
	payload := bytes.NewReader(payloadBytes)

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("accept", "text/plain")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("Api-Key", p.ApiKey)

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	fmt.Println(string(body))
}

// Read
func (p *Pinecone) DescribeIndex(name string) {

	url := fmt.Sprintf("https://controller.%s.pinecone.io/databases/%s", p.Environment, name)

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("accept", "application/json")
	req.Header.Add("Api-Key", p.ApiKey)

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	fmt.Println(string(body))

}
