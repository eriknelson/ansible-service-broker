//
// Copyright (c) 2017 Red Hat, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package server

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/openshift/ansible-service-broker/pkg/hydro/log"
	"github.com/openshift/ansible-service-broker/pkg/hydro/osb"
	"github.com/pborman/uuid"
	"net/http"
	"net/http/httputil"
	"strconv"
	"strings"
)

type handler struct {
	router *mux.Router
	broker osb.OpenServiceBroker
	config HandlerConfig
}

// HandlerConfig - Allows for some configuration of the OSB handler
type HandlerConfig struct {
	// RequestDebug - If true, will log requests
	RequestDebug bool
	// Prefix - The prefix that the router will be mounted behind
	Prefix string
}

// GorillaRouteHandler - gorilla route handler
// making the handler methods more testable by moving the reliance of mux.Vars()
// outside of the handlers themselves
type GorillaRouteHandler func(http.ResponseWriter, *http.Request)

// VarHandler - Variable route handler.
type VarHandler func(http.ResponseWriter, *http.Request, map[string]string)

// NewHandler - Creates a new
func NewOpenServiceBrokerHandler(broker osb.OpenServiceBroker, config HandlerConfig) *mux.Router {
	// TODO: Any auth stuff here?

	h := handler{
		router: mux.NewRouter(),
		broker: broker,
		config: config,
	}

	var r *mux.Router
	if h.config.Prefix == "/" {
		r = h.router
	} else {
		r = h.router.PathPrefix(config.Prefix).Subrouter()
	}

	r.HandleFunc("/v2/catalog", createVarHandler(h.catalog)).Methods("GET")
	r.HandleFunc("/v2/service_instances/{instance_uuid}",
		createVarHandler(h.provision)).Methods("PUT")
	r.HandleFunc("/v2/service_instances/{instance_uuid}",
		createVarHandler(h.getServiceInstance)).Methods("GET")
	r.HandleFunc("/v2/service_instances/{instance_uuid}",
		createVarHandler(h.deprovision)).Methods("DELETE")
	r.HandleFunc("/v2/service_instances/{instance_uuid}/service_bindings/{binding_uuid}",
		createVarHandler(h.bind)).Methods("PUT")
	r.HandleFunc("/v2/service_instances/{instance_uuid}/service_bindings/{binding_uuid}",
		createVarHandler(h.unbind)).Methods("DELETE")
	r.HandleFunc("/v2/service_instances/{instance_uuid}/service_bindings/{binding_uuid}",
		createVarHandler(h.getBindInstance)).Methods("GET")
	r.HandleFunc("/v2/service_instances/{instance_uuid}",
		createVarHandler(h.update)).Methods("PATCH")
	r.HandleFunc("/v2/service_instances/{instance_uuid}/last_operation",
		createVarHandler(h.lastoperation)).Methods("GET")
	r.HandleFunc("/v2/service_instances/{instance_uuid}/service_bindings/{binding_uuid}/last_operation",
		createVarHandler(h.lastoperation)).Methods("GET")

	// TODO: Push to broker server router extension
	//if brokerConfig.GetBool("broker.dev_broker") {
	//	s.HandleFunc("/v2/apb", createVarHandler(h.apbAddSpec)).Methods("POST")
	//	s.HandleFunc("/v2/apb/{spec_id}", createVarHandler(h.apbRemoveSpec)).Methods("DELETE")
	//	s.HandleFunc("/v2/apb", createVarHandler(h.apbRemoveSpecs)).Methods("DELETE")
	//}

	return r // handlers.LoggingHandler(os.Stdout, r)
}

func createVarHandler(r VarHandler) GorillaRouteHandler {
	return func(writer http.ResponseWriter, request *http.Request) {
		r(writer, request, mux.Vars(request))
	}
}

func (h handler) catalog(w http.ResponseWriter, r *http.Request, params map[string]string) {
	defer r.Body.Close()
	h.printRequest(r)

	resp, err := h.broker.Catalog()

	writeDefaultResponse(w, http.StatusOK, resp, err)
}

