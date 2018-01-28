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

package server

import (
	"net/http"
	"github.com/gorilla/mux"
	"github.com/openshift/ansible-service-broker/pkg/hydro/osb"
	"github.com/openshift/ansible-service-broker/pkg/hydro/log"
	asb "github.com/openshift/ansible-service-broker/pkg/broker"
)

type AnsibleBrokerServer struct {
	prefix string
	kubernetesServer *KubernetesServer
}

func NewAnsibleBrokerServer() *AnsibleBrokerServer {
	ansibleBrokerPrefix := "/ansible-service-broker"
	return &AnsibleBrokerServer{
		prefix: ansibleBrokerPrefix,
		kubernetesServer: NewKubernetesServer(ansibleBrokerPrefix),
	}
}

func (s *AnsibleBrokerServer) Prefix() string {
	return s.prefix
}

func (s *AnsibleBrokerServer) ExtendRouter(broker osb.OpenServiceBroker, router *mux.Router) {
	// TODO: Need to return an error in case something goes wrong?
	if _, ok := broker.(asb.AnsibleBroker); !ok {
		log.Debug("Could not casst the OSB to an Ansible Broker...")
		log.Error("Something went wrong trying to extend the ansible broker server's router!")
	}
}

func (s *AnsibleBrokerServer) StartServer(h *http.Handler) {
	s.kubernetesServer.StartServer(h)
}
