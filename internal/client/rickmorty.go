package client

import (
	"Users/tendusmac/Desktop/NEU/Akamai/rickmorty-api/internal/model"
	"Users/tendusmac/Desktop/NEU/Akamai/rickmorty-api/internal/util"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type Response struct {
	Info struct {
		Next string `json:"next"`
	}
	Results []model.Character `json:"results"`
}

func GetCharsWithFilters() ([]model.Character, error) {
	const maxRetries = 5
	var allChars []model.Character
	baseURL := "https://rickandmortyapi.com/api/character/?species=Human&status=Alive"
	// baseURL := "https://invalid.rickandmortyapi.com/api/character/?species=Human&status=Alive"
	client := &http.Client{Timeout: 10 * time.Second}

	for url := baseURL; url != ""; {
		var resp *http.Response
		var err error

		// Retry wrapper
		for attempt := 0; attempt < maxRetries; attempt++ {
			resp, err = client.Get(url)

			if err == nil && resp.StatusCode < 400 {
				// Success
				break
			}

			if err != nil || util.ShouldRetry(resp.StatusCode) {
				wait := util.Backoff(attempt)
				if err != nil {
					log.Printf("Retrying request due to error (attempt %d/%d): %v. Backing off for %v\n", attempt+1, maxRetries, err, wait)
				} else {
					log.Printf("Retrying request due to status code %d (attempt %d/%d). Backing off for %v\n", resp.StatusCode, attempt+1, maxRetries, wait)
				}
				time.Sleep(wait)
				continue
			}
			break
		}

		if err != nil {
			return nil, fmt.Errorf("failed to fetch data: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
		}

		body, _ := io.ReadAll(resp.Body)
		var apiResp Response
		if err := json.Unmarshal(body, &apiResp); err != nil {
			return nil, err
		}

		for _, c := range apiResp.Results {
			if c.Origin.Name != "" && containsEarth(c.Origin.Name) {
				allChars = append(allChars, c)
			}
		}

		url = apiResp.Info.Next
		// time.Sleep(200 * time.Millisecond) // space out pagination
	}

	return allChars, nil
}

func containsEarth(origin string) bool {
	return len(origin) >= 5 && origin[:5] == "Earth"
}
