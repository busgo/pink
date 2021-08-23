package main

import (
	"flag"
	client "github.com/busgo/pink-go"
	"github.com/busgo/pink-go/etcd"
	"github.com/busgo/pink-user/conf"
	"github.com/busgo/pink-user/job"
	"go.uber.org/dig"
	"time"
)

var confFile = flag.String("conf", "./trade-docker.yaml", "conf file")

func buildApp() *dig.Container {

	c := dig.New()
	c.Provide(func() (*conf.App, error) {
		return conf.NewAppConf(*confFile)
	})

	c.Provide(func(app *conf.App) (*etcd.Cli, error) {
		return etcd.NewEtcdCli(&etcd.CliConfig{
			Endpoints:   app.Etcd.Endpoints,
			UserName:    app.Etcd.UserName,
			Password:    app.Etcd.Password,
			DialTimeout: time.Second * time.Duration(app.Etcd.DialTimeout),
		})
	})

	c.Provide(func(app *conf.App, cli *etcd.Cli) *client.PinkClient {
		return client.NewPinkClient(cli, app.Group)
	})
	c.Provide(job.NewManaged)
	return c
}

func main() {

	flag.Parse()

	app := buildApp()

	err := app.Invoke(func(managed *job.Managed) error {
		managed.AddJob(new(job.UserLogJob))
		managed.AddJob(new(job.UserSignInJob))

		return managed.Run()

	})

	if err != nil {
		panic(err)
	}
}
