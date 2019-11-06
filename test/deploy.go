package main

import (
  "context"
  "crypto/ecdsa"
  "fmt"
  "log"
  "math/big"
  "github.com/ethereum/go-ethereum/accounts/abi/bind"
  "github.com/ethereum/go-ethereum/crypto"
  "github.com/ethereum/go-ethereum/ethclient"
  instapay "github.com/sslab-instapay/instapay-go-server/contracts"
)

func main() {
  client, err := ethclient.Dial("ws://141.223.121.139:8881")
  if err != nil {
    log.Fatal(err)
  }

  privateKey, err := crypto.HexToECDSA("e113ff405699b7779fbe278ee237f2988b1e6769d586d8803860d49f28359fbd")
  if err != nil {
    log.Fatal(err)
  }

  publicKey := privateKey.Public()
  publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
  if !ok {
    log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
  }

  address := crypto.PubkeyToAddress(*publicKeyECDSA)

  nonce, err := client.PendingNonceAt(context.Background(), address)
  if err != nil {
    log.Fatal(err)
  }

  gasPrice, err := client.SuggestGasPrice(context.Background())
  if err != nil {
    log.Fatal(err)
  }

  auth := bind.NewKeyedTransactor(privateKey)
  auth.Nonce = big.NewInt(int64(nonce))
  auth.Value = big.NewInt(0)     // in wei
  auth.GasLimit = uint64(10000000) // in units
  auth.GasPrice = gasPrice

  address, tx, instance, err := instapay.DeployInstapay(auth, client)
  if err != nil {
    log.Fatal(err)
  }

  fmt.Println(address.Hex())
  fmt.Println(tx.Hash().Hex())

  _ = instance
}
