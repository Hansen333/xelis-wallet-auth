# XELIS Wallet Authentication

## Overview

This application demonstrates a simple authentication process with the XELIS blockchain wallet. It securely verifies the ownership of a wallet address and validates its ability to sign data. The process ensures that the wallet and its associated public key are legitimate, and that the data signed by the wallet can be verified using the corresponding public key.

## Process Outline

### 1. Handshake with the Wallet
- Establish a WebSocket connection with the wallet.
- Send application-specific metadata during a handshake process to authorize interaction with the wallet.

### 2. Query the Wallet Address
- Retrieve the wallet's address using the wallet API.
- The wallet address serves as an identifier and is associated with the public key used for cryptographic operations.

### 3. Obtain the Public Key
- Query the blockchain daemon to extract the public key corresponding to the wallet address.
- The public key is required for verifying signatures created by the wallet.

### 4. Sign Data Using the Wallet
- Request the wallet to sign a specific piece of data.
- The wallet generates a signature using its private key, which corresponds to the retrieved public key.

### 5. Verify the Signature
- Use the retrieved public key to verify the authenticity of the signature.
- This step confirms that the wallet indeed signed the data, thereby authenticating its ownership.

This authentication process leverages standard cryptographic techniques, ensuring secure communication and validation of wallet ownership without exposing private keys or sensitive information.

## Run Instructions

To run the application, execute the following command in your terminal:

```bash
go run main.go wallet.go daemon.go utils.go
```