package hippo

import (
	"encoding/json"
	"github.com/jhaals/yopass/pkg/server"
	"github.com/jhaals/yopass/pkg/yopass"

	"hippos/walrus"
)

var hippoSecret yopass.Secret

type Hippo struct {
	redis  server.Database
	walrus *walrus.WalrusClient
}

func NewHippo(walrusURL string, redisURL string) (server.Database, error) {
	redis, err := server.NewRedis(redisURL)
	if err != nil {
		return nil, err
	}
	return &Hippo{redis, walrus.NewWalrusClient(walrusURL)}, err
}

func (d *Hippo) Get(key string) (yopass.Secret, error) {
	metaSecret, err := d.redis.Get(key)
	if err != nil {
		return metaSecret, err
	}
	v, err := d.walrus.GetData(metaSecret.Message)
	if err != nil {
		return metaSecret, err
	}

	var s yopass.Secret
	if err := json.Unmarshal([]byte(v), &s); err != nil {
		return s, err
	}

	return s, nil
}

func (d *Hippo) Put(key string, secret yopass.Secret) error {
	data, err := secret.ToJSON()
	if err != nil {
		return err
	}

	blobID, err := d.walrus.PutData(data)
	if err != nil {
		return err
	}

	secret.Message = blobID
	return d.redis.Put(key, secret)
}

func (d *Hippo) Delete(key string) (bool, error) {
	return d.redis.Delete(key)
}
