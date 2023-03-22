package config

import (
	"encoding/json"
	"log"
	"os"
	"strings"
	"sync"
)

// Configuration 项目配置
type Configuration struct {
	// gtp apikey
	ApiKey string `json:"api_key"`
	// 自动通过好友
	AutoPass bool     `json:"auto_pass"`
	Backend  string   `json:"backend"`
	Model    string   `json:"model"`
	KeyWords string   `json:"key_words"`
	Keys     []string `json:"-"`
}

var config *Configuration
var once sync.Once

// LoadConfig 加载配置
func LoadConfig() *Configuration {
	once.Do(func() {
		// 从文件中读取
		config = &Configuration{}
		f, err := os.Open("config.json")
		if err != nil {
			log.Fatalf("open config err: %v", err)
			return
		}
		defer f.Close()
		encoder := json.NewDecoder(f)
		err = encoder.Decode(config)
		if err != nil {
			log.Fatalf("decode config err: %v", err)
			return
		}

		// 如果环境变量有配置，读取环境变量
		ApiKey := os.Getenv("ApiKey")
		AutoPass := os.Getenv("AutoPass")
		if ApiKey != "" {
			config.ApiKey = ApiKey
		}
		if AutoPass == "true" {
			config.AutoPass = true
		}
		if config.KeyWords == "" {
			config.KeyWords = "投资,理财,股票,基金,保险,信托,债券,期货,外汇,黄金,银行,银行卡,信用卡,贷款,房贷,车贷,留学,移民,税,钱,支付宝,银行,款,价格,资金"
		}
		config.Keys = strings.Split(config.KeyWords, ",")
	})
	return config
}
