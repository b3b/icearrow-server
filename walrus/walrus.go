package walrus

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

// WalrusClient represents a client to interact with the Walrus HTTP API
type WalrusClient struct {
	baseURL string
}

func NewWalrusClient(baseURL string) *WalrusClient {
	return &WalrusClient{baseURL: baseURL}
}

func (c *WalrusClient) PutData(data []byte) (string, error) {
	req, err := http.NewRequest(http.MethodPut, c.baseURL+"/v1/store", bytes.NewReader(data))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/octet-stream")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	if alreadyCertified, ok := result["alreadyCertified"].(map[string]interface{}); ok {
		if blobID, ok := alreadyCertified["blobId"].(string); ok {
			return blobID, nil
		}
	}
	if newlyCreated, ok := result["newlyCreated"].(map[string]interface{}); ok {
		if blobObject, ok := newlyCreated["blobObject"].(map[string]interface{}); ok {
			if blobID, ok := blobObject["blobId"].(string); ok {
				return blobID, nil
			}
		}
	}

	return "", nil
}

func (c *WalrusClient) GetData(blobID string) ([]byte, error) {
	resp, err := http.Get(c.baseURL + "/v1/" + blobID)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
