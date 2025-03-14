package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/strata-io/service-extension/orchestrator"
)

type User struct {
	ID    int    `json:"id"`
	Phone string `json:"phone"`
}

func ServeSE(api orchestrator.Orchestrator) error {
	// Initialize the HTTP client that can be reused. Reuse is important
	// as the client will pool the underlying TCP connections to enable
	// connection reuse. Reusing a client also helps ensures the system
	// is not overloaded by the opening of too many connections.
	apiHttp := api.HTTP()
	err := apiHttp.SetClient(
		"myClient",
		&http.Client{
			Timeout: time.Second * 30,
		},
	)
	if err != nil {
		return err
	}

	// This client can also be accessed from other service extensions
	// and used to make HTTP requests concurrently.
	client, err := apiHttp.GetClient("myClient")
	if err != nil {
		return fmt.Errorf("unable to get client: %w", err)
	}

	resp, err := client.Get("https://jsonplaceholder.typicode.com/users")
	if err != nil {
		return fmt.Errorf("unable to fetch users: %w", err)
	}
	// Always be sure to close the response body regardless of whether it is read
	// or not. Closing the body releases resources and memory used for the
	// connection. Closing the body is critical to enable connection reuse.
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("unable to read response body: %w", err)
	}

	var users []User
	err = json.Unmarshal(data, &users)
	if err != nil {
		return fmt.Errorf("unable to unmarshal response body: %w", err)
	}

	for _, user := range users {
		fmt.Printf("ID: %d, Phone: %s\n", user.ID, user.Phone)
	}

	return nil
}