func (h handler) provision(w http.ResponseWriter, r *http.Request, params map[string]string) {
	defer r.Body.Close()
	h.printRequest(r)

	instanceUUID := uuid.Parse(params["instance_uuid"])
	if instanceUUID == nil {
		writeResponse(w, http.StatusBadRequest, osb.ErrorResponse{Description: "invalid instance_uuid"})
		return
	}

	// ignore the error, if async can't be parsed it will be false
	async, _ := strconv.ParseBool(r.FormValue("accepts_incomplete"))

	var req *osb.ProvisionRequest
	err := readRequest(r, &req)

	if err != nil {
		writeResponse(w, http.StatusBadRequest, osb.ErrorResponse{Description: "could not read request: " + err.Error()})
		return
	}

	// TODO: Need to push this into the ASB!
	// if !h.brokerConfig.GetBool("broker.auto_escalate") {
	//	userInfo, ok := r.Context().Value(UserInfoContext).(broker.UserInfo)
	//	if !ok {
	//		log.Debugf("unable to retrieve user info from request context")
	//		// if no user, we should error out with bad request.
	//		writeResponse(w, http.StatusBadRequest, broker.osb.ErrorResponse{
	//			Description: "Invalid user info from originating origin header.",
	//		})
	//		return
	//	}
	//
	//	if ok, status, err := h.validateUser(userInfo, req.Context.Namespace); !ok {
	//		writeResponse(w, status, broker.osb.ErrorResponse{Description: err.Error()})
	//		return
	//	}
	//} else {
	//	log.Debugf("Auto Escalate has been set to true, we are escalating permissions")
	//}

	resp, err := h.broker.Provision(instanceUUID, req, async)

	if err != nil {
		switch err {
		case osb.ErrorDuplicate:
			writeResponse(w, http.StatusConflict, osb.ProvisionResponse{})
		case osb.ErrorProvisionInProgress:
			writeResponse(w, http.StatusAccepted, resp)
		case osb.ErrorAlreadyProvisioned:
			writeResponse(w, http.StatusOK, resp)
		case osb.ErrorNotFound:
			writeResponse(w, http.StatusBadRequest, osb.ErrorResponse{Description: err.Error()})
		default:
			writeResponse(w, http.StatusBadRequest, osb.ErrorResponse{Description: err.Error()})
		}
	} else if async {
		writeDefaultResponse(w, http.StatusAccepted, resp, err)
	} else {
		writeDefaultResponse(w, http.StatusCreated, resp, err)
	}
}

