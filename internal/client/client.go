// Package client provides HTTP client functionality for communicating with the deCONZ REST API.
// It offers generic functions for making GET, POST, and PUT requests with JSON data,
// and automatically handles serialization and deserialization of request and response data.
package client

import (
	"bytes"
	"encoding/json"
	"net/http"
)

// parseResponse parses an HTTP response body into the specified type.
// This is a generic helper function used by the public request functions.
//
// Type Parameters:
//   - R: The type to parse the response into
//
// Parameters:
//   - resp: The HTTP response to parse
//
// Returns:
//   - *R: A pointer to the parsed response data
//   - error: An error if the response could not be parsed
func parseResponse[R interface{}](resp *http.Response) (*R, error) {
	responseData := new(R)
	if err := json.NewDecoder(resp.Body).Decode(&responseData); err != nil {
		return nil, err
	}

	return responseData, nil
}

// Post makes an HTTP POST request with JSON data and parses the response.
// This function is used for creating resources or requesting actions from the deCONZ API.
//
// Type Parameters:
//   - R: The type to parse the response into
//
// Parameters:
//   - url: The URL to send the request to
//   - data: The data to send in the request body (will be serialized to JSON)
//
// Returns:
//   - *R: A pointer to the parsed response data
//   - error: An error if the request failed or the response could not be parsed
func Post[R interface{}](url string, data any) (*R, error) {
	// Serialize the request data to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	// Send the POST request
	resp, err := http.Post(url, "application/json", bytes.NewReader(jsonData))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Parse the response
	return parseResponse[R](resp)
}

// Put makes an HTTP PUT request with JSON data and parses the response.
// This function is used for updating resources in the deCONZ API.
//
// Type Parameters:
//   - R: The type to parse the response into
//
// Parameters:
//   - url: The URL to send the request to
//   - data: The data to send in the request body (will be serialized to JSON)
//
// Returns:
//   - *R: A pointer to the parsed response data
//   - error: An error if the request failed or the response could not be parsed
func Put[R interface{}](url string, data any) (*R, error) {
	// Serialize the request data to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	// Create a new PUT request
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(jsonData))
	if err != nil {
		return nil, err
	}

	// Set the content type header
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Parse the response
	return parseResponse[R](resp)
}

// Get makes an HTTP GET request and parses the response.
// This function is used for retrieving resources from the deCONZ API.
//
// Type Parameters:
//   - R: The type to parse the response into
//
// Parameters:
//   - url: The URL to send the request to
//
// Returns:
//   - *R: A pointer to the parsed response data
//   - error: An error if the request failed or the response could not be parsed
func Get[R interface{}](url string) (*R, error) {
	// Send the GET request
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Parse the response
	return parseResponse[R](resp)
}
