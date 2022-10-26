package eth_test

import (
	"cloud-run-sample/eth"
	"encoding/hex"
	"fmt"
	"testing"
	"time"
)

func TestSign(t *testing.T) {
	//	SAVE BUT DO NOT SHARE THIS (Private Key): 0x04b6f500172977d66c285dadfea69dc75bf451ead1075ef4fa3197d933178318
	//	Public Key: 0x0480f8258d856bd4df28035439df9deb89b8b76fd0e589399d11f907a3a65783bd4f0d81ce831841ad56865420429ff98f9440791485c252d2071571a4734574da
	//  Address: 0xc08c89c9cd320847AED4ec5AA04338D7e932E53D
	privateKey, address := eth.WalletFromPrivate("04b6f500172977d66c285dadfea69dc75bf451ead1075ef4fa3197d933178318")

	fmt.Println("Address:", address.Hex())

	timestamp := int(time.Now().Unix())
	messageReplaced := fmt.Sprintf("%s::%d", "2", timestamp)

	hashMessage := eth.SignMessage([]byte(messageReplaced))
	signature := eth.Sign(hashMessage, privateKey)

	signatureHex := "0x" + hex.EncodeToString(signature)

	fmt.Println(eth.RecoverSig(signatureHex, []byte(messageReplaced)).Hex())

	fmt.Println(eth.VerifySig(address.Hex(), signatureHex, []byte(messageReplaced)))

	fmt.Printf("%s::%s::%s::%d", address.Hex(), "2", signatureHex, timestamp)
}