func (h handler) deprovision(w http.ResponseWriter, r *http.Request, params map[string]string) {
	defer r.Body.Close()
	h.printRequest(r)

	instanceUUID := uuid.Parse(params["instance_uuid"])
	if instanceUUID == nil {
		writeResponse(w, http.StatusBadRequest, osb.ErrorResponse{Description: "invalid instance_uuid"})
		return
	}

	// ignore the error, if async can't be parsed it will be false
	async, _ := strconv.ParseBool(r.FormValue("accepts_incomplete"))

	planID := r.FormValue("plan_id")
	if planID == "" {
		writeResponse(w, http.StatusBadRequest, osb.ErrorResponse{Description: "deprovision request missing plan_id query parameter"})
	}

	serviceInstance, err := h.broker.GetServiceInstance(instanceUUID)
	if err != nil {
		switch err {
		case osb.ErrorNotFound:
			writeResponse(w, http.StatusGone, osb.DeprovisionResponse{})
			return
		default:
			writeResponse(w, http.StatusInternalServerError, osb.ErrorResponse{Description: err.Error()})
			return
		}
	}

	// TODO: Need to push this into ASB
	//nsDeleted, err := isNamespaceDeleted(serviceInstance.Context.Namespace)
	//if err != nil {
	//	writeResponse(w, http.StatusInternalServerError, broker.osb.ErrorResponse{Description: err.Error()})
	//	return
	//}
	//
	//if !h.brokerConfig.GetBool("broker.auto_escalate") {
	//	userInfo, ok := r.Context().Value(UserInfoContext).(broker.UserInfo)
	//	if !ok {
	//		log.Debugf("unable to retrieve user info from request context")
	//		// if no user, we should error out with bad request.
	//		writeResponse(w, http.StatusBadRequest, broker.osb.ErrorResponse{
	//			Description: "Invalid user info from originating origin header.",
	//		})
	//		return
	//	}
	//
	//	if !nsDeleted {
	//		ok, status, err := h.validateUser(userInfo, serviceInstance.Context.Namespace)
	//		if !ok {
	//			writeResponse(w, status, broker.osb.ErrorResponse{Description: err.Error()})
	//			return
	//		}
	//	}
	//} else {
	//	log.Debugf("Auto Escalate has been set to true, we are escalating permissions")
	//}

	resp, err := h.broker.Deprovision(*serviceInstance, planID, async)

	if err != nil {
		switch err {
		case osb.ErrorNotFound:
			writeResponse(w, http.StatusGone, osb.DeprovisionResponse{})
			return
		case osb.ErrorBindingExists:
			writeResponse(w, http.StatusBadRequest, osb.DeprovisionResponse{})
			return
		case osb.ErrorDeprovisionInProgress:
			writeResponse(w, http.StatusAccepted, osb.DeprovisionResponse{})
			return
		}
	} else if async {
		writeDefaultResponse(w, http.StatusAccepted, resp, err)
	} else {
		writeDefaultResponse(w, http.StatusCreated, resp, err)
	}
}

func (h handler) getServiceInstance(w http.ResponseWriter, r *http.Request, params map[string]string) {
	defer r.Body.Close()
	h.printRequest(r)

	instanceUUID := uuid.Parse(params["instance_uuid"])
	if instanceUUID == nil {
		writeResponse(w, http.StatusBadRequest, osb.ErrorResponse{Description: "invalid instance_uuid"})
		return
	}

	serviceInstance, err := h.broker.GetServiceInstance(instanceUUID)
	if err != nil {
		switch err {
		case osb.ErrorNotFound: // return 404
			writeResponse(w, http.StatusNotFound, osb.ErrorResponse{Description: err.Error()})
		default: // return 422
			writeResponse(w, http.StatusUnprocessableEntity, osb.ErrorResponse{Description: err.Error()})
		}
		return
	}

	response := osb.ServiceInstanceResponse{
		ServiceID:  serviceInstance.ID.String(),
		PlanID:     serviceInstance.PlanID.String(),
		Parameters: *serviceInstance.Parameters,
	}

	writeDefaultResponse(w, http.StatusOK, response, err)
}

