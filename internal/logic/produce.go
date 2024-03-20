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
		calls                    = make([]w3types.RPCCaller, 0)
		res                      = make([]*ContentFrom, 0)
		GlobalSlots, GlobalQueue int
	)
	l.OkAddrList = make([]common.Address, 0)
	for index, _addr := range l.svcCtx.AddrList {
		addr := _addr
		res = append(res, &ContentFrom{
			addr: addr,
		})
		calls = append(calls, txpool.ContentFrom(addr).Returns(&res[index].ContentFromResponse))
	}

	l.svcCtx.W3Cli.Call(calls...)

	for _, item := range res {
		GlobalSlots += len(item.Pending)
		GlobalQueue += len(item.Queued)

		if len(item.Pending) < 1 && GlobalSlots < 1280 && len(item.Queued) < 32 && GlobalQueue < 256 {
			l.OkAddrList = append(l.OkAddrList, item.addr)
		} else {
			// logx.Infof("addr: %v Pending) %v queue %v freega:%vGwei tipgas:%vGwei  ")
		}
	}

	return nil

}

func (l *Producer) BatchListNoceByAddr() error {

	var (
		calls = make([]w3types.RPCCaller, 0)
		res   = make([]*addrNonce, 0)
	)
	for index, _addr := range l.OkAddrList {
		addr := _addr
		res = append(res, &addrNonce{
			addr: addr,
		})
		calls = append(calls, eth.Nonce(addr, nil).Returns(&res[index].nonce))
	}
	l.svcCtx.W3Cli.Call(calls...)
	l.addrNonceList = res

	return nil

}

func (l *Producer) SendTxByAddrList() error {
	var (
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
		calls = append(calls, eth.SendTx(tx).Returns(&txHash))

	}

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

	for {
		// 1.查询余额
		err := l.CheckAddressBalance()
		if err != nil {
			logx.Errorf("%v", err)
			os.Exit(-1)
		}
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

		time.Sleep(1000 * time.Millisecond)
	}

}

func (l *Producer) Stop() {
	logx.Infof("stop Producer \n")
}
