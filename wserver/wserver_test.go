package wserver

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/btceasypay/bitcoinpay/crypto/bip32"
	"github.com/btceasypay/bitcoinpay/crypto/bip39"
	"github.com/btceasypay/bitcoinpay/crypto/ecc/secp256k1"
	"github.com/btceasypay/bitcoinpay/crypto/seed"
	"github.com/btceasypay/btp-desktop-wallet/utils"

	"github.com/btceasypay/bitcoinpay/bx"
)

func TestSeed(t *testing.T) {

	activeNetParams := utils.GetNetParams("testnet")

	seedBuf, err := seed.GenerateSeed(uint16(32))
	if err != nil {
		t.Log(fmt.Errorf("GenerateSeed err: %s", err))
		return
	}
	seed := hex.EncodeToString(seedBuf)
	t.Log("seed", seed)

	mnemonic, err := bip39.NewMnemonic(seedBuf)

	t.Log("mnemonic", mnemonic, err)

	s3, err := bip39.EntropyFromMnemonic(mnemonic)
	t.Log("ok", hex.EncodeToString(s3) == seed, err)

	//import master key addr
	seedKey, err := bip32.NewMasterKey(seedBuf)
	if err != nil {
		t.Logf("createWallet NewMasterKey err: %s", err)
		return
	}
	t.Logf("createWallet import master key: %x\n", seedKey.Key)

	pri, _ := secp256k1.PrivKeyFromBytes(seedKey.Key)
	wif, err := utils.NewWIF(pri, activeNetParams, true)
	if err != nil {
		t.Logf("createWallet private key decode failed: %s", err)
		return
	}
	if !wif.IsForNet(activeNetParams) {
		t.Logf("createWallet Key is not intended for: %s %s", activeNetParams.Name, err)
		return
	}

	t.Log(wif)

}

func TestGan(t *testing.T) {
	//func GenerateAddr() string {

	ver := "mainnet"

	s, err := bx.NewEntropy(32)
	if err != nil {
		t.Logf("An error occurred generating s，%s", err)
	}
	prv, err := bx.EcNew("secp256k1", s)
	if err != nil {
		t.Logf("An error occurred generating private key，%s", err)
	}
	pub, err := bx.EcPrivateKeyToEcPublicKey(false, prv)
	if err != nil {
		t.Logf("An error occurred generating public key，%s", err)
	}
	addr, err := bx.EcPubKeyToAddress(ver, pub)
	if err != nil {
		t.Logf("An error occurred generating address，%s", err)
	}

	t.Log(addr)

}
