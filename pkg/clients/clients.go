package clients

import (
	"sync"
)

var instances struct {
	Etcd       etcdClientResult
	Kubernetes kubernetesClientResult
	Docker     dockerClientResult
}

var once struct {
	Etcd       sync.Once
	Kubernetes sync.Once
	Docker     sync.Once
}
