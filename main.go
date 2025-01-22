package main

import (
	"log"
)

// Authentication Process Overview
//
// The idea behind performing a simple authentication in this application is to securely verify
// the ownership of a wallet address and validate its ability to sign data. This process ensures
// that the wallet and its associated public key are legitimate and that the data signed by the wallet
// can be verified using the corresponding public key.
//
// Process Outline:
// 1. Handshake with the Wallet:
//    - Establish a WebSocket connection with the wallet and send application-specific metadata during a handshake process.
//      This step ensures that the application is authorized to interact with the wallet.
//
// 2. Query the Wallet Address:
//    - Retrieve the wallet's address using the wallet API. This address serves as the identifier for the wallet
//      and is associated with the public key used for cryptographic operations.
//
// 3. Obtain the Public Key:
//    - Query the blockchain daemon to extract the public key corresponding to the wallet address.
//      The public key is necessary for verifying signatures created by the wallet.
//
// 4. Sign Data Using the Wallet:
//    - Request the wallet to sign a specific piece of data. The wallet generates a signature using its private key,
//      which corresponds to the retrieved public key.
//
// 5. Verify the Signature:
//    - Use the retrieved public key to verify the authenticity of the signature. This step confirms that the wallet
//      indeed signed the data, thereby authenticating its ownership.
//
// This authentication process leverages standard cryptographic techniques, ensuring secure communication
// and validation of wallet ownership without exposing private keys or sensitive information.

func main() {
	// Wallet and daemon WebSocket URLs
	walletWebSocketURL := "ws://localhost:44325/xswd"
	daemonWebSocketURL := "wss://xelis-node.mysrv.cloud/json_rpc"

	// Connect to wallet WebSocket
	walletConn := connectToWebSocket(walletWebSocketURL)
	defer walletConn.Close()

	// Perform handshake with the wallet
	if err := performHandshake(walletConn); err != nil {
		log.Fatalf("Handshake with wallet failed: %v", err)
	}

	// Query the wallet for its address
	address, err := QueryWalletAddress(walletConn)
	if err != nil {
		log.Fatalf("Failed to query wallet address: %v", err)
	}
	log.Printf("Wallet address: %s", address)

	// Connect to daemon WebSocket
	daemonConn := connectToWebSocket(daemonWebSocketURL)
	defer daemonConn.Close()

	// Query the node for the public key
	publicKey, err := QueryPublicKey(daemonConn, address)
	if err != nil {
		log.Fatalf("Failed to query public key: %v", err)
	}
	log.Printf("Public key for wallet %s: %s", address, publicKey)

	// Data to sign
	data := "Hello, XELIS!"

	// Request the wallet to sign the data
	signature, err := SignData(walletConn, data)
	if err != nil {
		log.Fatalf("Failed to sign data: %v", err)
	}
	log.Printf("Signature: %s", signature)

	// Verify the signature using the public key
	isValid, err := VerifySignature(publicKey, signature, data)
	if err != nil {
		log.Fatalf("Failed to verify signature: %v", err)
	}
	if isValid {
		log.Println("Signature is valid!")
	} else {
		log.Println("Signature is invalid!")
	}
	// Keep the application running
	select {}
}
