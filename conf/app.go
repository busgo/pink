package conf

import (
	"github.com/spf13/viper"
	"log"
	"os"
)

type AppConf struct {
	Etcd struct {
		Endpoints   []string `yaml:"endpoints"`
		UserName    string   `yaml:"user_name"`
		Password    string   `yaml:"password"`
		DialTimeout int64    `yaml:"dial_timeout"`
	}

	Port int32 `yaml:"port"`
}

func NewAppConf(filePath string) (*AppConf, error) {

	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("can not open file:%s", filePath)
		return nil, err
	}
	v := viper.New()
	v.SetConfigType("yaml")
	err = v.ReadConfig(file)
	if err != nil {
		return nil, err
	}
	c := new(AppConf)
	err = v.Unmarshal(c)
	return c, err
}
