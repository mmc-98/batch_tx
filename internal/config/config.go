package config

import (
	"github.com/zeromicro/go-zero/core/service"
)

type Config struct {
	service.ServiceConf
	// Redis redis.RedisConf
	//
	// DB struct {
	// 	DataSource string
	// }
	// Cache cache.CacheConf
	//
	// // KqPusherConf struct {
	// // 	Brokers []string
	// // 	Topic   string
	// // }
	// DqConf dq.DqConf
	Eth struct {
		Url     string
		Key     string
		Num     int
		ChainID int
		ToAddr  string
		Value   string
		Time    int
	}
	// Vault struct {
	// 	Address *vault.Config
	// 	Token   string
	// }
}
