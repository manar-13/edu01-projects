package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func fetchJSON[T any](url string, target *T) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to GET from %s: %v", url, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read body: %v", err)
	}

	return json.Unmarshal(body, target)
}

func fetchJSONCtx[T any](ctx context.Context, url string, target *T) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to GET from %s: %v", url, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read body: %v", err)
	}

	return json.Unmarshal(body, target)
}
