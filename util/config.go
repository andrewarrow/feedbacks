package util

import "github.com/spf13/viper"
import "fmt"
import "strings"

var AllConfig Config
var HostToPort = map[string]string{}
var Hosts = []string{}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
}

type PathConfig struct {
	Prefix string
}
type HttpConfig struct {
	Ports string
	Hosts string
}

type Config struct {
	Db   DatabaseConfig `mapstructure:"database"`
	Path PathConfig     `mapstructure:"paths"`
	Http HttpConfig     `mapstructure:"http"`
}

func InitConfig() bool {

	v := viper.New()
	v.AddConfigPath("/")
	v.AddConfigPath(".")
	v.AddConfigPath("..")
	v.SetConfigName("conf")
	v.SetConfigType("toml")
	v.AutomaticEnv()
	if err := v.ReadInConfig(); err != nil {
		fmt.Println(err)
		return false
	}
	if err := v.Unmarshal(&AllConfig); err != nil {
		fmt.Printf("couldn't read config: %s", err)
	}

	ports := strings.Split(AllConfig.Http.Ports, ",")
	Hosts = strings.Split(AllConfig.Http.Hosts, ",")

	for i, port := range ports {
		HostToPort[Hosts[i]] = port
	}

	return true
}
