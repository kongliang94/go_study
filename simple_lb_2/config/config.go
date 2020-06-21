package config

import (
	"fmt"
	"log"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

//RuntimeViper runtime config
var RuntimeViper *viper.Viper

func init() {
	RuntimeViper = viper.New()
	RuntimeViper.SetConfigType("toml")
	RuntimeViper.SetConfigName("cfg")                   // name of config file (without extension)
	RuntimeViper.AddConfigPath("/etc/proxy/simple_lb/") // path to look for the config file in
	RuntimeViper.AddConfigPath("./config/")             // optionally look for config in the working directory
	if err := RuntimeViper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}

	// 监听配置文件的改变，实现热部署
	RuntimeViper.WatchConfig()
	RuntimeViper.OnConfigChange(func(e fsnotify.Event) {
		log.Printf("config file changed:%s", e.Name)
	})
}
