# X1 Blockchain(Testnet) 压力测试程序
X1 Blockchain(Testnet)压力测试程序是一个用于测试X1网络的工具，它可以模拟大量的交易和用户行为，以评估网络的性能和稳定性。本文档将指导您如何使用X1 Blockchain压力测试程序。


## 1.安装
```shell
git clone https://github.com/mmc-98/batch_tx.git
cd batch_tx
make all
```
## 2.配置
2.1 生成测试私钥和主地址
```shell
make start.generate
#mnemonic: bar tail load speak suggest dial canyon small assist clay boost amazing page kidney mom napkin yellow theory liberty buyer theory follow utility remain
#master address: 0x43cb784b6027948830062b336064432036a0e7a6
```

2.2 修改测试私钥
``` 
把2.1步骤生成的私钥(mnemonic)添加到etc/batch_tx.yaml中的key字段
其中字段解释:
  Url: /root/.x1/x1.ipc  # rpc地址(如你没有验证节点可以修改为官方的rpc地址:https://x1-devnet.xen.network)
  Key: ""                # 助记词
  Num: 100               # 并发数
  Value: "1 eth"         # 单个账号xn数量
```
2.3 通过小狐狸浏览器转账对量数量xn到主地址
```
转帐对应的xn到2.1步骤生成的主地址(master address)
列子: 转帐xn到主地址0x43cb784b6027948830062b336064432036a0e7a6,总量为num*vaule(100*1+1)总共101个

```
2.4 分发测试代币到其他账号
```shell
make start.send
```
 

## 3.运行
```shell
start.tx
```
 