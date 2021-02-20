package config

import (
	"testing"

	"github.com/btceasypay/btp-desktop-wallet/rpc/client"
)

func TestSave(t *testing.T) {

	cfg := NewDefaultConfig()

	cfg.BtpdList = make(map[string]*client.Config)
	cfg.BtpdList["local"] = &client.Config{
		RPCServer: "2.2.2.2:8080",
	}
	cfg.BtpdSelect = "local"

	t.Log(cfg.Save("config.toml"))
}
