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

package osb

import (
	"errors"
	schema "github.com/lestrrat/go-jsschema"
	"github.com/pborman/uuid"
)

var (
	// ErrorAlreadyProvisioned - Error for when an service instance has already been provisioned
	ErrorAlreadyProvisioned = errors.New("already provisioned")
	// ErrorDuplicate - Error for when a duplicate service instance already exists
	ErrorDuplicate = errors.New("duplicate instance")
	// ErrorNotFound  - Error for when a service instance is not found. (either etcd or kubernetes)
	ErrorNotFound = errors.New("not found")
	// ErrorBindingExists - Error for when deprovision is called on a service
	// instance with active bindings, or bind requested for already-existing
	// binding
	ErrorBindingExists = errors.New("binding exists")
	// ErrorProvisionInProgress - Error for when provision is called on a service instance that has a provision job in progress
	ErrorProvisionInProgress = errors.New("provision in progress")
	// ErrorDeprovisionInProgress - Error for when deprovision is called on a service instance that has a deprovision job in progress
	ErrorDeprovisionInProgress = errors.New("deprovision in progress")
	// ErrorUpdateInProgress - Error for when update is called on a service instance that has an update job in progress
	ErrorUpdateInProgress = errors.New("update in progress")
	// ErrorPlanNotFound - Error for when plan for update not found
	ErrorPlanNotFound = errors.New("plan not found")
	// ErrorParameterNotUpdatable - Error for when parameter in update request is not updatable
	ErrorParameterNotUpdatable = errors.New("parameter not updatable")
	// ErrorParameterNotFound - Error for when a parameter for update is not found
	ErrorParameterNotFound = errors.New("parameter not found")
	// ErrorPlanUpdateNotPossible - Error when a Plan Update request cannot be satisfied
	ErrorPlanUpdateNotPossible = errors.New("plan update not possible")
	// ErrorForbidden - Should be returned by broker handler if the user does not have sufficient permissions
	ErrorForbidden = errors.New("User does not have sufficient permissions")
)

// Parameters - generic string to object or value parameter
type Parameters map[string]interface{}

// Context - Determines the context in which the service is running
type Context struct {
	Platform  string `json:"platform"`
	// TODO: Namespace doesn't belong on an OSB type, need to find a better place for it in ASB.
	Namespace string `json:"namespace"`
}

// ServiceInstance - Service Instance describes a running service.
type ServiceInstance struct {
	ID         uuid.UUID       `json:"id"`
	PlanID     uuid.UUID       `json:"plan_id"`
	Context    *Context        `json:"context"`
	// TODO: Parameters doesn't need to be a pointer, it's a map (ref) naturally
	Parameters *Parameters     `json:"parameters"`
	BindingIDs map[string]bool `json:"binding_ids"`
}

// BindInstance - Binding Instance describes a completed binding
type BindInstance struct {
	ID         uuid.UUID   `json:"id"`
	ServiceID  uuid.UUID   `json:"service_id"`
	Parameters *Parameters `json:"parameters"`
}

// DashboardClient - Dashboard Client to be returned
// based on https://github.com/openservicebrokerapi/servicebroker/blob/master/spec.md#dashboard-client-object
type DashboardClient struct {
	ID          string `json:"id"`
	Secret      string `json:"secret"`
	RedirectURI string `json:"redirect_uri"`
}

// Service - Service object to be returned.
// based on https://github.com/openservicebrokerapi/servicebroker/blob/master/spec.md#service-objects
type Service struct {
	Name            string                 `json:"name"`
	ID              string                 `json:"id"`
	Description     string                 `json:"description"`
	Tags            []string               `json:"tags,omitempty"`
	Requires        []string               `json:"requires,omitempty"`
	Bindable        bool                   `json:"bindable"`
	Metadata        map[string]interface{} `json:"metadata,omitempty"`
	DashboardClient *DashboardClient       `json:"dashboard_client,omitempty"`
	PlanUpdatable   bool                   `json:"plan_updateable,omitempty"`
	Plans           []Plan                 `json:"plans"`
}

// Plan - Plan to be returned
// based on https://github.com/openservicebrokerapi/servicebroker/blob/master/spec.md#plan-object
type Plan struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	Free        bool                   `json:"free,omitempty"`
	Bindable    bool                   `json:"bindable,omitempty"`
	Schemas     Schema                 `json:"schemas,omitempty"`
	UpdatesTo   []string               `json:"updates_to,omitempty"`
}

// ServiceInstanceSchema - Schema definitions for creating and updating a service instance.
// Toyed with the idea of making an InputParameters
// that was a *schema.Schema
// based on 2.13 of the open service broker api. https://github.com/avade/servicebroker/blob/cda8c57b6a4bb7eaee84be20bb52dc155269758a/spec.md
type ServiceInstanceSchema struct {
	Create map[string]*schema.Schema `json:"create"`
	Update map[string]*schema.Schema `json:"update"`
}

// ServiceBindingSchema - Schema definitions for creating a service binidng.
type ServiceBindingSchema struct {
	Create map[string]*schema.Schema `json:"create"`
}

// Schema  - Schema to be returned
// based on 2.13 of the open service broker api. https://github.com/avade/servicebroker/blob/cda8c57b6a4bb7eaee84be20bb52dc155269758a/spec.md
type Schema struct {
	ServiceInstance ServiceInstanceSchema `json:"service_instance"`
	ServiceBinding  ServiceBindingSchema  `json:"service_binding"`
}

