package conf

import (
	"github.com/spf13/viper"
	"log"
	"os"
)

type App struct {
	Etcd struct {
		Endpoints   []string `yaml:"endpoints"`
		UserName    string   `yaml:"user_name"`
		Password    string   `yaml:"password"`
		DialTimeout int64    `yaml:"dial_timeout"`
	}

	Group string `yaml:"group"`
}

func NewAppConf(filePath string) (*App, error) {

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
	c := new(App)
	err = v.Unmarshal(c)
	return c, err
}
