package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

// QueryPublicKey retrieves the public key for a given wallet address
func QueryPublicKey(conn *websocket.Conn, walletAddress string) (string, error) {

	// Define the parameters
	params := map[string]interface{}{
		"address": walletAddress,
		"as_hex":  false, // Set to true to get the public key as hex
	}

	// Create a JSON-RPC request
	request := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "extract_key_from_address",
		"params":  params,
	}

	// Log the outgoing request for debugging
	// requestJSON, _ := json.MarshalIndent(request, "", "  ")
	// log.Printf("Sending Request: %s", string(requestJSON))

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
		JSONRPC string `json:"jsonrpc"`
		ID      int    `json:"id"`
		Result  struct {
			Bytes []byte `json:"bytes"`
		} `json:"result"`
		Error map[string]interface{} `json:"error,omitempty"`
	}
	if err := json.Unmarshal(message, &response); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	// Check for errors in the response
	if response.Error != nil {
		return "", fmt.Errorf("RPC error: %v", response.Error)
	}

	// Convert the byte array to a hexadecimal string
	publicKey := hex.EncodeToString(response.Result.Bytes)

	// Return the public key
	return publicKey, nil
}
