package svc

import (
	"batch_tx/internal/config"
	"crypto/ecdsa"
	"fmt"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/lmittmann/w3"
	"github.com/okx/go-wallet-sdk/example"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/mr"
)

type ServiceContext struct {
	Config config.Config
	// Consumer   dq.Consumer
	W3Cli      *w3.Client
	Lock       sync.RWMutex
	AddressKey map[common.Address]*ecdsa.PrivateKey
	AddrList   []common.Address
	// VaultCli *vault.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	// VaultCli, err := vault.NewClient(c.Vault.Address)
	// if err != nil {
	// 	panic(err)
	// }
	// VaultCli.SetToken(c.Vault.Token)
	s := &ServiceContext{
		Config: c,
		// Consumer:   dq.NewConsumer(c.DqConf),
		W3Cli:      w3.MustDial(c.Eth.Url),
		Lock:       sync.RWMutex{},
		AddressKey: make(map[common.Address]*ecdsa.PrivateKey),
		AddrList:   make([]common.Address, 0),
		// VaultCli: VaultCli,
	}
	s.SetAddrKeyAndAddrList()
	return s
}

func (s *ServiceContext) SetAddrKeyAndAddrList() {
	var (
		funcs []func() error
	)

	for i := 0; i < s.Config.Eth.Num; i++ {
		hdPath := example.GetDerivedPath(i)
		derivePrivateKey, _ := example.GetDerivedPrivateKey(s.Config.Eth.Key, hdPath)
		address := example.GetNewAddress(derivePrivateKey)

		privateKeyByte, err := hexutil.Decode(fmt.Sprintf("0x%v", derivePrivateKey))
		if err != nil {
			logx.Errorf("err:%v", err)
		}

		privateKey, err := crypto.ToECDSA(privateKeyByte)

		funcs = append(funcs, func() error {
			s.Lock.Lock()
			defer s.Lock.Unlock()
			s.AddressKey[common.HexToAddress(address)] = privateKey
			s.AddrList = append(s.AddrList, common.HexToAddress(address))
			return nil
		})
	}
	err := mr.Finish(funcs...)
	if err != nil {
		logx.Errorf("client.Call err:%v", err)
	}

}
