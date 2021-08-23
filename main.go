package main

import (
	"flag"
	"github.com/busgo/pink/collector"
	"github.com/busgo/pink/conf"
	"github.com/busgo/pink/executor"
	"github.com/busgo/pink/http"
	"github.com/busgo/pink/http/handler"
	"github.com/busgo/pink/node"
	"github.com/busgo/pink/pkg/bus"
	"github.com/busgo/pink/pkg/etcd"
	"github.com/busgo/pink/schedule"
	"go.uber.org/dig"
	"log"
	"time"
)

var confFile = flag.String("conf", "./app.yaml", "conf file")

func BuildApp() *dig.Container {

	app := dig.New()
	_ = app.Provide(func() (*conf.AppConf, error) {
		return conf.NewAppConf(*confFile)
	})
	_ = app.Provide(func(app *conf.AppConf) (*etcd.Cli, error) {
		return etcd.NewEtcdCli(&etcd.CliConfig{
			Endpoints:   app.Etcd.Endpoints,
			UserName:    app.Etcd.UserName,
			Password:    app.Etcd.Password,
			DialTimeout: time.Second * time.Duration(app.Etcd.DialTimeout),
		})
	})
	_ = app.Provide(bus.NewEventBus)
	_ = app.Provide(collector.NewPinkCollector)
	_ = app.Provide(executor.NewPinkGroupManaged)
	_ = app.Provide(executor.NewPinkExecutor)
	_ = app.Provide(schedule.NewPinkScheduler)
	_ = app.Provide(node.NewPinkNode)
	_ = app.Provide(handler.NewPinkWebHandler)
	_ = app.Provide(http.NewHttpService)
	return app
}

func main() {

	flag.Parse()
	app := BuildApp()
	err := app.Invoke(func(s *http.Service) error {
		log.Printf("http service start run......")
		return s.Run()
	})
	if err != nil {
		panic(err)
	}

}
