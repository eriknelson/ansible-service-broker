package clients

import (
	docker "github.com/fsouza/go-dockerclient"
	logging "github.com/op/go-logging"
)

type dockerClientResult struct {
	client *docker.Client
	err    error
}

func Docker(log *logging.Logger) (*docker.Client, error) {
	once.Docker.Do(func() {
		client, err := newDocker(log)
		if err != nil {
			log.Error("An error occurred while initializing Docker client:")
			log.Error(err.Error())
			instances.Docker = dockerClientResult{nil, err}
		}
		instances.Docker = dockerClientResult{client, nil}
	})

	err := instances.Docker.err
	if err != nil {
		log.Error("Something went wrong initializing Docker!")
		log.Error(err.Error())
		return nil, err
	}

	return instances.Docker.client, nil
}

func newDocker(log *logging.Logger) (*docker.Client, error) {
	return docker.NewClient(DockerSocket)
}
