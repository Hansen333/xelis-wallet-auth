package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

// ApplicationData defines the structure for application-specific information
type ApplicationData struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	URL         string            `json:"url"`
	Permissions map[string]string `json:"permissions"`
	Signature   *string           `json:"signature"` // Optional, can be null
}

// RPCRequest defines the structure of a JSON-RPC request
type RPCRequest struct {
	JSONRPC         string           `json:"jsonrpc"`
	ID              int              `json:"id"`
	Method          string           `json:"method"`
	Params          any              `json:"params,omitempty"`
	ApplicationData *ApplicationData `json:"ApplicationData,omitempty"`
}

func performHandshake(conn *websocket.Conn) error {
	// Define permissions (empty for this example)
	permissions := map[string]string{}

	// Define the ApplicationData
	appData := ApplicationData{
		ID:          "0000006b2aec4651b82111816ed599d1b72176c425128c66b2ab945552437dc9",
		Name:        "MyXELISApp",
		Description: "An example application integrating with XELIS wallet.",
		URL:         "https://myxelisapp.example.com",
		Permissions: permissions,
		Signature:   nil,
	}

	// Send ApplicationData as the initial message
	if err := conn.WriteJSON(appData); err != nil {
		return fmt.Errorf("failed to send ApplicationData: %w", err)
	}

	// Wait for the response
	_, message, err := conn.ReadMessage()
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	// Log the raw response for debugging
	log.Printf("Raw response: %s", string(message))

	// Parse the response
	var response struct {
		ID      any    `json:"id"`
		JSONRPC string `json:"jsonrpc"`
		Result  struct {
			Message string `json:"message"`
			Success bool   `json:"success"`
		} `json:"result"`
		Error map[string]interface{} `json:"error,omitempty"`
	}
	if err := json.Unmarshal(message, &response); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	// Check for errors in the response
	if response.Error != nil {
		return fmt.Errorf("handshake error: %v", response.Error)
	}

	// Verify successful handshake
	if !response.Result.Success {
		return fmt.Errorf("handshake failed: %s", response.Result.Message)
	}

	log.Printf("Handshake successful: %s", response.Result.Message)
	return nil
}
func QueryWalletAddress(conn *websocket.Conn) (string, error) {
	// Create a JSON-RPC request for the "wallet.get_address" method
	request := RPCRequest{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "wallet.get_address",
	}

	// Send the request
	if err := conn.WriteJSON(request); err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}

	// Wait for the response
	_, message, err := conn.ReadMessage()
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	// Log the raw response for debugging
	log.Printf("Raw response: %s", string(message))

	// Parse the response
	var response struct {
		JSONRPC string                 `json:"jsonrpc"`
		ID      int                    `json:"id"`
		Result  string                 `json:"result"`
		Error   map[string]interface{} `json:"error,omitempty"`
	}
	if err := json.Unmarshal(message, &response); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	// Check for errors in the response
	if response.Error != nil {
		return "", fmt.Errorf("RPC error: %v", response.Error)
	}

	// Return the wallet address from the result field
	return response.Result, nil
}
func SignData(conn *websocket.Conn, data string) (string, error) {
	// Create a JSON-RPC request for the "sign_data" method
	request := RPCRequest{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "wallet.sign_data",
		Params: map[string]string{
			"data": data,
		},
	}

	// Send the request
	if err := conn.WriteJSON(request); err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}

	// Wait for the response
	_, message, err := conn.ReadMessage()
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	// Log the raw response for debugging
	log.Printf("Raw response: %s", string(message))

	// Parse the response
	var response struct {
		JSONRPC string                 `json:"jsonrpc"`
		ID      int                    `json:"id"`
		Result  string                 `json:"result"` // Parse result as a plain string
		Error   map[string]interface{} `json:"error,omitempty"`
	}
	if err := json.Unmarshal(message, &response); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	// Check for errors in the response
	if response.Error != nil {
		return "", fmt.Errorf("RPC error: %v", response.Error)
	}

	// Return the signature from the result
	return response.Result, nil
}
