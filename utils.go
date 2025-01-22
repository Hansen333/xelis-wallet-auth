package main

import (
	"encoding/hex"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
	"golang.org/x/crypto/ed25519"
)

// connectToWebSocket establishes a WebSocket connection to the specified URL
func connectToWebSocket(url string) *websocket.Conn {
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatalf("Failed to connect to WebSocket at %s: %v", url, err)
	}
	log.Printf("Connected to WebSocket at %s", url)
	return conn
}

func VerifySignature(publicKeyHex, signatureHex, data string) (bool, error) {
	// Decode the public key and signature from hexadecimal
	publicKey, err := hex.DecodeString(publicKeyHex)
	if err != nil {
		return false, fmt.Errorf("failed to decode public key: %w", err)
	}

	if len(publicKey) != 32 {
		return false, fmt.Errorf("invalid public key length: expected 32 bytes, got %d", len(publicKey))
	}

	signature, err := hex.DecodeString(signatureHex)
	if err != nil {
		return false, fmt.Errorf("failed to decode signature: %w", err)
	}

	// Log the decoded values for debugging
	log.Printf("Decoded Public Key: %x", publicKey)
	log.Printf("Decoded Signature: %x", signature)
	log.Printf("Data to Verify: %s", data)

	// Verify the signature
	isValid := ed25519.Verify(publicKey, []byte(data), signature)
	if !isValid {
		log.Println("Signature verification failed!")
	}
	return isValid, nil
}
