package ges

import (
	"fmt"
	"sync"

	"github.com/elastic/go-elasticsearch/v7"
)

/***************************
    @author: tiansheng.ren
    @date: 2022/12/16
    @desc:

***************************/

var (
	rawESClient *elasticsearch.Client = nil
	once                              = sync.Once{}
)

func InitDefaultClient(c *elasticsearch.Client) error {
	err := fmt.Errorf("duplicate init elasticsearch")
	once.Do(func() {
		rawESClient = c
		err = nil
	})

	return err
}

func InitClientWithCfg(addrs []string, user, pwd string) error {
	cfg := elasticsearch.Config{
		Addresses: addrs,
		Username:  user,
		Password:  pwd,
	}
	c, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return fmt.Errorf("elastic init err: %v", err)

	}

	rawESClient = c
	return nil
}
