package config

import (
	"os"
	log "github.com/cihub/seelog"
	homedir "github.com/mitchellh/go-homedir"
	"path/filepath"
	"encoding/json"
)
/**
  Proxy config is used for storing a url to ECR registry mapping.
  This is useful when setting up a reverse proxy for ECR to cache
  container data.
  {"proxies":{"ecrproxy.myasws.com":"12345678910.dkr.ecr.us-east-1.amazonaws.com"}}
 */
type ProxyConfig struct {
	Proxies map[string]string `json:"proxies"`
}

func GetProxyConfig() (*ProxyConfig, error) {
	config := ProxyConfig{}
	dir, _ := homedir.Expand(GetCacheDir())

	//Save it the same as the cache dir now
	proxyFile := filepath.Join(dir, "ecr_proxy.json")
	if _, err := os.Stat(proxyFile); err != nil {
		log.Info("No Proxy config found")
		return nil, err
	}

	reader, err := os.Open(proxyFile)
	if (err != nil) {
		return nil, err
	}

	defer reader.Close()
	if err := json.NewDecoder(reader).Decode(&config); err != nil {
		log.Errorf("Fail to load config with error %s", err)
		return nil, err
	}

	return &config, nil
}

func GetRegistryURL(serverUrl string) string {
	proxyConfig, err := GetProxyConfig()
	if err != nil {
		return serverUrl
	}
	if val, ok := proxyConfig.Proxies[serverUrl]; ok{
		return val
	}
	//Return the input if no mapping found
	return serverUrl
}