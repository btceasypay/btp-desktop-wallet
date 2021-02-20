package wserver

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/julienschmidt/httprouter"

	qJson "github.com/btceasypay/bitcoinpay/core/json"
	"github.com/btceasypay/bitcoinpay/log"

	"github.com/btceasypay/btp-desktop-wallet/assets"
	"github.com/btceasypay/btp-desktop-wallet/btpd"
	"github.com/btceasypay/btp-desktop-wallet/config"
	"github.com/btceasypay/btp-desktop-wallet/rpc/client"
	"github.com/btceasypay/btp-desktop-wallet/rpc/server"
	"github.com/btceasypay/btp-desktop-wallet/utils"
	"github.com/btceasypay/btp-desktop-wallet/wallet"
)

//WalletServer wallet api server
type WalletServer struct {
	cfg *config.Config

	WtLoader *wallet.Loader
	Wt       *wallet.Wallet

	RPCSvr *server.RpcServer

	exitCh chan bool

	BtpdStatus *qJson.InfoNodeResult
}

//NewWalletServer make a wallet api server
func NewWalletServer(cfg *config.Config) (wSvr *WalletServer, err error) {

	// qitmeed
	var btpdSelect *client.Config
	if cfg.BtpdSelect != "" {
		for _, item := range cfg.Btpds {
			if item.Name == cfg.BtpdSelect {
				btpdSelect = item
			}
		}
	}
	if len(cfg.Btpds) < 1 {
		return nil, fmt.Errorf("config btpds not found %s", "")
	}
	if btpdSelect == nil {
		cfg.BtpdSelect = cfg.Btpds[0].Name
	}

	activeNetParams := utils.GetNetParams(cfg.Network)
	dbDir := filepath.Join(cfg.AppDataDir, cfg.Network)
	wtLoader := wallet.NewLoader(activeNetParams, dbDir, 250, cfg)
	wtExist, err := wtLoader.WalletExists()
	if err != nil {
		return nil, fmt.Errorf("load wallet err: %s", err)
	}
	if !wtExist && !cfg.UI {
		return nil, fmt.Errorf("not wallet exist,please run crate command")
	}

	wSvr = &WalletServer{
		cfg: cfg,

		WtLoader: wtLoader,
		exitCh:   make(chan bool),
	}

	RPCSvrCfg := &server.Config{
		RPCUser:       cfg.RPCUser,
		RPCPass:       cfg.RPCPass,
		RPCCert:       cfg.RPCCert,
		RPCKey:        cfg.RPCKey,
		RPCMaxClients: 100,
		DisableRPC:    false,
		DisableTLS:    cfg.DisableTLS,
	}

	wSvr.RPCSvr, err = server.NewRPCServer(RPCSvrCfg)
	if err != nil {
		return nil, fmt.Errorf("NewWallet: %s", err)
	}

	if wSvr.cfg.UI {
		//ui rpc
		wSvr.RPCSvr.RegisterService("ui", NewAPI(cfg, wSvr))
	} else {
		err = wSvr.OpenWallet(cfg.WalletPass)
		if err != nil {
			return nil, err
		}
	}

	return
}

