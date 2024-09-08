//go:build e2e

package walrus

import (
	"bytes"
	"fmt"
	"testing"
	"time"
)

func TestWalrusClient_E2E(t *testing.T) {
	baseURL := "http://127.0.0.1:31415"

	client := NewWalrusClient(baseURL)

	generateTimestampedData := func() []byte {
		data := fmt.Sprintf("test hippos data %s", time.Now().Format(time.RFC3339))
		return []byte(data)
	}

	testPutAndGet := func(t *testing.T, data []byte) {
		blobID, err := client.PutData(data)
		if err != nil {
			t.Fatalf("Failed to put data: %v", err)
		}
		fmt.Printf("Blob ID: %s\n", blobID)

		gotData, err := client.GetData(blobID)
		if err != nil {
			t.Fatalf("Failed to get data: %v", err)
		}

		if !bytes.Equal(gotData, data) {
			t.Errorf("Expected GET response %s, but got %s", data, gotData)
		}
	}

	data1 := generateTimestampedData()
	testPutAndGet(t, data1)

	time.Sleep(1 * time.Second)

	data2 := generateTimestampedData()
	testPutAndGet(t, data2)

	data3 := data1
	testPutAndGet(t, data3)
}
