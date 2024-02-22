package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func getTesEnv(key string) string {
	if key == "WALLET_DB_NAME" {
		return "test.db"
	}
	if key == "WALLET_HOST" {
		return "0.0.0.0"
	}
	if key == "WALLET_PORT" {
		return "8082"
	}
	return ""
}

func TestCreateWallet(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	t.Cleanup(cancel)
	t.Cleanup(func() {
		os.Remove(getTesEnv("WALLET_DB_NAME"))
	})

	t.Log("Start the server")
	go run(ctx, io.Discard, getTesEnv)
	err := waitForReady(ctx, 3*time.Second, "http://0.0.0.0:8082/health")
	assert.NoError(t, err, "waiting for ready should have no err")

	t.Log("Create a wallet with 100 amount")
	reqBody := bytes.Buffer{}
	reqBody.WriteString(`{"amount": 100}`)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "http://0.0.0.0:8082/wallet/andrei", &reqBody)
	req.Header.Set("Content-Type", "application/json")
	client := http.Client{}
	resp, err := client.Do(req)

	assert.NoError(t, err, "doing request to http://0.0.0.0:8082/wallet/andrei should succeed")
	defer resp.Body.Close()

	t.Log("Check create wallet response")
	data, err := io.ReadAll(resp.Body)
	assert.NoError(t, err, "data should be read from response body")
	assert.JSONEqf(t, `{"balance": 100, "userId": "andrei"}`, string(data), "response should contain wallet response")
}

func waitForReady(
	ctx context.Context,
	timeout time.Duration,
	endpoint string,
) error {
	client := http.Client{}
	startTime := time.Now()
	for {
		req, err := http.NewRequestWithContext(
			ctx,
			http.MethodGet,
			endpoint,
			nil,
		)
		if err != nil {
			return fmt.Errorf("failed to create request: %w", err)
		}

		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("Error making request: %s\n", err.Error())
			continue
		}
		if resp.StatusCode == http.StatusOK {
			resp.Body.Close()
			return nil
		}
		resp.Body.Close()

		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if time.Since(startTime) >= timeout {
				return fmt.Errorf("timeout reached while waiting for endpoint")
			}
			// wait a little while between checks
			time.Sleep(250 * time.Millisecond)
		}
	}
}
