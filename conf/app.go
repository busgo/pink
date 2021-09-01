package conf

import (
	"github.com/busgo/pink/pkg/log"
	"github.com/spf13/viper"
	"os"
)

type AppConf struct {
	Etcd struct {
		Endpoints   []string
		UserName    string
		Password    string
		DialTimeout int64
	}
	Mysql struct {
		DSN string
	}
	Log struct {
		FileName    string
		Level       int32
		ServiceName string
	}
	Port int32
}

func NewAppConf(filePath string) (*AppConf, error) {

	file, err := os.Open(filePath)
	if err != nil {
		log.Errorf("can not open file:%s", filePath)
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
	if err == nil {
		_ = log.NewLoggerSugar(c.Log.ServiceName, c.Log.FileName, c.Log.Level)
	}
	return c, err
}
