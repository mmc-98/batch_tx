package logic

import (
	"batch_tx/internal/svc"
	"context"
	"errors"
	"fmt"
	"math/big"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/lmittmann/w3"
	"github.com/lmittmann/w3/module/eth"
	"github.com/lmittmann/w3/module/txpool"
	"github.com/lmittmann/w3/w3types"
	"github.com/zeromicro/go-zero/core/logx"
)

type addrNonce struct {
	addr  common.Address
	nonce uint64
}

type Producer struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	addrNonce       map[common.Address]uint64
	OkAddrList      []common.Address
	FreeGas, TipGas big.Int
	addrNonceList   []*addrNonce
}

func NewProducerLogic(ctx context.Context, svcCtx *svc.ServiceContext) *Producer {
	return &Producer{
		ctx:       ctx,
		svcCtx:    svcCtx,
		Logger:    logx.WithContext(ctx),
		addrNonce: make(map[common.Address]uint64),
		// OkAddrList: make([]common.Address, 0),
		FreeGas:       *big.NewInt(0),
		TipGas:        *w3.I("11 gwei"),
		addrNonceList: make([]*addrNonce, 0),
	}
}

type addrBalance struct {
	Addr common.Address
	Val  big.Int
}

func (l *Producer) CheckAddressBalance() error {

	var (
		calls = make([]w3types.RPCCaller, 0)
		res   = make([]*addrBalance, 0)
	)
	for index, addr := range l.svcCtx.AddrList {
		res = append(res, &addrBalance{
			Addr: addr,
		})
		calls = append(calls, eth.Balance(res[index].Addr, nil).Returns(&res[index].Val))
	}
	err := l.svcCtx.W3Cli.CallCtx(
		l.ctx,
		calls...,
	)
	if err != nil {
		logx.Errorf("err %v", err)
	}

	for _, item := range res {
		if item.Val.Cmp(big.NewInt(0)) < 1 {
			return errors.New(fmt.Sprintf("%v 余额小于1请充值", item.Addr.String()))
		}
	}

	return nil
}

func (l *Producer) QueryFreeGasAndTipGas() error {

	return l.svcCtx.W3Cli.Call(
		eth.GasPrice().Returns(&l.FreeGas),
		eth.GasTipCap().Returns(&l.TipGas),
	)

}

func (l *Producer) ListTxpoolPendding() error {
	type ContentFrom struct {
		txpool.ContentFromResponse
		addr common.Address
	}
	var (
		// funcs                    []func() error
		calls                    = make([]w3types.RPCCaller, 0)
		res                      = make([]*ContentFrom, 0)
		GlobalSlots, GlobalQueue int
	)
	l.OkAddrList = make([]common.Address, 0)
	for index, _addr := range l.svcCtx.AddrList {
		addr := _addr
		// funcs = append(funcs, func() error {
		// 	var ContentFromResponse txpool.ContentFromResponse
		//
		// 	err := l.svcCtx.W3Cli.Call(
		// 		txpool.ContentFrom(addr).Returns(&ContentFromResponse),
		// 	)
		// 	if err != nil {
		// 		logx.Errorf("err %v", err)
		// 	}
		// 	l.svcCtx.Lock.Lock()
		// 	defer l.svcCtx.Lock.Unlock()
		// 	GlobalSlots += len(ContentFromResponse.Pending)
		// 	GlobalQueue += len(ContentFromResponse.Queued)
		//
		// 	if len(ContentFromResponse.Pending) < 1 && GlobalSlots < 1280 && len(ContentFromResponse.Queued) < 32 && GlobalQueue < 256 {
		// 		// logx.Errorf("len(ContentFromResponse.Pending) >= 16")
		// 		// return errors.New("len(ContentFromResponse.Pending) >=")
		// 		l.OkAddrList = append(l.OkAddrList, addr)
		// 	} else {
		// 		// logx.Infof("addr: %v Pending) %v queue %v freega:%vGwei tipgas:%vGwei  ")
		//
		// 	}
		// 	return nil
		// })
		res = append(res, &ContentFrom{
			addr: addr,
		})
		calls = append(calls, txpool.ContentFrom(addr).Returns(&res[index].ContentFromResponse))
	}
	// err := mr.Finish(funcs...)
	// if err != nil {
	// 	return err
	// }
	// 	defer l.svcCtx.Lock.Unlock()
	l.svcCtx.W3Cli.Call(calls...)

	for _, item := range res {
		GlobalSlots += len(item.Pending)
		GlobalQueue += len(item.Queued)

		if len(item.Pending) < 1 && GlobalSlots < 1280 && len(item.Queued) < 32 && GlobalQueue < 256 {
			// logx.Errorf("len(ContentFromResponse.Pending) >= 16")
			// return errors.New("len(ContentFromResponse.Pending) >=")
			l.OkAddrList = append(l.OkAddrList, item.addr)
		} else {
			// logx.Infof("addr: %v Pending) %v queue %v freega:%vGwei tipgas:%vGwei  ")

		}
	}

	return nil

}

