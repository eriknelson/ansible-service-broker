//
// Copyright (c) 2017 Red Hat, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package main

import (
	logging "github.com/op/go-logging"
	"runtime/debug"
	//"github.com/openshift/ansible-service-broker/pkg/app"
	"fmt"
	"github.com/openshift/ansible-service-broker/pkg/clients"
	//"k8s.io/apimachinery/pkg/runtime"
	"bufio"
	"bytes"
	"io"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/pkg/api/v1"
	rest "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
	"k8s.io/kubernetes/pkg/api"
	"os"
)

const (
	PodName       = "asb-1-zzwlt"
	PodNamespace  = "ansible-service-broker"
	ContainerName = "asb"
)

func main() {
	TestCommand := []string{"whoami"}

	log := NewLogger()

	k8sClient, err := clients.Kubernetes(log)
	if err != nil {
		fmt.Printf("error creating k8s client:")
		fmt.Printf("%+v", err.Error())
		return
	}

	clientConfig := k8sClient.ClientConfig
	clientConfig.GroupVersion = &v1.SchemeGroupVersion
	clientConfig.NegotiatedSerializer =
		serializer.DirectCodecFactory{CodecFactory: scheme.Codecs}

	// NOTE: kubectl exec simply sets the API path to /api when where is no
	// Group, which is the case for pod exec.
	clientConfig.APIPath = "/api"

	log.Infof("%+v\n", clientConfig)
	log.Infof("%s", string(debug.Stack()))

	restClient, err := rest.RESTClientFor(clientConfig)
	if err != nil {
		fmt.Printf("error creating rest client:")
		fmt.Printf("%+v", err.Error())
		return
	}

	req := restClient.Post().
		Resource("pods").
		Name(PodName).
		Namespace(PodNamespace).
		SubResource("exec").
		Param("container", ContainerName)

	req.VersionedParams(&api.PodExecOptions{
		Container: ContainerName,
		Command:   TestCommand,
		Stdin:     false,
		Stdout:    true,
		Stderr:    true,
		TTY:       false,
	}, api.ParameterCodec)

	exec, err := remotecommand.NewExecutor(clientConfig, "POST", req.URL())
	if err != nil {
		fmt.Printf("error getting new remotecommand executor")
		fmt.Printf("%+v", err.Error())
	}

	var stdoutBuffer, stderrBuffer bytes.Buffer
	stdoutWriter := bufio.NewWriter(&stdoutBuffer)
	stderrWriter := bufio.NewWriter(&stderrBuffer)

	err = exec.Stream(remotecommand.StreamOptions{
		Stdout: stdoutWriter,
		Stderr: stderrWriter,
	})

	// Flush?

	fmt.Printf("stdoutBuffer: [%s]", stdoutBuffer.String())
	fmt.Printf("stderrBuffer: [%s]", stderrBuffer.String())

	if err != nil {
		log.Error("Bad shit happened")
		log.Errorf("%+v", err.Error())
	}

	////////////////////////////////////////////////////////////
	// TODO:
	// try/finally to make sure we clean things up cleanly?
	//if stopsignal {
	//app.stop() // Stuff like close open files
	//}
	////////////////////////////////////////////////////////////
}

const LOGGING_MODULE = "cluster-sbx"

func NewLogger() *logging.Logger {
	logger := logging.MustGetLogger(LOGGING_MODULE)
	var backends []logging.Backend

	colorFormatter := logging.MustStringFormatter(
		"%{color}[%{time}] [%{level}] %{message}%{color:reset}",
	)

	standardFormatter := logging.MustStringFormatter(
		"[%{time}] [%{level}] %{message}",
	)

	var formattedBackend = func(writer io.Writer, isColored bool) logging.Backend {
		backend := logging.NewLogBackend(writer, "", 0)
		formatter := standardFormatter
		if isColored {
			formatter = colorFormatter
		}
		return logging.NewBackendFormatter(backend, formatter)
	}

	backends = append(backends, formattedBackend(os.Stdout, true /* isColored */))

	multiBackend := logging.MultiLogger(backends...)
	logger.SetBackend(multiBackend)
	logging.SetLevel(logging.DEBUG, LOGGING_MODULE)

	return logger
}
