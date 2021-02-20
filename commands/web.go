package commands

import (
	"fmt"
	"github.com/btceasypay/bitcoinpay/log"
	"github.com/btceasypay/btp-desktop-wallet/config"
	"github.com/btceasypay/btp-desktop-wallet/wserver"
	"github.com/spf13/cobra"
)

// web mode

var WebCmd = &cobra.Command{
	Use:   "web",
	Short: "web administration UI",
	Example: `
		Enter web mode
		`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("web model")
		btpMain(userConf)
	},
}

func AddWebCommand() {

}

func btpMain(cfg *config.Config) {
	log.Trace("Btp Main")
	wsvr, err := wserver.NewWalletServer(cfg)
	if err != nil {
		log.Error(fmt.Sprintf("NewWalletServer err: %s", err))
		return
	}
	wsvr.Start()

	exitCh := make(chan int)
	<-exitCh
}