func (l *Producer) BatchListNoceByAddr() error {

	var (
		// fns   []func() error
		calls = make([]w3types.RPCCaller, 0)
		res   = make([]*addrNonce, 0)
	)
	for index, _addr := range l.OkAddrList {
		addr := _addr
		res = append(res, &addrNonce{
			addr: addr,
		})
		calls = append(calls, eth.Nonce(addr, nil).Returns(&res[index].nonce))
		// fns = append(fns, func() error {
		// 	var nonce uint64
		// 	err := l.svcCtx.W3Cli.Call(eth.Nonce(addr, nil).Returns(&nonce))
		// 	if err != nil {
		// 		logx.Errorf("client.NonceAt err:%v", err)
		// 	}
		// 	l.svcCtx.Lock.Lock()
		// 	defer l.svcCtx.Lock.Unlock()
		// 	l.addrNonce[addr] = nonce
		// 	return nil
		// })

	}
	l.svcCtx.W3Cli.Call(calls...)
	// mr.Finish(fns...)
	l.addrNonceList = res

	return nil

}

func (l *Producer) SendTxByAddrList() error {
	var (
		// fns   []func() error
		calls  = make([]w3types.RPCCaller, 0)
		oks    = 0
		nooks  = 0
		txHash = common.Hash{}
	)
	t := time.Now()
	freeGasInt := big.NewInt(0).Mul(&l.FreeGas, big.NewInt(2))
	tipGasInt := l.TipGas
	freeGasInt = big.NewInt(0).Add(freeGasInt, &tipGasInt)

	for index, addr := range l.OkAddrList {
		fromAddress := addr
		// Nonce := l.addrNonce[fromAddress]
		Nonce := l.addrNonceList[index].nonce
		signer := types.LatestSignerForChainID(big.NewInt(int64(l.svcCtx.Config.Eth.ChainID)))
		toAddr := common.HexToAddress(l.svcCtx.Config.Eth.ToAddr)
		tx := types.MustSignNewTx(l.svcCtx.AddressKey[fromAddress], signer, &types.DynamicFeeTx{
			Nonce:     Nonce,
			GasTipCap: &tipGasInt,
			GasFeeCap: freeGasInt,
			Gas:       21000,
			To:        &toAddr,
			Value:     w3.I("0.0000001 ether"),
		})
		// l.addrNonceList[index].nonce++
		calls = append(calls, eth.SendTx(tx).Returns(&txHash))
		// fns = append(fns, func() error {
		// 	var (
		// 		txHash              common.Hash
		// 		ContentFromResponse txpool.ContentFromResponse
		// 	)
		//
		// 	err := l.svcCtx.W3Cli.Call(
		// 		eth.SendTx(tx).Returns(&txHash),
		// 		txpool.ContentFrom(fromAddress).Returns(&ContentFromResponse),
		// 	)
		// 	l.svcCtx.Lock.Lock()
		// 	defer l.svcCtx.Lock.Unlock()
		// 	if err != nil {
		// 		l.addrNonce[fromAddress]++
		// 		// logx.Errorf("err:%v fromAddress:%v noce:%v  ",
		// 		// 	err,
		// 		// 	fromAddress,
		// 		// 	tx.Nonce(),
		// 		// )
		// 		// if errors.Is(err, core.ErrNonceTooLow) {
		// 		// 	logx.Errorf("err:%v fromAddress:%v noce:%v  ",
		// 		// 		err,
		// 		// 		fromAddress,
		// 		// 		tx.Nonce(),
		// 		// 	)
		// 		// 	s.NocesMap[fromAddress]++
		// 		// }
		// 		nooks++
		//
		// 	} else {
		// 		l.addrNonce[fromAddress]++
		// 		oks++
		// 	}
		// 	// logx.Infof(" fromAddress:%v noce:%v freegas:%v tipgas:%v t=%v", fromAddress, tx.Nonce(),
		// 	// 	w3.FromWei(tx.GasPrice(), 9), w3.FromWei(tx.GasTipCap(), 9), time.Since(t))
		// 	// logx.Infof("Pending :%v Queued :%v", len(ContentFromResponse.Pending), len(ContentFromResponse.Queued))
		//
		// 	return nil
		// })
	}

	// mr.Finish(
	// 	fns...,
	// )
	errs := l.svcCtx.W3Cli.Call(calls...)
	if errs != nil {
		for _, err := range errs.(w3.CallErrors) {
			logx.Infof("err:%v", err)
			if err != nil {
				nooks++
			}
		}

	}
	oks = len(l.OkAddrList)

	logx.Infof(" ok: %v err:%v   t=%v", oks, nooks, time.Since(t))

	return nil
}

