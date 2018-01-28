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
	"encoding/base64"
	"github.com/openshift/ansible-service-broker/pkg/apb"
	"gopkg.in/yaml.v2"
	"github.com/openshift/ansible-service-broker/pkg/version"
)

type AnsibleBrokerServer struct {
	prefix string
	kubernetesServer *KubernetesServer
	broker *asb.AnsibleBroker
}

func NewAnsibleBrokerServer(broker *asb.AnsibleBroker) *AnsibleBrokerServer {
	ansibleBrokerPrefix := "/ansible-service-broker"
	return &AnsibleBrokerServer{
		broker: broker,
		prefix: ansibleBrokerPrefix,
		kubernetesServer: NewKubernetesServer(ansibleBrokerPrefix),
	}
}

func (s *AnsibleBrokerServer) Prefix() string {
	return s.prefix
}

func (s *AnsibleBrokerServer) ExtendRouter(broker osb.OpenServiceBroker, router *mux.Router) {
	ansibleBroker, ok := broker.(asb.AnsibleBroker)

	if !ok {
		log.Debug("Could not cast the OSB to an Ansible Broker...")
		log.Error("Something went wrong trying to extend the ansible broker server's router!")
		// TODO: Panic? Return error? What does hydro expect from brokers that error out while extending
		// the router?
	}

	if ansibleBroker.IsDevelopmentBroker() {
		router.HandleFunc("/v2/apb",
			createVarHandler(s.apbAddSpec)).Methods("POST")
		router.HandleFunc("/v2/apb/{spec_id}",
			createVarHandler(s.apbRemoveSpec)).Methods("DELETE")
		router.HandleFunc("/v2/apb",
			createVarHandler(s.apbRemoveSpecs)).Methods("DELETE")
	}
}

func (s *AnsibleBrokerServer) StartServer(h *http.Handler) {
	s.kubernetesServer.StartServer(h)
}

// apbAddSpec - Development only route. Will be used by for local developers to add images to the catalog.
func (s *AnsibleBrokerServer) apbAddSpec(w http.ResponseWriter, r *http.Request, params map[string]string) {
	log.Debug("handler::apbAddSpec")

	// Decode
	spec64Yaml := r.FormValue("apbSpec")
	if spec64Yaml == "" {
		log.Errorf("Could not find form parameter apbSpec")
		writeResponse(w, http.StatusBadRequest, asb.ErrorResponse{Description: "Could not parameter apbSpec"})
		return
	}
	decodedSpecYaml, err := base64.StdEncoding.DecodeString(spec64Yaml)
	if err != nil {
		log.Errorf("Something went wrong decoding spec from encoded string - %v err -%v", spec64Yaml, err)
		writeResponse(w, http.StatusBadRequest, osb.ErrorResponse{Description: "Invalid parameter encoding"})
		return
	}
	log.Debug("Successfully decoded pushed spec:")
	log.Debugf("%s", decodedSpecYaml)

	var spec apb.Spec
	if err = yaml.Unmarshal([]byte(decodedSpecYaml), &spec); err != nil {
		log.Errorf("Unable to decode yaml - %v to spec err - %v", decodedSpecYaml, err)
		writeResponse(w, http.StatusBadRequest, osb.ErrorResponse{Description: "Invalid parameter yaml"})
		return
	}
	log.Infof("Assuming pushed APB runtime version [%v]", version.MaxRuntimeVersion)
	spec.Runtime = version.MaxRuntimeVersion

	log.Debug("Unmarshalled into apb.Spec:")
	log.Debugf("%+v", spec)

	resp, err := s.broker.AddSpec(spec)
	if err != nil {
		log.Errorf("An error occurred while trying to add a spec via apb push:")
		log.Errorf("%s", err.Error())
		writeResponse(w, http.StatusInternalServerError,
			asb.ErrorResponse{Description: err.Error()})
		return
	}

	writeDefaultResponse(w, http.StatusOK, resp, err)
}

func (s *AnsibleBrokerServer) apbRemoveSpec(w http.ResponseWriter, r *http.Request, params map[string]string) {
	specID := params["spec_id"]

	var err error
	if specID != "" {
		err = s.broker.RemoveSpec(specID)
	} else {
		log.Errorf("Unable to find spec id in request")
		writeResponse(w, http.StatusBadRequest, osb.ErrorResponse{Description: "No Spec/service id found."})
		return
	}
	writeDefaultResponse(w, http.StatusNoContent, struct{}{}, err)
}

func (s *AnsibleBrokerServer) apbRemoveSpecs(w http.ResponseWriter, r *http.Request, params map[string]string) {
	err := s.broker.RemoveSpecs()
	writeDefaultResponse(w, http.StatusNoContent, struct{}{}, err)
}