func (h handler) bind(w http.ResponseWriter, r *http.Request, params map[string]string) {
	defer r.Body.Close()
	h.printRequest(r)

	instanceUUID := uuid.Parse(params["instance_uuid"])
	if instanceUUID == nil {
		writeResponse(w, http.StatusBadRequest, osb.ErrorResponse{Description: "invalid instance_uuid"})
		return
	}

	bindingUUID := uuid.Parse(params["binding_uuid"])
	if bindingUUID == nil {
		writeResponse(w, http.StatusBadRequest, osb.ErrorResponse{Description: "invalid binding_uuid"})
		return
	}

	// ignore the error, if async can't be parsed it will be false
	async, _ := strconv.ParseBool(r.FormValue("accepts_incomplete"))

	// TODO: Push to ASB
	//if !async && h.brokerConfig.GetBool("broker.launch_apb_on_bind") {
	//	log.Warning("launch_apb_on_bind is enabled, but accepts_incomplete is false, binding may fail")
	//}

	var req *osb.BindRequest
	if err := readRequest(r, &req); err != nil {
		writeResponse(w, http.StatusInternalServerError, osb.ErrorResponse{Description: err.Error()})
		return
	}

	serviceInstance, err := h.broker.GetServiceInstance(instanceUUID)
	if err != nil {
		switch err {
		case osb.ErrorNotFound:
			writeResponse(w, http.StatusBadRequest, osb.ErrorResponse{Description: err.Error()})
		default:
			writeResponse(w, http.StatusInternalServerError, osb.ErrorResponse{Description: err.Error()})
		}
		return
	}

	// TODO: Push to ASB
	//if !h.brokerConfig.GetBool("broker.auto_escalate") {
	//	userInfo, ok := r.Context().Value(UserInfoContext).(broker.UserInfo)
	//	if !ok {
	//		log.Debugf("unable to retrieve user info from request context")
	//		// if no user, we should error out with bad request.
	//		writeResponse(w, http.StatusBadRequest, broker.osb.ErrorResponse{
	//			Description: "Invalid user info from originating origin header.",
	//		})
	//		return
	//	}
	//
	//	if ok, status, err := h.validateUser(userInfo, serviceInstance.Context.Namespace); !ok {
	//		writeResponse(w, status, broker.osb.ErrorResponse{Description: err.Error()})
	//		return
	//	}
	//} else {
	//	log.Debugf("Auto Escalate has been set to true, we are escalating permissions")
	//}

	resp, ranAsync, err := h.broker.Bind(*serviceInstance, bindingUUID, req, async)

	if err != nil {
		switch err {
		case osb.ErrorDuplicate:
			writeResponse(w, http.StatusConflict, osb.BindResponse{})
		case osb.ErrorBindingExists:
			writeResponse(w, http.StatusOK, resp)
		case osb.ErrorNotFound:
			writeResponse(w, http.StatusBadRequest, osb.ErrorResponse{Description: err.Error()})
		default:
			writeResponse(w, http.StatusBadRequest, osb.ErrorResponse{Description: err.Error()})
		}
		return
	}
	if ranAsync {
		writeDefaultResponse(w, http.StatusAccepted, resp, err)
	} else {
		writeDefaultResponse(w, http.StatusCreated, resp, err)
	}
}

func (h handler) unbind(w http.ResponseWriter, r *http.Request, params map[string]string) {
	defer r.Body.Close()
	h.printRequest(r)

	instanceUUID := uuid.Parse(params["instance_uuid"])
	if instanceUUID == nil {
		writeResponse(w, http.StatusBadRequest, osb.ErrorResponse{Description: "invalid instance_uuid"})
		return
	}

	bindingUUID := uuid.Parse(params["binding_uuid"])
	if bindingUUID == nil {
		writeResponse(w, http.StatusBadRequest, osb.ErrorResponse{Description: "invalid binding_uuid"})
		return
	}
	planID := r.FormValue("plan_id")
	if planID == "" {
		writeResponse(w, http.StatusBadRequest, osb.ErrorResponse{Description: "unbind request missing plan_id query parameter"})
	}

	// ignore the error, if async can't be parsed it will be false
	async, _ := strconv.ParseBool(r.FormValue("accepts_incomplete"))

	// TODO: push to ASB
	//if !async && h.brokerConfig.GetBool("broker.launch_apb_on_bind") {
	//	log.Warning("launch_apb_on_bind is enabled, but accepts_incomplete is false, unbinding may fail")
	//}

	serviceInstance, err := h.broker.GetServiceInstance(instanceUUID)
	if err != nil {
		switch err {
		case osb.ErrorNotFound:
			writeResponse(w, http.StatusGone, nil)
		default:
			writeResponse(w, http.StatusInternalServerError, osb.ErrorResponse{Description: err.Error()})
		}
		return
	}

	bindInstance, err := h.broker.GetBindInstance(bindingUUID)
	if err != nil {
		switch err {
		case osb.ErrorNotFound:
			writeResponse(w, http.StatusGone, nil)
		default:
			writeResponse(w, http.StatusInternalServerError, osb.ErrorResponse{Description: err.Error()})
		}
		return
	}

	// TODO: Push to ASB
	//nsDeleted, err := isNamespaceDeleted(serviceInstance.Context.Namespace)
	//if err != nil {
	//	writeResponse(w, http.StatusInternalServerError, broker.osb.ErrorResponse{Description: err.Error()})
	//	return
	//}

	//if !h.brokerConfig.GetBool("broker.auto_escalate") {
	//	userInfo, ok := r.Context().Value(UserInfoContext).(broker.UserInfo)
	//	if !ok {
	//		log.Debugf("unable to retrieve user info from request context")
	//		// if no user, we should error out with bad request.
	//		writeResponse(w, http.StatusBadRequest, broker.osb.ErrorResponse{
	//			Description: "Invalid user info from originating origin header.",
	//		})
	//		return
	//	}
	//	if !nsDeleted {
	//		if ok, status, err := h.validateUser(userInfo, serviceInstance.Context.Namespace); !ok {
	//			writeResponse(w, status, broker.osb.ErrorResponse{Description: err.Error()})
	//			return
	//		}
	//	}
	//} else {
	//	log.Debugf("Auto Escalate has been set to true, we are escalating permissions")
	//}

	resp, err := h.broker.Unbind(*serviceInstance, *bindInstance, planID, async)

	if err != nil {
		switch err {
		case osb.ErrorNotFound: // return 404
			log.Debugf("Binding not found.")
			writeResponse(w, http.StatusNotFound, osb.ErrorResponse{Description: err.Error()})
		default: // return 500
			log.Errorf("Unknown error: %v", err)
			writeResponse(w, http.StatusInternalServerError, osb.ErrorResponse{Description: err.Error()})
		}
		return
	}

	writeDefaultResponse(w, http.StatusOK, resp, err)
	return
}

