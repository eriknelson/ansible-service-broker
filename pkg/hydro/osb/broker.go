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
	"context"
	"github.com/pborman/uuid"
)

type OpenServiceBroker interface {
	Catalog() (*CatalogResponse, error)
	Provision(uuid.UUID, *ProvisionRequest, bool, context.Context) (*ProvisionResponse, error)
	Deprovision(ServiceInstance, string, bool, context.Context) (*DeprovisionResponse, error)
	Bind(ServiceInstance, uuid.UUID, *BindRequest, bool, context.Context) (*BindResponse, bool, error)
	Unbind(ServiceInstance, BindInstance, string, bool, context.Context) (*UnbindResponse, error)
	Update(uuid.UUID, *UpdateRequest, bool, context.Context) (*UpdateResponse, error)
	LastOperation(uuid.UUID, *LastOperationRequest) (*LastOperationResponse, error)
	GetServiceInstance(uuid.UUID) (*ServiceInstance, error)
	GetBindInstance(uuid.UUID) (*BindInstance, error)

	// NOT OSB OPS, from ASB
	//Recover() (string, error)
	//GetServiceInstance(uuid.UUID) (apb.ServiceInstance, error)
	//GetBindInstance(uuid.UUID) (apb.BindInstance, error)
	//GetBind(apb.ServiceInstance, uuid.UUID) (*BindResponse, error)
}
