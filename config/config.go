package config

import (
	"encoding/json"
	"flag"
	"io/ioutil"
)

type EnvConfig struct {
	Log struct {
		Access      string `json:"access"`
		Application string `json:"application"`
		Compress    bool   `json:"compress"`
		MaxSize     int    `json:"max_size"`
	} `json:"log"`
	Mysql struct {
		ConnnectTimeout int `json:"connnect_timeout"`
		Database        struct {
			Master struct {
				Dbname string `json:"dbname"`
				Dbpass string `json:"dbpass"`
				Dbuser string `json:"dbuser"`
				Host   string `json:"host"`
				Port   int    `json:"port"`
			} `json:"master"`
			Slave []struct {
				Dbname string `json:"dbname"`
				Dbpass string `json:"dbpass"`
				Dbuser string `json:"dbuser"`
				Host   string `json:"host"`
				Port   int    `json:"port"`
			} `json:"slave"`
		} `json:"database"`
		MaxIdleConnections int `json:"max_idle_connections"`
		MaxOpenConections  int `json:"max_open_conections"`
		ReadTimeout        int `json:"read_timeout"`
		WriteTimeout       int `json:"write_timeout"`
	} `json:"mysql"`
	ProductLine struct {
		Rongcrypt struct {
			Acl struct {
				Access      string `json:"access"`
				Secret      string `json:"secret"`
				VerifyToken int    `json:"verify_token"`
			} `json:"acl"`
			Machine []string `json:"machine"`
			Talk    struct {
				ConnectionTimeoutMs int `json:"connection_timeout_ms"`
				ReadTimeoutMs       int `json:"read_timeout_ms"`
				WriteTimeoutMs      int `json:"write_timeout_ms"`
			} `json:"talk"`
		} `json:"rongcrypt"`
	} `json:"product_line"`
	Redis struct {
		Machine []struct {
			Addr     string `json:"addr"`
			Password string `json:"password"`
		} `json:"machine"`
		MaxActive           int `json:"max_active"`
		MaxIdle             int `json:"max_idle"`
		ReadTimeoutSeconds  int `json:"read_timeout_seconds"`
		WriteTimeoutSeconds int `json:"write_timeout_seconds"`
	} `json:"redis"`
}

var GlobalEnv EnvConfig

func LoadEnv() {
	jsonFile := flag.String("c", "./config/env.json", "the path of the env file")
	flag.Parse()
	bytes, err := ioutil.ReadFile(*jsonFile)
	if err != nil {
		panic("read file error")
	}
	if err = json.Unmarshal(bytes, &GlobalEnv); err != nil {
		panic("json unmarshal error")
	}
}