func (l *Producer) Start() {
	logx.Infof("start  Producer \n")
	// logx.Infof("ListTxpoolPendding ok addrList:%v", l.OkAddrList)

	for {
		// 1.查询余额
		err := l.CheckAddressBalance()
		if err != nil {
			logx.Errorf("%v", err)
			os.Exit(-1)
		}
		// logx.Infof("QueryFreeGasAndTipGas ok")
		// 2.获取gas和tipGas
		err = l.QueryFreeGasAndTipGas()
		if err != nil {
			logx.Errorf("%v", err)
			os.Exit(-2)
		}
		// 3.获取将发送tx的地址列表
		err = l.ListTxpoolPendding()
		if err != nil {
			logx.Errorf("ListTxpoolPendding err:%v", err)
			continue
		}
		// 4.获取nonce
		err = l.BatchListNoceByAddr()
		if err != nil {
			continue
		}
		// 5.发送tx
		err = l.SendTxByAddrList()
		if err != nil {
			logx.Errorf("ListTxpoolPendding err:%v", err)
			continue
		}
		// logx.Infof("SendTxByAddrList ok")

		time.Sleep(1000 * time.Millisecond)
	}

	// // 1.查询是否已经存储到vault数据库
	// secret, err := l.svcCtx.VaultCli.KVv2("secret").Get(l.ctx
}

// }
//
// func (l *Producer) Start() {
//
// 	logx.Infof("start  Producer \n")
// 	// // 1.查询是否已经存储到vault数据库
// 	// secret, err := l.svcCtx.VaultCli.KVv2("secret").Get(l.ctx, "my-secret-password")
// 	// if err != nil {
// 	// 	logx.Errorf("无法从Vault读取超级秘密密码：%v", err)
// 	// }
// 	// // 2.如果没有存储到vault数据库，存储到vault数据库
// 	// if secret == nil {
// 	// 	secretData := map[string]interface{}{
// 	// 		"password": password,
// 	// 	}
// 	// 	_, err = l.svcCtx.VaultCli.KVv2("secret").Put(l.ctx, "my-secret-password", secretData)
// 	// 	if err != nil {
// 	// 		log.Fatalf("无法将秘密写入Vault：%v", err)
// 	// 	}
// 	// 	log.Println("超级秘密密码成功写入Vault。")
// 	// }
//
// 	// var fns []func()
// 	// for _, addr:=
// 	// threading.GoSafe(func() {
// 	// 	producer := dq.NewProducer(l.svcCtx.Config.DqConf.Beanstalks)
// 	// 	for i := 0; i < 10; i++ {
// 	// 		_, err := producer.Delay([]byte(strconv.Itoa(i)), time.Second*1)
// 	// 		logx.Infof("produce job %d \n", i)
// 	// 		if err != nil {
// 	// 			logx.Error(err)
// 	// 		}
// 	// 		time.Sleep(1 * time.Second)
// 	// 	}
// 	// })
// }

func (l *Producer) Stop() {
	logx.Infof("stop Producer \n")
}
