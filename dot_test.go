package gsrpc_test

import (
	"encoding/hex"
	"fmt"
	"github.com/centrifuge/go-substrate-rpc-client/v2/extra"
	rpc "github.com/centrifuge/go-substrate-rpc-client/v2/method"
	"math/big"
	"testing"
)

func TestSystem(t *testing.T) {
	//dotsdk, err := rpc.NewDotSdk("wss://rpc.polkadot.io")
	dotsdk, err := rpc.NewDotSdk("wss://westend-rpc.polkadot.io")
	//dotsdk, err := rpc.NewDotSdk("ws://127.0.0.1:9944")
	if err != nil {
		panic(err)
	}

	propreties, err := dotsdk.Symbol()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", propreties)

	err = dotsdk.CheckChain(rpc.Polkadot)
	if err != nil {
		panic(err)
	}

	height, err := dotsdk.GetHieght()
	if err != nil {
		panic(err)
	}
	fmt.Printf("当前节点高度:%v\n", height)

	pair, err := dotsdk.GenerateNewAddress(rpc.Polkadot)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", pair)

	err = dotsdk.ValidAddress("12xtAYsRUrmbniiWQqJtECiBQrMn8AypQcXhnQAc6RB6XkLW", rpc.Polkadot)
	if err != nil {
		panic(err)
	}

	account, err := dotsdk.GetAccountInfo("12xtAYsRUrmbniiWQqJtECiBQrMn8AypQcXhnQAc6RB6XkLW")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", account)

	txs, err := dotsdk.GetBlockTransaction(4124848, rpc.Polkadot)
	if err != nil {
		panic(err)
	}
	for _, tx := range txs {
		fmt.Printf("%+v\n", tx)
	}

	hash, fee, err := dotsdk.Transfer("0xebc62b16438a7009ced6be3953474ffd8d1b884c7b8a5e359f717e55a0faf018", rpc.Westend, "5EqHe4TTbtqD7zCBrV8HpQWjekEEmFjb88m1tDGGVipvC1Pg", big.NewInt(200000000000), 9)
	if err != nil {
		panic(err)
	}
	fmt.Printf("tx hash is %v, pay fee is %v", hash.Hex(), fee)

	fee, err = dotsdk.GetPartialFee("0x412c4a3c51e6dd7b4bc60567694d2f6f0d4075217aa05166ed7dab616828dea5", rpc.Polkadot, "5EqHe4TTbtqD7zCBrV8HpQWjekEEmFjb88m1tDGGVipvC1Pg", big.NewInt(2000000000000), 9)
	if err != nil {
		panic(err)
	}
	fmt.Printf("fee1: %+v\n", fee)

	fee2, err := dotsdk.GetPartialFee("0x412c4a3c51e6dd7b4bc60567694d2f6f0d4075217aa05166ed7dab616828dea5", rpc.Polkadot, "5EqHe4TTbtqD7zCBrV8HpQWjekEEmFjb88m1tDGGVipvC1Pg", big.NewInt(2000000000000), 9)
	if err != nil {
		panic(err)
	}
	fmt.Printf("fee2: %+v\n", fee2)

	err = dotsdk.CheckSync()
	if err != nil {
		panic(err)
	}

}

func TestAddressFromPub(t *testing.T) {
	pub := "a69958eee5de0cb8fb250eba9c4b4ab1675468e68e49a5ebcac22fa9340fe938"
	pubBytes, _ := hex.DecodeString(pub)
	fmt.Println(pubBytes)
	//pubBytes = append([]byte{0xff},pubBytes...)
	address, err := extra.Encode(pubBytes, extra.PolkadotPrefix)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(address)
}

func TestVliadAddress(t *testing.T) {
	err := extra.VerityAddress("5HpZZQKbdd8iUnedUCASdY9xyAMrrzmsmWY5J4nzTBjbbnbX", []byte{42})
	if err != nil {
		panic(err)
	}
}

func TestTransferECDSA(t *testing.T) {
	dotsdk, err := rpc.NewDotSdk("wss://westend-rpc.polkadot.io")
	if err != nil {
		t.Fatal(err)
	}

	//types.SetSerDeOptions(types.SerDeOptions{NoPalletIndices: false})

	hash, fee, err := dotsdk.TransferECDSA("52feabaab39178b12be74ff1f7cdcf18812d68ce2dfb648a5783be4e92924220", rpc.Westend, "5FvcbAPyzmMsYrGgmecLG6vPMqzaCr1V91hngtuYddc8VL3s", big.NewInt(100000000000), 0)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("tx hash is %v, pay fee is %v", hash.Hex(), fee)
}
