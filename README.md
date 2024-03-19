# X1 Blockchain(Testnet) 压力测试程序
X1 Blockchain(Testnet)压力测试程序是一个用于测试X1网络的工具，它可以模拟大量的交易和用户行为，以评估网络的性能和稳定性。本文档将指导您如何使用X1 Blockchain压力测试程序。


## 1.安装
```shell
go build .
```
## 2.配置
配置文件在etc/batch_tx.yaml中，您可以根据需要进行修改。
Key为必填字段
```shell
Key: 字段填写测试私钥
num: 字段填写并发数
```


## 3.运行
```shell
./batch_tx
```
 