func (wSvr *WalletServer) run() {
	defer func() {
		if rev := recover(); rev != nil {
			log.Trace("WalletServer.run", "WalletServer run recover ", rev)
			go wSvr.run()
		}
	}()
	go func() {
		for {
			select {
			case <-wSvr.exitCh:
				os.Exit(1)
			}
		}
	}()
	log.Trace("WalletServer run")

	router := httprouter.New()

	if wSvr.cfg.UI {
		staticF, err := assets.GetStatic()
		if err != nil {
			log.Error("server run ", "err ", err)
			return
		}
		myStaticF := assets.NewMyStatic(staticF)

		myStaticF.AddFilter("/config.js", func() []byte {

			//update config.js
			tmpl := `
			//config
			window.BtpConfig = {
				RPCAddr: "{{api_url}}",
				RPCUser: "{{rpc_user}}",
				RPCPass: "{{rpc_pass}}"
			};
			`
			tmpl = strings.Replace(tmpl, "{{api_url}}", "http://"+wSvr.cfg.Listeners[0]+"/api", -1)
			tmpl = strings.Replace(tmpl, "{{rpc_user}}", wSvr.cfg.RPCUser, -1)
			tmpl = strings.Replace(tmpl, "{{rpc_pass}}", wSvr.cfg.RPCPass, -1)

			return []byte(tmpl)
		})

		router.ServeFiles("/app/*filepath", myStaticF)
		router.GET("/", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
			http.Redirect(w, r, "app/index.html", http.StatusMovedPermanently)
		})

		//ajx post options
		router.OPTIONS("/api", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
			if r.Header.Get("Origin") == "http://127.0.0.1:8080" {
				w.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:8080")
			} else {
				w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
			}
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Methods", "GET,HEAD,OPTIONS,POST,PUT")
			w.Header().Set("Access-Control-Allow-Headers", "Access-Control-Allow-Headers, Origin, Accept, X-Requested-With, Content-Type, Access-Control-Request-Method, Access-Control-Request-Headers, Authorization")
			return
		})
	}

	router.POST("/api", wSvr.HandleAPI)

	for _, addr := range wSvr.cfg.Listeners {
		go func() {
			log.Trace("WalletServer listening on", "addr", addr)
			err := http.ListenAndServe(addr, router)
			if err != nil {
				log.Error("server listen", " err", err)
				wSvr.exitCh <- true
				return
			}
		}()
	}

}

// HandleAPI RPC Method
func (wSvr *WalletServer) HandleAPI(ResW http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if r.Header.Get("Origin") == "http://127.0.0.1:8080" {
		ResW.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:8080")
	} else {
		ResW.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
	}

	ResW.Header().Set("Access-Control-Allow-Credentials", "true")
	ResW.Header().Set("Access-Control-Allow-Methods", "GET,HEAD,OPTIONS,POST,PUT")
	ResW.Header().Set("Access-Control-Allow-Headers", "Access-Control-Allow-Headers, Origin, Accept, X-Requested-With, Content-Type, Access-Control-Request-Method, Access-Control-Request-Headers, Authorization")

	wSvr.RPCSvr.HandleFunc(ResW, r)
}

// Start server
func (wSvr *WalletServer) Start() error {
	log.Trace("WalletServer start")

	wSvr.RPCSvr.Start()

	go wSvr.run()

	//open home in web browser
	if wSvr.cfg.UI {
		utils.OpenBrowser("http://" + wSvr.cfg.Listeners[0])
	}

	return nil
}

// RegAPI if wallet open
func (wSvr *WalletServer) RegAPI() {
	//wallet rpc
	wSvr.RPCSvr.RegisterService("wallet", wallet.NewAPI(wSvr.cfg, wSvr.Wt))

	//btpd rpc
	btpD := btpd.NewBtpd(wSvr.Wt, wSvr.cfg.BtpdSelect)
	wSvr.RPCSvr.RegisterService("btpd", btpd.NewAPI(wSvr.cfg, btpD))
}

// OpenWallet load wallet and start rpc
func (wSvr *WalletServer) OpenWallet(pass string) error {
	if wSvr.Wt != nil {
		log.Trace("OpenWallet: wallet already open")
		return nil
	}
	walletPubPassBuf := []byte(pass)
	wt, err := wSvr.WtLoader.OpenExistingWallet(walletPubPassBuf, false)
	if err != nil {
		return fmt.Errorf("OpenWallet OpenExistingWallet err: %s", err)
	}
	wSvr.Wt = wt
	log.Trace("OpenWallet ok")

	wSvr.WtLoader.RunAfterLoad(func(w *wallet.Wallet) {
		w.Start()
	})

	wSvr.RegAPI()
	log.Trace("OpenWallet ok and reg api")

	return nil
}