// -- Req/Response types

// CatalogResponse - Response for the catalog call.
// Defined here https://github.com/openservicebrokerapi/servicebroker/blob/v2.12/spec.md#response
type CatalogResponse struct {
	Services []Service `json:"services"`
}

// ProvisionRequest - Request for provision
// Defined here https://github.com/openservicebrokerapi/servicebroker/blob/v2.12/spec.md#request-2
type ProvisionRequest struct {
	OrganizationID    uuid.UUID  `json:"organization_guid"`
	PlanID            string     `json:"plan_id"`
	ServiceID         string     `json:"service_id"`
	SpaceID           uuid.UUID  `json:"space_guid"`
	Context           Context    `json:"context"`
	Parameters        Parameters `json:"parameters,omitempty"`
	AcceptsIncomplete bool       `json:"accepts_incomplete,omitempty"`
}

// ProvisionResponse - Response for provison
// Defined here https://github.com/openservicebrokerapi/servicebroker/blob/v2.12/spec.md#response-2
type ProvisionResponse struct {
	DashboardURL string `json:"dashboard_url,omitempty"`
	Operation    string `json:"operation,omitempty"`
}

// DeprovisionResponse - Response for a deprovision
// Defined here https://github.com/openservicebrokerapi/servicebroker/blob/v2.12/spec.md#response-6
type DeprovisionResponse struct {
	Operation string `json:"operation,omitempty"`
}

// BindRequest - Request for a bind
// Defined here https://github.com/openservicebrokerapi/servicebroker/blob/v2.12/spec.md#request-4
type BindRequest struct {
	ServiceID string `json:"service_id"`
	PlanID    string `json:"plan_id"`
	// Deprecated: AppID deprecated in favor of BindResource.AppID
	AppID uuid.UUID `json:"app_guid,omitempty"`

	BindResource struct {
		AppID uuid.UUID `json:"app_guid,omitempty"`
		Route string    `json:"route,omitempty"`
	} `json:"bind_resource,omitempty"`
	Parameters Parameters `json:"parameters,omitempty"`
}

// BindResponse - Response for a bind
// Defined here https://github.com/openservicebrokerapi/servicebroker/blob/v2.12/spec.md#response-4
type BindResponse struct {
	Credentials     map[string]interface{} `json:"credentials,omitempty"`
	SyslogDrainURL  string                 `json:"syslog_drain_url,omitempty"`
	RouteServiceURL string                 `json:"route_service_url,omitempty"`
	VolumeMounts    []interface{}          `json:"volume_mounts,omitempty"`
	Operation       string                 `json:"operation,omitempty"`
}

// UnbindResponse - Response for unbinding
// Defined here https://github.com/openservicebrokerapi/servicebroker/blob/v2.12/spec.md#response-5
type UnbindResponse struct {
	Operation string `json:"operation,omitempty"`
}

// UpdateRequest - Request for an update for a service instance.
// Defined here https://github.com/openservicebrokerapi/servicebroker/blob/v2.12/spec.md#request-3
type UpdateRequest struct {
	ServiceID      string            `json:"service_id"`
	PlanID         string            `json:"plan_id,omitempty"`
	Parameters     map[string]string `json:"parameters,omitempty"`
	PreviousValues struct {
		PlanID         string    `json:"plan_id,omitempty"`
		ServiceID      string    `json:"service_id,omitempty"`
		OrganizationID uuid.UUID `json:"organization_id,omitempty"`
		SpaceID        uuid.UUID `json:"space_id,omitempty"`
	} `json:"previous_values,omitempty"`
	Context           Context `json:"context"`
	AcceptsIncomplete bool    `json:"accepts_incomplete,omitempty"`
}

// UpdateResponse - Response for an update for a service instance.
// Defined here https://github.com/openservicebrokerapi/servicebroker/blob/v2.12/spec.md#response-3
type UpdateResponse struct {
	Operation string `json:"operation,omitempty"`
}

// LastOperationRequest - Request to obtain state information about an action that was taken
// Defined in more detail here https://github.com/openservicebrokerapi/servicebroker/blob/v2.12/spec.md#polling-last-operation
type LastOperationRequest struct {
	ServiceID string `json:"service_id"`
	PlanID    string `json:"plan_id"`
	Operation string `json:"operation"`
}

const (
	// LastOperationStateInProgress - In Progress state for last operation.
	LastOperationStateInProgress = "in progress"
	// LastOperationStateSucceeded - Succeeded state for the last operation.
	LastOperationStateSucceeded = "succeeded"
	// LastOperationStateFailed - Failed state for the last operation.
	LastOperationStateFailed = "failed"
)

// LastOperationResponse - Response for the laster operation request.
// Defined here https://github.com/openservicebrokerapi/servicebroker/blob/v2.12/spec.md#response-1
type LastOperationResponse struct {
	State       string `json:"state"`
	Description string             `json:"description,omitempty"`
}

// ServiceInstanceResponse - The response for a get service instance request
type ServiceInstanceResponse struct {
	ServiceID    string     `json:"service_id"`
	PlanID       string     `json:"plan_id"`
	DashboardURL string     `json:"dashboard_url,omitempty"`
	Parameters   Parameters `json:"parameters,omitempty"`
}

// ErrorResponse - Error response for all broker errors
// Defined here https://github.com/openservicebrokerapi/servicebroker/blob/v2.12/spec.md#broker-errors
type ErrorResponse struct {
	Description string `json:"description"`
}
