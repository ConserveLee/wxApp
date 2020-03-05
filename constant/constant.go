package constant

import (
	"encoding/json"
	"log"
	"os"
)

type config struct {
	ConfigMap map[string]string
}

var ServerConfig config

// 初始化config文件
func init() {
	conf, err := os.Open("./conf/config.conf")
	defer conf.Close()
	if err != nil {
		log.Println(err.Error())
	}

	ServerConfig.ConfigMap = make(map[string]string)
	err = json.NewDecoder(conf).Decode(&ServerConfig.ConfigMap)
	if err != nil {
		log.Println(err.Error())
	}
}

// 获取config的值
func (config *config) Get(key string) string {
	if val, ok := config.ConfigMap[key]; ok {
		return val
	}
	log.Printf("配置文件没有该键: %s!\n", key)
	return ""
}