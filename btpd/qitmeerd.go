package btpd

import (
	"fmt"
	"strconv"
	"time"

	"github.com/btceasypay/btp-desktop-wallet/wallet"

	"github.com/btceasypay/bitcoinpay/log"
)

// Btpd mgr
type Btpd struct {
	Status *Status // *qJson.InfoNodeResult
	Wt     *wallet.Wallet
}

// NewBtpd make btpd
func NewBtpd(wt *wallet.Wallet, name string) *Btpd {
	d := &Btpd{
		Wt:     wt,
		Status: &Status{Network: wt.ChainParams().Name, CurrentName: name},
	}
	d.Start()
	return d
}

// Start run
func (btpd *Btpd) Start() {
	go btpd.GetStatus()
}

// GetStatus get current btpd status
func (btpd *Btpd) GetStatus() {
	defer func() {
		if rev := recover(); rev != nil {
			go btpd.GetStatus()
			log.Error("btpd GetStatus recover", "recover", rev)
		}
	}()

	ticker := time.NewTicker(30 * time.Second)
	for {
		select {
		case <-ticker.C:
			if btpd.Wt.HttpClient == nil {
				log.Debug("btpd GetNodeInfo,but HttpClient nil")
				continue
			}

			nodeInfo, err := btpd.Wt.HttpClient.GetNodeInfo()
			if err != nil {
				btpd.Status.err = fmt.Sprintf("getNodeInfo err: %v", err)
				log.Error("btpd GetNodeInfo err", "err", err)
				continue
			}

			btpd.Status.err = ""
			btpd.Status.MainOrder = nodeInfo.GraphState.MainOrder
			btpd.Status.MainHeight = nodeInfo.GraphState.MainHeight
			btpd.Status.Blake2bdDiff = strconv.FormatFloat(nodeInfo.PowDiff.Blake2bdDiff, 'f', 2, 64)
			btpd.Status.CuckarooDiff = nodeInfo.PowDiff.CuckarooDiff
			btpd.Status.CuckatooDiff = nodeInfo.PowDiff.CuckatooDiff
		}
	}
}
