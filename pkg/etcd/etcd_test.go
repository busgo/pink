package etcd

import (
	"context"
	"testing"
	"time"
)

func TestCli_GetWithPrefix(t *testing.T) {

	cli, err := NewEtcdCli(&CliConfig{
		Endpoints:   []string{"127.0.0.1:2379"},
		UserName:    "",
		Password:    "",
		DialTimeout: time.Second * 5,
	})

	if err != nil {
		panic(err)
	}

	cli.c.Txn(context.Background()).If()
}
