package config

import (
	"encoding/json"
	"log"
	"os"
	"sync"
	"time"
)

// Configuration 项目配置
type Configuration struct {
	// gtp apikey
	ApiKey string `json:"api_key"`
	// 自动通过好友
	AutoPass bool `json:"auto_pass"`
	// 是否使用代理
	UseProxy bool `json:"use_proxy"`
	// 代理地址
	ProxyUrl string `json:"proxy_url"`
	// 是否使用gptTurbo
	UseGPT bool `json:"use_gpt"`
	// 会话超时时间
	SessionTimeout time.Duration `json:"session_timeout"`
}

var config *Configuration
var once sync.Once

// LoadConfig 加载配置
func LoadConfig() *Configuration {
	once.Do(func() {
		// 从文件中读取
		config = &Configuration{
			SessionTimeout: 1,
		}
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
		SessionTimeout := os.Getenv("SessionTimeout")
		UseProxy := os.Getenv("UseProxy")
		ProxyUrl := os.Getenv("ProxyUrl")
		UseGPT := os.Getenv("UseGPT")
		if ProxyUrl != "" {
			config.ProxyUrl = ProxyUrl
		}
		if UseProxy != "" {
			if UseProxy == "true" {
				config.UseProxy = true
			} else {
				config.UseProxy = false
			}
		}
		if UseGPT != "" {
			if UseGPT == "true" {
				config.UseGPT = true
			} else {
				config.UseGPT = false
			}
		}
		if ApiKey != "" {
			config.ApiKey = ApiKey
		}
		if AutoPass == "true" {
			config.AutoPass = true
		}
		if SessionTimeout != "" {
			duration, err := time.ParseDuration(SessionTimeout)
			if err != nil {
				log.Fatalf("config decode session timeout err: %v ,get is %v", err, SessionTimeout)
				return
			}
			config.SessionTimeout = duration
		}
	})
	return config
}
