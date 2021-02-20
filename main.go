package main

import (
	"github.com/btceasypay/bitcoinpay/log"
	"github.com/btceasypay/btp-desktop-wallet/commands"
	"os"
)

func main() {
	if err := commands.RootCmd.Execute(); err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
}
