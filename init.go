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
	if c == nil && rawESClient == nil {
		return fmt.Errorf("elasticsearch client is nil")
	}
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

func NewClient(c *elasticsearch.Client) (GES, error) {
	if c == nil && rawESClient == nil {
		return nil, fmt.Errorf("elasticsearch client is nil")
	}
	if c == nil {
		c = rawESClient
	}

	return &ges{c: c}, nil
}

func NewClientCfg(addrs []string, user, pwd string) (GES, error) {
	cfg := elasticsearch.Config{
		Addresses: addrs,
		Username:  user,
		Password:  pwd,
	}
	c, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("elastic init err: %v", err)

	}

	return NewClient(c)
}

type GES interface {
	Client() Client
}

type ges struct {
	c *elasticsearch.Client
}

func (g *ges) Client() Client {
	return newES(g.c)
}

func newES(c *elasticsearch.Client) Client {
	return &es{
		isAgg:     false,
		sorts:     []string{},
		fields:    nil,
		indexName: "",
		from:      0,
		size:      0,
		cond:      cond{},
		agg:       make(map[string]interface{}, 0),
		c:         c,
	}
}
