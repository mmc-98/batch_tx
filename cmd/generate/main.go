package main

import (
	"fmt"

	"github.com/okx/go-wallet-sdk/example"
)

func main() {
	mnemonic, err := example.GenerateMnemonic()
	if err != nil {
		panic(err)
	}
	fmt.Printf("mnemonic: %v\n", mnemonic)
	hdPath := example.GetDerivedPath(0)
	derivePrivateKey, err := example.GetDerivedPrivateKey(mnemonic, hdPath)

	// get new address
	newAddress := example.GetNewAddress(derivePrivateKey)
	fmt.Printf("master address: %v\n", newAddress)
}
