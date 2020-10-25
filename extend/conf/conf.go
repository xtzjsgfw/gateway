package conf

import (
	"bytes"
	"fmt"
	"github.com/spf13/viper"
	"io/ioutil"
)

type server struct {
	RunMode   string `mapstructure:"runMode"`
	Port      string `mapstructure:"port"`
	JwtSecret string `mapstructure:"jwtSecret"`
	JwtExpire int    `mapstructure:"jwtExpire"`
}

var ServerConf = &server{}

type database struct {
	DBType       string `mapstructure:"dbType"`
	DBName       string `mapstructure:"dbName"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	Host         string `mapstructure:"host"`
	Port         string `mapstructure:"port"`
	Debug        bool   `mapstructure:"debug"`
	MaxIdleConns int    `mapstructure:"maxIdleConns"`
	MaxOpenConns int    `mapstructure:"maxOpenConns"`
}

var DBConf = &database{}

type log struct {
	LogPath string `mapstructure:"logPath"`
}

var LogConf = &log{}

type redis struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
}

var RedisConf = &redis{}

type cluster struct {
	IP      string `mapstructure:"ip"`
	Port    string `mapstructure:"port"`
	SSLPort string `mapstructure:"sslPort"`
}

var ClusterConf = &cluster{}

type http struct {
	Addr           string `mapstructure:"addr"`
	ReadTimeout    int    `mapstructure:"readTimeout"`
	WriteTimeout   int    `mapstructure:"writeTimeout"`
	MaxHeaderBytes int    `mapstructure:"maxHeaderBytes"`
}

var HttpConf = &http{}

type https struct {
	Addr           string `mapstructure:"addr"`
	ReadTimeout    int    `mapstructure:"readTimeout"`
	WriteTimeout   int    `mapstructure:"writeTimeout"`
	MaxHeaderBytes int    `mapstructure:"maxHeaderBytes"`
}

var HttpsConf = &https{}

func Init() {
	viper.SetConfigType("yaml")
	confFile, err := ioutil.ReadFile("config/config.yaml")
	if err != nil {
		fmt.Println("读取配置文件失败")
	}
	err = viper.ReadConfig(bytes.NewBuffer(confFile))
	if err != nil {
		fmt.Println("Viper库读取配置文件失败")
	}
	viper.UnmarshalKey("server", ServerConf)
	viper.UnmarshalKey("database", DBConf)
	viper.UnmarshalKey("log", LogConf)
	viper.UnmarshalKey("redis", RedisConf)
	viper.UnmarshalKey("cluster", ClusterConf)
	viper.UnmarshalKey("http", HttpConf)
	viper.UnmarshalKey("https", HttpsConf)
}