func (h handler) getBindInstance(w http.ResponseWriter, r *http.Request, params map[string]string) {
	defer r.Body.Close()
	h.printRequest(r)

	instanceUUID := uuid.Parse(params["instance_uuid"])
	if instanceUUID == nil {
		writeResponse(w, http.StatusBadRequest, osb.ErrorResponse{Description: "invalid instance_uuid"})
		return
	}

	bindingUUID := uuid.Parse(params["binding_uuid"])
	if bindingUUID == nil {
		writeResponse(w, http.StatusBadRequest, osb.ErrorResponse{Description: "invalid binding_uuid"})
		return
	}

	bindInstance, err := h.broker.GetBindInstance(bindingUUID)

	if err != nil {
		switch err {
		case osb.ErrorNotFound:
			writeResponse(w, http.StatusNotFound, osb.ErrorResponse{Description: err.Error()})
		default:
			writeResponse(w, http.StatusBadRequest, osb.ErrorResponse{Description: err.Error()})
		}
		return
	}

	writeDefaultResponse(w, http.StatusOK, bindInstance, err)
}

func (h handler) update(w http.ResponseWriter, r *http.Request, params map[string]string) {
	defer r.Body.Close()
	h.printRequest(r)

	instanceUUID := uuid.Parse(params["instance_uuid"])
	if instanceUUID == nil {
		writeResponse(w, http.StatusBadRequest, osb.ErrorResponse{Description: "invalid instance_uuid"})
		return
	}

	var req *osb.UpdateRequest

	if err := readRequest(r, &req); err != nil {
		writeResponse(w, http.StatusBadRequest, osb.ErrorResponse{Description: err.Error()})
		return
	}

	// ignore the error, if async can't be parsed it will be false
	async, _ := strconv.ParseBool(r.FormValue("accepts_incomplete"))

	//if !h.brokerConfig.GetBool("broker.auto_escalate") {
	//	userInfo, ok := r.Context().Value(UserInfoContext).(broker.UserInfo)
	//	if !ok {
	//		log.Debugf("unable to retrieve user info from request context")
	//		// if no user, we should error out with bad request.
	//		writeResponse(w, http.StatusBadRequest, broker.osb.ErrorResponse{
	//			Description: "Invalid user info from originating origin header.",
	//		})
	//		return
	//	}

	//	if ok, status, err := h.validateUser(userInfo, req.Context.Namespace); !ok {
	//		writeResponse(w, status, broker.osb.ErrorResponse{Description: err.Error()})
	//		return
	//	}
	//} else {
	//	log.Debugf("Auto Escalate has been set to true, we are escalating permissions")
	//}

	resp, err := h.broker.Update(instanceUUID, req, async)

	if err != nil {
		switch err {
		case osb.ErrorUpdateInProgress:
			writeResponse(w, http.StatusAccepted, resp)
		case osb.ErrorNotFound:
			writeResponse(w, http.StatusBadRequest, osb.ErrorResponse{Description: err.Error()})
		default:
			writeResponse(w, http.StatusBadRequest, osb.ErrorResponse{Description: err.Error()})
		}
	} else if async {
		writeDefaultResponse(w, http.StatusAccepted, resp, err)
	} else {
		writeDefaultResponse(w, http.StatusOK, resp, err)
	}
}

