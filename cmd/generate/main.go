package main

import (
	"batch_tx/internal/config"
	"flag"
	"fmt"
	"os"
	"path"
	"text/template"

	"github.com/okx/go-wallet-sdk/example"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
)

const (
	regularPerm = 0o666
)

type ServiceContext struct {
	Config config.Config
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
	}
}

func (l *ServiceContext) generateConf(tpl string, dst string, data any) error {
	file, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, regularPerm)
	defer file.Close()
	if err != nil {
		return err
	}

	wd, err := os.Getwd()
	if err != nil {
		fmt.Println("错误:", err)
		return err
	}
	// fmt.Println("当前工作目录:", wd)
	tpll := path.Join(wd, tpl)
	t, err := template.ParseFiles(tpll)
	if err != nil {
		return err
	}

	return t.Execute(file, data)
}

func (l *ServiceContext) generateMnemonicAndMasterAddr() (string, string, error) {
	mnemonic, err := example.GenerateMnemonic()
	if err != nil {
		return "", "", err
	}
	// fmt.Printf("mnemonic: %v\n", mnemonic)
	hdPath := example.GetDerivedPath(0)
	derivePrivateKey, err := example.GetDerivedPrivateKey(mnemonic, hdPath)

	// get new address
	newAddress := example.GetNewAddress(derivePrivateKey)
	return mnemonic, newAddress, nil
	// fmt.Printf("master address: %v\n", newAddress)
}

var (
	configFile = flag.String("f", "build/etc/batch_tx.yaml", "the config file")
	tplFile    = flag.String("s", "cmd/generate/batch_tx.tpl", "the tpl config file")
)

type mnemonicData struct {
	Mnemonic string
}

func main() {
	flag.Parse()
	logx.DisableStat()
	// 配置
	var c config.Config
	conf.MustLoad(*configFile, &c)
	s := NewServiceContext(c)

	mnemonic, masterAddr, err := s.generateMnemonicAndMasterAddr()
	if err != nil {
		logx.Errorf("生成错误")
		return
	}

	err = s.generateConf(*tplFile, *configFile, &mnemonicData{
		Mnemonic: mnemonic,
	})
	if err != nil {
		logx.Errorf("生成模版错误")
		return
	}

	fmt.Printf("mnemonic: %v \n master address: %v \n", mnemonic, masterAddr)
}
