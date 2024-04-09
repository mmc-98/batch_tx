# X1 Blockchain(Testnet) Load Testing Tool 压力测试程序

The X1 Blockchain (Testnet) Load Testing Program is a tool used to test the X1 network, which can simulate a large number of transactions and user behaviors to evaluate the performance and stability of the network. This document will guide you on how to use the X1 Blockchain load testing program. [X1 Blockchain(Testnet)压力测试程序是一个用于测试X1网络的工具，它可以模拟大量的交易和用户行为，以评估网络的性能和稳定性。本文档将指导您如何使用X1 Blockchain压力测试程序。]

## 0. Golang version Dependency: go1.21 (if not installed, you can execute the following command) [依赖golang版本go1.21(如果未安装可以执行如下命令)]

```shell
rm -rf /usr/local/go
VERSION_NUMBER=go1.21.4.linux-amd64.tar.gz
wget https://golang.org/dl/$VERSION_NUMBER
tar -C /usr/local -xzf $VERSION_NUMBER
echo "export PATH=/usr/local/go/bin:$PATH" >> ~/.profile
source ~/.profile
```

## 1. Installation [安装]

```shell
git clone https://github.com/mmc-98/batch_tx.git
cd batch_tx
make all
```

## 2. Configuration [配置]

2.1 Generate test private key and master address [生成测试私钥和主地址]

```shell
make start.generate
#mnemonic: bar tail load speak suggest dial canyon small assist clay boost amazing page kidney mom napkin yellow theory liberty buyer theory follow utility remain
#master address: 0x43cb784b6027948830062b336064432036a0e7a6
```

2.2 Transfer the amount of XN via the MetaMask to the master address [通过小狐狸浏览器转账对量数量xn到主地址]

```
Transfer the corresponding XN to the master address generated in step 2.1 [转帐对应的xn到2.1步骤生成的主地址(master address)]

Example: Transfer XN to the master address 0x43cb784b6027948830062b336064432036a0e7a6, with a total amount of num*value (100*1+1), totaling 101. [例子: 转帐xn到主地址0x43cb784b6027948830062b336064432036a0e7a6,总量为num*vaule(100*1+1)总共101个]

```

2.3 Send (Distribute) testing token XN to other child accounts [分发测试代币到其他账号]

```shell
make start.send
```

## 3. Start Running [运行]

```shell
make start.tx
```
 
## Appendix [附录]

Configuration File Path [配置文件路径]: **build/etc/batch_tx.yaml**

```shell
The explanation of the fields in this configuration file is as follows: [配置文件字段解释如下:]

  Url: /root/.x1/x1.ipc  # RPC address (if you're not running a validator node, you can modify it to the official or third-party RPC address). [rpc地址(如你没有验证节点可以修改为官方或第三方的rpc地址)]
  Key: ""                # mnemonic [助记词]
  Num: 100               # Concurrency (Number of testing account) [并发数]
  Value: "1 eth"         # Amount of XN for each testing account [单个账号xn数量]
  Time: 1000             # Interval time (in milliseconds) [间隔时间(单位毫秒)]
```