package main

import (
	"fmt"
	"log"

	daemon "github.com/xelis-project/xelis-go-sdk/daemon"
	sig "github.com/xelis-project/xelis-go-sdk/signature"
	w "github.com/xelis-project/xelis-go-sdk/wallet"
	x "github.com/xelis-project/xelis-go-sdk/xswd"
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

func performHandshake(conn *x.XSWD) error {

	// Define the ApplicationData
	appData := x.ApplicationData{
		ID:          "0000006b2aec4651b82111816ed599d1b72176c425128c66b2ab945552437dc9",
		Name:        "MyXELISApp",
		Description: "An example application integrating with XELIS wallet.",
		Permissions: make(map[string]x.Permission),
	}

	response, err := conn.Authorize(appData)

	if err != nil {
		log.Fatal(err)
	}
	// log.Printf("%+v", response)

	// Check for errors in the response
	if response.Error != nil {
		return fmt.Errorf("handshake error: %v", response.Error)
	}

	log.Printf("Handshake successful: %s", response.Result)
	return nil
}

func QueryWalletAddress(conn *x.XSWD) (address string, err error) {

	address, err = conn.Wallet.GetAddress(w.GetAddressParams{})

	// Check for errors in the response
	if err != nil {
		return "", fmt.Errorf("RPC error: %v", err.Error())
	}

	// Return the wallet address from the result field
	return
}

// QueryPublicKey retrieves the public key for a given wallet address
func QueryPublicKey(conn *daemon.WebSocket, walletAddress string) (string, error) {
	// Prepare the parameters
	params := daemon.ExtractKeyFromAddressParams{
		Address: walletAddress,
		AsHex:   true,
	}

	// Call the method to extract the key
	result, err := conn.ExtractKeyFromAddress(params)
	if err != nil {
		return "", fmt.Errorf("failed to extract key from address: %v", err)
	}

	// Ensure result is a map
	resultMap, ok := result.(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("unexpected result type: %T", result)
	}

	// Extract the "hex" value
	hexValue, ok := resultMap["hex"]
	if !ok {
		return "", fmt.Errorf("hex key not found in the result")
	}

	// Ensure the "hex" value is a string
	hexStr, ok := hexValue.(string)
	if !ok {
		return "", fmt.Errorf("hex value is not a string, got: %T", hexValue)
	}

	return hexStr, nil
}

func main() {
	// Wallet and daemon WebSocket URLs
	walletWebSocketURL := "ws://localhost:44325/xswd"
	daemonWebSocketURL := "wss://xelis-node.mysrv.cloud/json_rpc"

	// Connect to wallet WebSocket
	walletConn, err := x.NewXSWD(walletWebSocketURL)
	if err != nil {
		log.Fatalf("Failed to connect to wallet XSWD: %v", err)
	}
	defer walletConn.Close()

	// Perform handshake with the wallet
	if err := performHandshake(walletConn); err != nil {
		log.Fatalf("Handshake with wallet failed: %v", err)
	}

	// Query the wallet for its address
	address, err := walletConn.Wallet.GetAddress(w.GetAddressParams{})

	if err != nil {
		log.Fatalf("Failed to query wallet address: %v", err)
	}
	log.Printf("Wallet address: %s", address)

	// Connect to daemon WebSocket
	daemonConn, err := daemon.NewWebSocket(daemonWebSocketURL)
	if err != nil {
		log.Fatalf("Failed to connect to daemon websocket: %v", err.Error())
	}
	defer daemonConn.Close()

	// Query the node for the public key
	publicKey, err := QueryPublicKey(daemonConn, address)
	if err != nil {
		log.Fatalf("Failed to query public key: %v", err)
	}
	log.Printf("Public key for wallet %s: %+v", address, publicKey)

	// Data to sign
	// data := map[string]interface{}{"hello": "world"}
	// // Convert the map to JSON bytes
	// dataBytes, err := json.Marshal(data)
	// if err != nil {
	// 	fmt.Println("Error marshaling data:", err)
	// 	return
	// }

	var dataBytes []byte

	// Request the wallet to sign the data
	signature, err := walletConn.Wallet.SignData(dataBytes)
	if err != nil {
		log.Fatalf("Failed to sign data: %v", err)
	}
	log.Printf("Signature: %s", signature)

	// Verify the signature using the public key
	isValid, err := sig.Verify2(publicKey, signature, dataBytes)
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
