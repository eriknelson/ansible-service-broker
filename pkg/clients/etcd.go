package clients

import (
	"fmt"
	"time"

	logging "github.com/op/go-logging"

	etcd "github.com/coreos/etcd/client"
)

type EtcdConfig struct {
	EtcdHost string `yaml:"etcd_host"`
	EtcdPort string `yaml:"etcd_port"`
}

func Etcd(config EtcdConfig, log *logging.Logger) (*etcd.Client, error) {
	once.Etcd.Do(func() {
		client, err := newEtcd(config, log)
		if err != nil {
			log.Error("An error occurred while initializing Etcd client:")
			log.Error(err.Error())
			instances.Etcd = clientResult{nil, err}
		}
		instances.Etcd = clientResult{client, nil}
	})

	err := instances.Etcd.err
	if err != nil {
		log.Error("Something went wrong initializing Etcd!")
		log.Error(err.Error())
		return nil, err
	}

	return instances.Etcd.client, nil
}

func newEtcd(config EtcdConfig, log *logging.Logger) (*etcd.Client, error) {
	// TODO: Config validation
	endpoints := []string{etcdEndpoint(config.EtcdHost, config.EtcdPort)}

	log.Info("== ETCD CX ==")
	log.Infof("EtcdHost: %s", config.EtcdHost)
	log.Infof("EtcdPort: %s", config.EtcdPort)
	log.Infof("Endpoints: %v", endpoints)

	etcdClient, err := etcd.New(etcd.Config{
		Endpoints:               endpoints,
		Transport:               etcd.DefaultTransport,
		HeaderTimeoutPerRequest: time.Second,
	})
	if err != nil {
		return nil, err
	}

	return &etcdClient, err
}

func etcdEndpoint(host string, port string) string {
	return fmt.Sprintf("http://%s:%s", host, port)
}
