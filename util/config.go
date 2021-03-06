package util

import "github.com/spf13/viper"
import "fmt"

var AllConfig Config

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
}

type PathConfig struct {
	Prefix string
	Sites  string
}
type HttpConfig struct {
	Secret string
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

	return true
}
