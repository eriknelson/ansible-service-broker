package clients

import (
	logging "github.com/op/go-logging"
	restclient "k8s.io/client-go/rest"

	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"k8s.io/kubernetes/pkg/client/clientset_generated/clientset"
)

type kubernetesClientResult struct {
	client *clientset.Clientset
	err    error
}

func Kubernetes(log *logging.Logger) (*clientset.Clientset, error) {
	once.Kubernetes.Do(func() {
		client, err := newKubernetes(log)
		if err != nil {
			log.Error("An error occurred while initializing Kubernetes client:")
			log.Error(err.Error())
			instances.Kubernetes = kubernetesClientResult{nil, err}
		}
		instances.Kubernetes = kubernetesClientResult{client, nil}
	})

	err := instances.Etcd.err
	if err != nil {
		log.Error("Something went wrong initializing Kubernetes!")
		log.Error(err.Error())
		return nil, err
	}

	return instances.Kubernetes.client, nil
}

func createClientConfigFromFile(configPath string) (*restclient.Config, error) {
	clientConfig, err := clientcmd.LoadFromFile(configPath)
	if err != nil {
		return nil, err
	}

	config, err := clientcmd.NewDefaultClientConfig(*clientConfig, &clientcmd.ConfigOverrides{}).ClientConfig()
	if err != nil {
		return nil, err
	}
	return config, nil
}

func newKubernetes(log *logging.Logger) (*clientset.Clientset, error) {
	// NOTE: Both the external and internal client object are using the same
	// clientset library. Internal clientset normally uses a different
	// library
	clientConfig, err := restclient.InClusterConfig()
	if err != nil {
		log.Warning("Failed to create a InternalClientSet: %v.", err)

		log.Debug("Checking for a local Cluster Config")
		clientConfig, err = createClientConfigFromFile(homedir.HomeDir() + "/.kube/config")
		if err != nil {
			log.Error("Failed to create LocalClientSet")
			return nil, err
		}
	}

	clientset, err := clientset.NewForConfig(clientConfig)
	if err != nil {
		log.Error("Failed to create LocalClientSet")
		return nil, err
	}

	return clientset, err
}