func (h handler) lastoperation(w http.ResponseWriter, r *http.Request, params map[string]string) {
	defer r.Body.Close()
	h.printRequest(r)

	instanceUUID := uuid.Parse(params["instance_uuid"])
	if instanceUUID == nil {
		writeResponse(w, http.StatusBadRequest, osb.ErrorResponse{Description: "invalid instance_uuid"})
		return
	}

	if strings.Index(r.URL.Path, "/service_bindings/") > 0 {
		bindingUUID := uuid.Parse(params["binding_uuid"])
		if bindingUUID == nil {
			writeResponse(w, http.StatusBadRequest, osb.ErrorResponse{Description: "invalid binding_uuid"})
			return
		}

		// let's see if the bindInstance exists or not. We don't need the
		// actual instance, just need to know if it is there.
		_, err := h.broker.GetBindInstance(bindingUUID)
		if err != nil {
			switch err {
			case osb.ErrorNotFound:
				writeResponse(w, http.StatusGone, nil)
			default:
				writeResponse(w, http.StatusInternalServerError, osb.ErrorResponse{Description: err.Error()})
			}
			return
		}
	}

	req := osb.LastOperationRequest{}

	// operation is expected
	if op := r.FormValue("operation"); op != "" {
		req.Operation = op
	} else {
		errmsg := fmt.Sprintf("operation not supplied for a last_operation with instance_uuid [%s]", instanceUUID)
		log.Error(errmsg)
		writeResponse(w, http.StatusBadRequest, osb.ErrorResponse{Description: errmsg})
		return
	}

	// service_id is optional
	if serviceID := r.FormValue("service_id"); serviceID != "" {
		req.ServiceID = serviceID
	}

	// plan_id is optional
	if planID := r.FormValue("plan_id"); planID != "" {
		req.PlanID = planID
	}

	resp, err := h.broker.LastOperation(instanceUUID, &req)

	if err != nil {
		switch err {
		default:
			writeResponse(w, http.StatusInternalServerError, osb.ErrorResponse{Description: err.Error()})
		}
		return
	}

	writeDefaultResponse(w, http.StatusOK, resp, err)
}

// printRequest - will print the request with the body.
func (h handler) printRequest(req *http.Request) {
	if h.config.RequestDebug {
		b, err := httputil.DumpRequest(req, true)
		if err != nil {
			log.Errorf("unable to dump request to log: %v", err)
		}
		log.Infof("Request: %q", b)
	}
}

func readRequest(r *http.Request, obj interface{}) error {
	if r.Header.Get("Content-Type") != "application/json" {
		return errors.New("error: invalid content-type")
	}

	return json.NewDecoder(r.Body).Decode(&obj)
}

func writeResponse(w http.ResponseWriter, code int, obj interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	b, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	i := bytes.Buffer{}
	json.Indent(&i, b, "", "  ")
	i.WriteString("\n")
	_, err = w.Write(i.Bytes())
	return err
}
func writeDefaultResponse(w http.ResponseWriter, code int, resp interface{}, err error) error {
	if err == nil {
		return writeResponse(w, code, resp)
	}

	return writeResponse(w, http.StatusInternalServerError, osb.ErrorResponse{Description: err.Error()})
}
