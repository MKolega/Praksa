package client

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func FetchData(url string) (io.ReadCloser, error) {
	r, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data from %s: %v", url, err)
	}
	return r.Body, nil
}

func ProcessData(url string, out interface{}) error {
	resp, err := FetchData(url)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			log.Println("failed to close response body:", err)
		}
	}(resp)

	if err := json.NewDecoder(resp).Decode(out); err != nil {
		return fmt.Errorf("failed to decode or process data from %s: %v", url, err)
	}

	log.Printf("Successfully processed data from %s.", url)
	return nil
}
