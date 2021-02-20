package btpd

import (
	"fmt"

	"github.com/btceasypay/bitcoinpay/log"

	"github.com/btceasypay/btp-desktop-wallet/config"
	"github.com/btceasypay/btp-desktop-wallet/rpc/client"
	"github.com/btceasypay/btp-desktop-wallet/wallet"
)

// API to mgr btpd
type API struct {
	cfg  *config.Config
	btpd *Btpd
}

// NewAPI api make
func NewAPI(cfg *config.Config, btpd *Btpd) *API {
	return &API{
		cfg:  cfg,
		btpd: btpd,
	}
}

// List list all btpd
func (api *API) List() ([]*client.Config, error) {
	return api.cfg.Btpds, nil
}

// Add a new btpd conf
func (api *API) Add(name string,
	RPCServer string, RPCUser string, RPCPassword string,
	RPCCert string, NoTLS bool, TLSSkipVerify bool,
	Proxy string, ProxyUser string, ProxyPass string) error {

	for _, item := range api.cfg.Btpds {
		if item.Name == name {
			return nil
		}
	}

	api.cfg.Btpds = append(api.cfg.Btpds, &client.Config{
		Name:          name,
		RPCUser:       RPCUser,
		RPCPassword:   RPCPassword,
		RPCServer:     RPCServer,
		RPCCert:       RPCCert,
		NoTLS:         NoTLS,
		TLSSkipVerify: TLSSkipVerify,
		Proxy:         Proxy,
		ProxyUser:     ProxyUser,
		ProxyPass:     ProxyPass,
	})

	api.cfg.Save(api.cfg.ConfigFile)

	return nil
}

// Del a btpd conf
func (api *API) Del(name string) error {
	nameP := -1
	for i, item := range api.cfg.Btpds {
		if item.Name == name {
			nameP = i
			break
		}
	}
	if nameP == -1 {
		return nil
	}

	api.cfg.Btpds = append(api.cfg.Btpds[:nameP], api.cfg.Btpds[nameP+1:]...)
	api.cfg.Save(api.cfg.ConfigFile)
	return nil
}

// Update btpd conf
func (api *API) Update(name string,
	RPCServer string, RPCUser string, RPCPassword string,
	RPCCert string, NoTLS bool, TLSSkipVerify bool,
	Proxy string, ProxyUser string, ProxyPass string) error {

	var updateBtpd *client.Config
	for _, item := range api.cfg.Btpds {
		if item.Name == name {
			updateBtpd = item
			break
		}
	}

	updateBtpd.RPCUser = RPCUser
	updateBtpd.RPCPassword = RPCPassword
	updateBtpd.RPCServer = RPCServer
	updateBtpd.RPCCert = RPCCert
	updateBtpd.NoTLS = NoTLS
	updateBtpd.TLSSkipVerify = TLSSkipVerify
	updateBtpd.Proxy = Proxy
	updateBtpd.ProxyUser = ProxyUser
	updateBtpd.ProxyPass = ProxyPass

	api.cfg.Save(api.cfg.ConfigFile)

	return nil
}

// Reset btpd rpc client
func (api *API) Reset(name string) error {

	if api.cfg.BtpdSelect == name {
		log.Trace("not reset btpd,it eq")
		return nil
	}

	var resetBtpd *client.Config
	for _, item := range api.cfg.Btpds {
		if item.Name == name {
			resetBtpd = item
			break
		}
	}

	if resetBtpd == nil {
		return fmt.Errorf("btpd %s not found", name)
	}

	hc, err := wallet.NewHtpcByCfg(resetBtpd)
	if err != nil {
		return fmt.Errorf("make rpc clent error: %s", err.Error())
	}

	api.cfg.BtpdSelect = name
	api.btpd.Status.CurrentName = name
	// update wallet httpclient
	api.btpd.Wt.HttpClient = hc

	return nil
}

// Status get btpd stats
func (api *API) Status() (*Status, error) {
	return api.btpd.Status, nil
}

//Status btpd status
type Status struct {
	Network      string
	CurrentName  string //current btpd name
	err          string //
	MainOrder    uint32
	MainHeight   uint32
	Blake2bdDiff string // float64
	CuckarooDiff float64
	CuckatooDiff float64
}
